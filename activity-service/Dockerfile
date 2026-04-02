# СТАДИЯ 1: Сборка (Builder)
FROM golang:1.26-alpine AS builder

# Устанавливаем необходимые инструменты для сборки (если нужны)
RUN apk add --no-cache git

WORKDIR /app

# Сначала копируем зависимости, чтобы Docker мог их закэшировать
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл.
# CGO_ENABLED=0 делает бинарник статическим (не нужны системные библиотеки)
# GOOS=linux говорит собирать под Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o activity-service ./cmd/server/main.go

# СТАДИЯ 2: Финальный образ (Runner)
# Используем минимальный образ Alpine (всего 5МБ)
FROM alpine:latest

# Добавляем сертификаты (нужны для HTTPS запросов, если будут)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем ТОЛЬКО скомпилированный файл из стадии builder
COPY --from=builder /app/activity-app .
# Копируем папку с миграциями, так как они нужны в рантайме
COPY --from=builder /app/migrations ./migrations

# Открываем порты
EXPOSE 8080 50051

# Запускаем наше приложение
CMD ["./activity-service"]