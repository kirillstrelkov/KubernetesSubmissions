FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /log_output_reader

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /log_output_reader .

CMD ["./log_output_reader"]
