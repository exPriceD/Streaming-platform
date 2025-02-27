# Используем официальный образ Go для сборки
FROM golang:1.21 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код внутрь контейнера
COPY . .

# Компилируем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o streaming-service ./cmd/main.go

# Используем минимальный образ для финального контейнера
FROM debian:bullseye-slim

# Устанавливаем необходимые зависимости (FFmpeg и Nginx RTMP)
RUN apt-get update && apt-get install -y \
    ffmpeg \
    nginx && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Копируем скомпилированный бинарник из builder-контейнера
COPY --from=builder /app/streaming-service /usr/local/bin/streaming-service

# Копируем конфигурацию Nginx RTMP
COPY configs/nginx-rtmp.conf /etc/nginx/nginx.conf

# Создаём директорию для HLS сегментов
RUN mkdir -p /var/www/hls && chown -R www-data:www-data /var/www/hls

# Открываем порты RTMP и HTTP
EXPOSE 1935 8080

# Запускаем Nginx RTMP и сам сервис
CMD service nginx start && streaming-service
