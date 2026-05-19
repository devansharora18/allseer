# allseer

Local programmable proxy in Go.

## Current Status

Current implementation includes a working baseline proxy:

- Go module setup
- Application boot and graceful shutdown flow
- Config loading and validation
- HTTP forward proxy support
- HTTPS CONNECT tunneling support
- Rule engine evaluation for request decisions
- Rule actions: allow, block, redirect (HTTP), and request-header mutation
- Domain-based ad/tracker blocking from dedicated blocklist file
- Persistent traffic logging to local SQLite database
- Local admin endpoint for recent traffic logs
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

Rules are loaded from `config/rules.example.yaml` by default, and ad/tracker domains are loaded from `config/ad_domains.txt`.
Override these through `config/config.json`.

Traffic logs are persisted by default in `allseer.db` (configured via `database.path` in `config/config.json`).

Query recent logs from the local admin endpoint:

```bash
curl "http://127.0.0.1:8080/admin/logs?limit=100"
```

Notes:
- Endpoint is local-control API intended for direct requests to the proxy process.
- Absolute proxy requests are not treated as admin API calls.
- `limit` defaults to 100 and is capped at 1000.

## Project Layout

- `cmd/allseer/main.go`: process entrypoint and lifecycle wiring
- `internal/app`: application container and server startup/shutdown
- `internal/config`: config schema, parsing, and validation
- `internal/proxy`: forwarding/tunneling engine and rule application
- `internal/rules`: rule model, evaluation engine, and file loader
- `config`: starter runtime config and sample rules

## Next Build Step

Implement response interception transforms and hot-reload for runtime rules.

