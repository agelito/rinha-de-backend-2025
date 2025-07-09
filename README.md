# rinha-de-backend-2025

This Go workspace demonstrates a distributed system using NATS for messaging, Protocol Buffers for serialization, and separate API and worker services.

## Operations

### 1. Generate Protobuf Files

Generate Go code from `.proto` files:

```sh
make proto
```

### 2. Run NATS

Start a local NATS server using Docker Compose:

```sh
docker compose up
```

### 3. Run API Service

Start the API service:

```sh
go run ./api
```

### 4. Run Payment Worker Service

Start the payment worker service:

```sh
go run ./payment-worker
```

## Benchmarking

To test the API and workers, run:

```sh
make benchmark
```

> **Note:** This requires the [`hey`](https://github.com/rakyll/hey) CLI tool to be installed.
