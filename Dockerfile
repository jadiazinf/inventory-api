# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

# Run Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY .env.example .env

EXPOSE 3000

CMD ["./main"]
