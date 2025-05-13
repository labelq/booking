# Используем официальный образ Node.js для сборки приложения
FROM node:18 AS build

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем package.json и package-lock.json (если есть)
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем весь проект в контейнер
COPY . .

# Строим проект для продакшн
RUN npm run build

# Используем nginx для сервирования статических файлов
FROM nginx:alpine

# Копируем собранные файлы в папку nginx
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]