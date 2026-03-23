# Build stage
FROM golang:alpine AS builder

# Set Go Proxy for speed in China
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build (CGO disabled, strip debug info, trim paths)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -v -o test-ebook-backend ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install CA certificates (for external APIs like PaddleOCR) and tzdata
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/test-ebook-backend .
COPY --from=builder /app/config.yaml .

# Create directories for persistence
RUN mkdir -p data uploads

EXPOSE 8182

CMD ["./test-ebook-backend"]
