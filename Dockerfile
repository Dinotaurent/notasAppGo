# Imagen base de go
FROM golang:tip-alpine3.23 AS builder
LABEL authors="dinotaurent"

RUN mkdir /app

COPY . /app

# Directorio de trabajo
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o notasApp ./cmd/api

# Opcional
RUN chmod +x /app/notasApp

# Construir la imagen
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/notasApp /app/notasApp

CMD ["/app/notasApp"]