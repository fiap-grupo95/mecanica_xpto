FROM golang:1.24-alpine as builder

WORKDIR /app

COPY .. .

RUN go mod download

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 go build -o /mecanica-xpto-api


FROM alpine:latest

WORKDIR /app

COPY --from=builder /mecanica-xpto-api .

EXPOSE 8080

CMD ["./mecanica-xpto-api"]