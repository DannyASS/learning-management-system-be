# ----------------------
# Stage 1: Build
# ----------------------
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache git ca-certificates && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./main.go

# ----------------------
# Stage 2: Runtime
# ----------------------
FROM alpine:latest
WORKDIR /app

# Install tini
RUN apk add --no-cache tini ca-certificates

COPY --from=builder /server /app/server
COPY --from=builder /app/internal/resources ./internal/resources

EXPOSE 8080

# Prefork-safe: tini menangani PID 1 & forward signal
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/server"]
