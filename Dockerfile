FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mecanica-xpto-api .


FROM alpine:3.22

WORKDIR /app

# Copy the .env file
COPY .env .

COPY --from=builder /app/cmd/api/mecanica-xpto-api .

EXPOSE 8080

CMD ["./mecanica-xpto-api"]