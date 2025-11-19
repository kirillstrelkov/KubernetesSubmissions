FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o /ping-pong-app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /ping-pong-app .

EXPOSE 8080

CMD ["./ping-pong-app"]
