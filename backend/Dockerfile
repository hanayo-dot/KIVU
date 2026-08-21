# Multi-stage Dockerfile for KIVU Go Backend

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Final production stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server ./server

EXPOSE 8080

CMD ["./server"]
