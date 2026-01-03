# Gunakan Go 1.24 alpine (compatible dengan go.mod)
FROM golang:1.24-alpine

WORKDIR /app
COPY . .

# Install git & ca-certificates jika repo pakai go get HTTPS
RUN apk add --no-cache git ca-certificates && go mod tidy
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
