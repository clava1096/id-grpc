# Этап сборки Go приложения
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/http cmd/api/main.go

# Этап создания финального образа
FROM alpine:3 AS http

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/http /app/service
# Копируем конфигурационные файлы (если есть)
# COPY --from=builder /app/config.yaml /app/

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser


CMD ["/app/service"]