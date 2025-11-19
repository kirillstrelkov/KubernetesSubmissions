FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go get && CGO_ENABLED=0 GOOS=linux go build -o /todo-backend

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /todo-backend ./

CMD ["./todo-backend"]
