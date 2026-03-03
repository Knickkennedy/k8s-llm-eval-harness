# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/eval-runner ./cmd/eval-runner

# Runtime stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/eval-runner .

# Default benchmark dir — overridden by ConfigMap volume mount in Kubernetes
COPY evaluations/benchmarks ./evaluations/benchmarks

ENTRYPOINT ["./eval-runner"]
