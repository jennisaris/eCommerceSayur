# ==============================================================================
# Multi-Stage Dockerfile for High-Performance Golang eCommerce Backend
# ==============================================================================

# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod ./
# Note: go.sum will be downloaded via go mod tidy
RUN go mod download || true

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl
ENV TZ=Asia/Jakarta

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/api"]
