# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the binary
COPY . .
RUN go build -o oneart-product-cart-service ./cmd/main.go

# Minimal image for production
FROM alpine:3.17
WORKDIR /app

COPY --from=build /app/oneart-product-cart-service /app/oneart-product-cart-service

EXPOSE 8080

# CMD without GOOGLE_APPLICATION_CREDENTIALS since Cloud Run handles it
CMD ["/app/oneart-product-cart-service"]
