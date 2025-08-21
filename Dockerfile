FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Instala swag CLI para gerar docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN go mod download

# Gera documentação Swagger (ajuste o caminho do main.go se necessário)
RUN swag init -g cmd/api/main.go -o docs

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 go build -o /mecanica-xpto-api

FROM alpine:3.22

WORKDIR /app

COPY .env-example .

COPY --from=builder /mecanica-xpto-api .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./mecanica-xpto-api"]
