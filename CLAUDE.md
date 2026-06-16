# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`ai-doc-platform` is a microservices-based AI documentation platform. Currently contains one service:

- **gateway-service** — Go HTTP gateway using [chi](https://github.com/go-chi/chi) router, listens on `:8080`

## Commands

All commands are run from within the service directory (e.g., `gateway-service/`).

```bash
# Run the server
go run ./cmd/server/main.go

# Build
go build ./...

# Test
go test ./...

# Run a single test
go test ./path/to/package -run TestFunctionName

# Lint (if golangci-lint is installed)
golangci-lint run
```

## Architecture

Each service follows standard Go project layout:
- `cmd/server/main.go` — entrypoint: router setup and server start
- `internal/` — private packages (config, handlers, etc.)

The module path is `github.com/vokhotin/ai-doc-platform/<service-name>`. Each service has its own `go.mod` — they are independent Go modules.

Router: `chi/v5` is the HTTP router. Register routes in `main.go` for now; move handlers to `internal/` as the service grows.
