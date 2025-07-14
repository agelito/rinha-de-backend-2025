FROM golang:1.24-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY payment-worker/ ./payment-worker/

# Workspace dependencies
COPY messages/ ./messages/
COPY api/ ./api/
COPY processor-watchdog/ ./processor-watchdog/

WORKDIR /workspace/payment-worker

RUN go build -o /app/payment-worker

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment-worker /app/payment-worker

CMD ["/app/payment-worker"]

