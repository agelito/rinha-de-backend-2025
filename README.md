# rinha-de-backend-2025

This Go workspace demonstrates a distributed system using NATS for messaging, Protocol Buffers for serialization, and separate API and worker services.

## Development

Creating a new Go module in the workspace.

### 1. Create module directory

```sh
mkdir my-module-name
cd my-module-name
```

### 2. Initialize module

> In the module directory `./my-module-name`

```sh
go mod init github.com/agelito/rinha-de-backend-2025/my-module-name
```

### 3. Add module to workspace

> In the workspace directory

```sh
go work use ./my-module-name
```

### 4. Build & Run

> Create a main entrypoint in the `my-module-name` module and try running it from the workspace root.

```sh
go run ./my-module-name
```

## Running Development

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
