# Use Go 1.24.2 to satisfy go.mod requirement
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authservice ./cmd/main.go

# --- Runtime Stage ---
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

COPY --from=builder /app/authservice .
COPY start.sh .
COPY config ./config
COPY migrations ./migrations

# 👇 Fix line endings and permissions
RUN sed -i 's/\r$//' start.sh && chmod +x start.sh authservice

ENTRYPOINT ["/app/start.sh"]


   
    