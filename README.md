# allseer

Local programmable proxy in Go.

## Current Status

Project initialization is complete with a runnable application skeleton:

- Go module setup
- Application boot and graceful shutdown flow
- Config loading and validation
- Proxy server entrypoint with HTTP and CONNECT interception stubs
- Initial rule engine contracts and matcher scaffolding
- Starter config and rule files

## Quick Start

1. Run the proxy:

```bash
go run ./cmd/allseer
```

2. Optional: provide a custom config file:

```bash
ALLSEER_CONFIG=config/config.json go run ./cmd/allseer
```

By default, the proxy listens on `127.0.0.1:8080`.

## Project Layout

- `cmd/allseer/main.go`: process entrypoint and lifecycle wiring
- `internal/app`: application container and server startup/shutdown
- `internal/config`: config schema, parsing, and validation
- `internal/proxy`: proxy server and interception stubs
- `internal/rules`: rule model and evaluation engine scaffold
- `config`: starter runtime config and sample rules

## Next Build Step

Implement forward proxying and CONNECT tunnel handling in `internal/proxy`.

