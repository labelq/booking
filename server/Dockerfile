# Используем официальный образ Go
FROM golang:1.24-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем curl для скачивания скрипта wait-for-it
RUN apk update && apk add --no-cache curl bash

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем все исходники в контейнер
COPY . .

# Копируем папку миграций в контейнер
COPY migrations /app/migrations

# Скачиваем скрипт wait-for-it.sh
RUN curl -o /usr/local/bin/wait-for-it https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh

# Даем скрипту права на выполнение
RUN chmod +x /usr/local/bin/wait-for-it

# Собираем приложение
RUN go build -o main .

# Устанавливаем права на выполнение
RUN chmod +x ./main

# Указываем порт
EXPOSE 8080

# Запускаем приложение, используя wait-for-it
CMD ["wait-for-it", "db:5432", "--", "./main"]