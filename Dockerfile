FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src .

RUN go build -o main .

FROM debian:bullseye-slim

COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 3000

CMD ["./main"]
