FROM ubuntu:latest
LABEL authors="alexander"

ENTRYPOINT ["top", "-b"]

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Final stage
FROM alpine:3.18

WORKDIR /app

# Устанавливаем зависимости времени выполнения
RUN apk --no-cache add ca-certificates

# Копируем бинарник из builder stage
COPY --from=builder /app/main .

# Копируем статические файлы, миграции и т.д.
COPY ./migrations ./migrations
COPY ./config ./config

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]