FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 go build -o /mecanica-xpto-api


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /mecanica-xpto-api .

EXPOSE 8080

CMD ["./mecanica-xpto-api"]