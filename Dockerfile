FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app

EXPOSE 3000

CMD ["go","run","main.go"]
