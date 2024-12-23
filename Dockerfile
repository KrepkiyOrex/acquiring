# 1. Используем базовый образ для сборки
FROM golang:1.23 AS builder

# 2. Устанавливаем рабочую директорию
WORKDIR /app

# 3. Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# 4. Устанавливаем air для автоматической перезагрузки
RUN go install github.com/air-verse/air@latest

# 5. Копируем весь проект в контейнер
COPY . ./

# 6. Указываем команду для запуска приложения с air
CMD ["air"]
