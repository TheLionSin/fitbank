# СТАДИЯ 1: Сборка (Builder)
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# 1. Кэшируем зависимости
COPY go.work ./
COPY activity-service/go.mod activity-service/go.sum ./activity-service/
COPY nutrition-service/go.mod nutrition-service/go.sum ./nutrition-service/

RUN go mod download

# 2. Копируем весь код
COPY . .

# 3. Переходим в папку питания
WORKDIR /app/nutrition-service

# Собираем
RUN CGO_ENABLED=0 GOOS=linux go build -o nutrition-service ./cmd/server/main.go


# СТАДИЯ 2: Финальный образ (Runner)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из подпапки builder
COPY --from=builder /app/nutrition-service/nutrition-service .

# КРИТИЧНО: Копируем миграции именно этого сервиса
COPY --from=builder /app/nutrition-service/migrations ./migrations

EXPOSE 8081 50052

CMD ["./nutrition-service"]