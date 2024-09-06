package main

import (
 "context"
 "log"
 "fmt"
 "net/http"

 "github.com/go-redis/redis/v8"
 "github.com/gorilla/mux"
 "github.com/gorilla/websocket"
)

var (
 ctx = context.Background()
 rdb = redis.NewClient(&redis.Options{
  Addr: "172.17.0.3:6379",
 })
 upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
   return true
  },
 }
)

func main() {

    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        fmt.Println("Error connecting to Redis:", err)
        return
    }
    fmt.Println("Connected to Redis:", pong)

 r := mux.NewRouter()
 r.HandleFunc("/ws/{channel}", handleWebSocket)

 fs := http.FileServer(http.Dir("./web"))
 r.PathPrefix("/").Handler(fs)

 log.Println("Server started on :3000")
 log.Fatal(http.ListenAndServe(":3000", r))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
 vars := mux.Vars(r)
 channel := vars["channel"]

 conn, err := upgrader.Upgrade(w, r, nil)
 if err != nil {
  log.Println(err)
  return
 }
 defer conn.Close()

 sub := rdb.Subscribe(ctx, channel)
 defer sub.Close()
 ch := sub.Channel()

 go func() {
  for {
   _, msg, err := conn.ReadMessage()
   if err != nil {
    log.Println("Read error:", err)
    return
   }

   if err := rdb.Publish(ctx, channel, string(msg)).Err(); err != nil {
    log.Println("Publish error:", err)
    return
   }
  }
 }()

 for msg := range ch {
  if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
   log.Println("Write error:", err)
   return
  }
 }
}