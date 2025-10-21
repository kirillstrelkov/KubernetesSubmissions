FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /log_output

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /log_output .

CMD ["./log_output"]
