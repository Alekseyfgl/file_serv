# Используем официальный образ Go (версию подставьте нужную)
FROM golang:1.23

# Создадим рабочую директорию
WORKDIR /app

# Скопируем файлы go.mod и go.sum, чтобы установить зависимости
COPY go.mod go.sum ./

# Загрузим модули и зависимости
RUN go mod download

# Скопируем весь проект в /app
COPY . .

# Собираем бинарник из cmd/main.go
RUN go build -o server ./cmd/main.go

# Слушаем порт приложения
EXPOSE 3000

# Запускаем бинарник
CMD ["./server"]
