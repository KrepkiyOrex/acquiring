# 1. Используем базовый образ для сборки
FROM golang:1.23 AS builder

# 2. Устанавливаем рабочую директорию
WORKDIR /app

# 3. Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# 4. Копируем весь проект в контейнер
COPY . ./

# 6. Сборка приложения
RUN go build -o app ./cmd

# 7. Создаем минимальный образ для запуска
FROM debian:bookworm-slim AS base

# 8. Копируем бинарник в финальный образ
COPY --from=builder /app/app /app

# 9. Указываем команду запуска
CMD ["./app"]
