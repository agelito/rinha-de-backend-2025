FROM golang:1.24-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY api/ ./api/

# Workspace dependencies
COPY messages/ ./messages/
COPY payment-worker/ ./payment-worker/
COPY processor-watchdog/ ./processor-watchdog/

WORKDIR /workspace/api

RUN go build -o /app/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api /app/api

EXPOSE 3001

CMD ["/app/api"]
