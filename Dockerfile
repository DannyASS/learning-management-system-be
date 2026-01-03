# ----------------------
# Stage 1: Build
# ----------------------
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod & go.sum, download dependencies
COPY go.mod go.sum ./
RUN apk add --no-cache git ca-certificates \
    && go mod download

# Copy seluruh project
COPY . .

# Build binary (sesuaikan path main.go jika ada di root)
RUN go build -o main ./main.go

# ----------------------
# Stage 2: Runtime
# ----------------------
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .

# Install ca-certificates supaya HTTPS bisa jalan
RUN apk add --no-cache ca-certificates

# Port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
