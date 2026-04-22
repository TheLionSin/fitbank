# СТАДИЯ 1: Сборка (Builder)
FROM golang:1.26-alpine AS builder

# Инструменты сборки
RUN apk add --no-cache git

WORKDIR /app

# 1. Кэшируем зависимости всего воркспейса
COPY go.work ./
COPY activity-service/go.mod activity-service/go.sum ./activity-service/
COPY nutrition-service/go.mod nutrition-service/go.sum ./nutrition-service/

# Скачиваем зависимости (этот слой закешируется)
RUN go mod download

# 2. Копируем весь код проекта
COPY . .

# 3. ПЕРЕХОДИМ В ПАПКУ ИМЕННО ЭТОГО СЕРВИСА
WORKDIR /app/activity-service

# Собираем бинарник (без tidy, он уже должен быть сделан в WSL)
RUN CGO_ENABLED=0 GOOS=linux go build -o activity-service ./cmd/server/main.go


# СТАДИЯ 2: Финальный образ (Runner)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник (он лежит в папке сервиса внутри builder)
COPY --from=builder /app/activity-service/activity-service .

# КРИТИЧНО: Копируем миграции из папки сервиса, а не из корня!
COPY --from=builder /app/activity-service/migrations ./migrations

EXPOSE 8080 50051

CMD ["./activity-service"]