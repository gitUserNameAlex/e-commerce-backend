# Dockerfile
# Используем официальный образ Golang для сборки
FROM golang:1.21 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем исходный код в рабочую директорию
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Используем образ alpine для запуска
FROM alpine:latest

# Копируем бинарник из билдера
COPY --from=builder /app/main .

# Запускаем бинарник
CMD ["./main"]