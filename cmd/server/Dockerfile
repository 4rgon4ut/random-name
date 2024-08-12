# Builer image
FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

# Stage 2: Create a lightweight image for running the application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

ENTRYPOINT ["./server"]
