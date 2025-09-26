FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем модули для кэширования
COPY go.mod ./
#COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение из cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w" \
    -o /go/bin/server ./cmd/server

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем бинарник
COPY --from=builder /go/bin/server .

# Создаем пользователя
RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 8080
CMD ["./server"]
