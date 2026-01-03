# Gunakan Go 1.24 alpine (sesuai go.mod)
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy file project
COPY go.mod go.sum ./
RUN apk add --no-cache git ca-certificates \
    && go mod download

COPY . .

# Build binary
RUN go build -o main ./cmd/main.go   # sesuaikan path main.go

# ----------------------
# Stage runtime (lebih ringan)
# ----------------------
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .
# Copy .env jika mau masuk image (opsional, lebih aman mount di docker run)
# COPY --from=builder /app/.env .

# Install ca-certificates supaya HTTPS bisa jalan
RUN apk add --no-cache ca-certificates

# Port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
