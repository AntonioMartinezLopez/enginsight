# Enginsight

A Go-based client-server application using JSON-RPC for communication.

## Prerequisites

- Go 1.26.0 or higher
- Task (for task automation)
- Mockery (for generating mocks)

## Getting Started

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
- `server/internal/domain/` - Domain logic and business rules

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

## Possible implementation improvements

1. **Context Handling** - Context parameters accepted but not checked for cancellation/timeouts
2. **Error Propagation** - Custom error codes/types not exposed through JSON-RPC layer to clients
