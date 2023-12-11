# Используем официальный образ Golang
FROM golang:latest

# Установка необходимых зависимостей
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/nats-io/stan.go
RUN go get -u github.com/lib/pq
RUN go get -u github.com/patrickmn/go-cache

# Копирование исходного кода приложения в контейнер
COPY . /app

# Установка рабочей директории
WORKDIR /app

# Сборка приложения
RUN go build -o main .

# Определение команды для запуска приложения
CMD ["/app/main"]
