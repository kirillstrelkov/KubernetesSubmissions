FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /greeter

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /greeter .

CMD ["./greeter"]
