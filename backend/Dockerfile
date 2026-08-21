# Multi-stage Dockerfile for KIVU Unified Application (Frontend + Backend)

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache git

# Cache Go dependencies
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

# Copy backend source code
COPY backend/ /app/backend/

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Final production stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy compiled backend binary
COPY --from=builder /app/server ./server

# Copy static frontend web assets
COPY frontend/ ./frontend/

EXPOSE 8080

CMD ["./server"]
