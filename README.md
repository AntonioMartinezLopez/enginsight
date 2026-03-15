# Enginsight

A Go-based client-server application using JSON-RPC for communication.

## Prerequisites

- Go 1.26.0 or higher
- Task (for task automation) - https://taskfile.dev/
- Mockery (for generating mocks) - https://vektra.github.io/mockery/latest/

## Getting Started

Important: Server applicaiton uses environment variables for setting host, port and rpc path. Please copy the existing `.env.template` file to provide those variables:

```bash
cp server/.env.template server/.env
```

### Running the Server

```bash
task run:server
```

### Available Tasks

View all available tasks:
```bash
task
```

Common tasks:
- `task lint` - Run linters
- `task fmt` - Format code
- `task run:server` - Start the server

## Project Structure

- `client/` - Client implementation
- `server/` - Server implementation with domain-driven design
- `jrpc/` - JSON-RPC communication layer
- `server/internal/domain/` - Domain/Business logic

## Testing

### Run Unit Tests

```bash
task test:server:unit
```

### Run Integration Tests

```bash
task test:server:integration
```

### Generate Mocks

```bash
task generate:mocks:server
```

## Current Security Issues

1. **Unlimited Input Size** - The server accepts messages of any size without validation (risk of memory exhaustion, blocking, server crash etc.)
2. **No HTTP Request Body Size Limits**
3. **No Rate Limiting**

## Possible implementation improvements

1. **Context Handling** - Context parameters accepted but not checked for cancellation/timeouts
2. **Error Propagation** - Custom error codes/types not exposed through JSON-RPC layer to clients
3. **Testing** - Mainly testing happy paths
