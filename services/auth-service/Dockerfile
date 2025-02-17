FROM golang:1.23.6-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/

RUN apk add --no-cache ca-certificates

RUN wget -qO /usr/local/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.16/grpc_health_probe-linux-amd64 \
    && chmod +x /usr/local/bin/grpc_health_probe

COPY --from=builder /app/auth-service .

EXPOSE 50051

HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD grpc_health_probe -addr=localhost:50051 || exit 1

CMD ["./auth-service"]