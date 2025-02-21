# Этап 1: Сборка приложения
FROM golang:1.23.6 AS builder

WORKDIR /app

# Копирование go.mod и go.sum из корня проекта
COPY go.mod go.sum ./
RUN go mod download

# Копирование всего исходного кода
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./services/auth-service/cmd/server/main.go

RUN GO111MODULE=on go get github.com/grpc-ecosystem/grpc-health-probe@v0.4.28 && \
    GO111MODULE=on go build -o /go/bin/grpc_health_probe github.com/grpc-ecosystem/grpc-health-probe

# Этап 2: Создание образа приложения
FROM alpine:latest AS final

WORKDIR /

COPY --from=builder /auth-service /auth-service
COPY --from=builder /go/bin/grpc_health_probe /usr/local/bin/grpc_health_probe
COPY services/auth-service/internal/config/config.dev.yaml /config/config.dev.yaml

RUN apk --no-cache add ca-certificates

EXPOSE 50051

ENTRYPOINT ["/auth-service"]
CMD ["--config=/config/config.docker.yaml"]