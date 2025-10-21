FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /todo-app .

EXPOSE 8080

CMD ["./todo-app"]
