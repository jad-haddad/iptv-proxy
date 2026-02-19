# AGENTS
#
# Guidance for agentic coding tools working in this repo.
# Keep this file up to date when build/lint/test or style norms change.

## Project snapshot

- Language: Go (module: github.com/jad-haddad/iptv-proxy)
- Go version: 1.22 (see go.mod)
- Entry point: cmd/iptv-proxy
- Runtime: HTTP server serving /lebanon.m3u, /epg.xml, /health
- Primary packages: internal/httpserver, internal/m3u, internal/epg

## Build, run, lint, test

### Local run

- Run server (dev): `go run ./cmd/iptv-proxy`
- Healthcheck command: `go run ./cmd/iptv-proxy healthcheck`

### Build

- Build binary: `go build -o ./iptv-proxy ./cmd/iptv-proxy`
- Docker build: `docker build -t iptv-proxy .`

### Lint / static checks

- Basic vet: `go vet ./...`
- Format all Go code: `gofmt -w ./cmd ./internal ./scripts`

Notes:
- There is no repo-specific linter configured (no golangci-lint config found).
- Prefer gofmt over gofmt -w on the whole repo if you need to minimize diffs.

### Tests

- All tests: `go test ./...`
- Single package tests: `go test ./internal/m3u`
- Single test by name (example): `go test ./internal/m3u -run TestFilter -count=1`
- Run subtests: `go test ./internal/epg -run TestFilter/CaseName -count=1`

### Integration check

- Endpoint validation script: `go run ./scripts/check_endpoints.go http://127.0.0.1:8080`
- What it checks: 200 OK responses and ETag 304 behavior for M3U/EPG.

## Configuration

Environment variables (with defaults):
- `M3U_URL`: https://iptv-org.github.io/iptv/countries/lb.m3u
- `EPG_URL`: https://mdag9904.github.io/lebanon-epg/epg.xml
- `MTV_REGEX`: `(?i)\bmtv\b.*\blebanon\b|mtv\s*lebanon|mtvlebanon`
- `MTV_TVG_ID`: `mtvlebanon.lb`
- `MTV_TVG_NAME`: `MTV Lebanon`
- `EPG_REFRESH_SECONDS`: `3600`
- `REQUEST_TIMEOUT_SECONDS`: `15`

## Code style and conventions

### Formatting

- Always use gofmt. The codebase follows standard Go formatting.
- Keep lines readable; avoid overly long inline expressions.

### Imports

- Group imports with standard Go ordering:
  1) standard library
  2) third-party
  3) module-local (`github.com/jad-haddad/iptv-proxy/...`)
- Use explicit package names; avoid dot imports.

### Package structure

- `cmd/iptv-proxy` contains the main entry point and CLI behavior.
- `internal/` contains all application logic (httpserver, m3u, epg, cache, config).
- `scripts/` hosts helper programs executed via `go run`.

### Naming

- Packages: short, lowercase, no underscores (e.g., `httpserver`, `m3u`).
- Types: PascalCase (`Config`, `Server`, `M3UCache`).
- Functions and variables: camelCase; exported only when needed.
- Constants: CamelCase; keep acronyms consistently cased (e.g., `EPG`, `M3U`).

### Types and data flow

- Prefer concrete types over interface{}.
- Keep structs small and cohesive (see `cache` structs).
- Prefer `time.Duration` for timeouts; convert from seconds at the boundary.
- Pass config via a `Config` struct instead of global vars.

### Error handling

- Return errors up the stack when possible; handle at the boundary.
- For HTTP handlers, use `http.Error` with a human-friendly message and status.
- For CLI/entry points, `log.Fatal` is acceptable for unrecoverable errors.
- Avoid panics in normal control flow.

### Concurrency and state

- Protect shared state with a mutex (see `cache` usage in `httpserver`).
- Prefer short critical sections; lock only around state updates.
- Avoid global mutable state when you can scope it to a server instance.

### HTTP behavior

- Always set `Content-Type` for responses (`application/xml`, `application/json`).
- Use strong ETags and honor `If-None-Match` for 304 responses.
- When upstream fetch fails, return `502 Bad Gateway` with a stable message.

### Parsing and data transforms

- Keep parsing logic isolated in `internal/m3u` and `internal/epg`.
- Favor deterministic output: stable ordering and predictable formatting.
- Avoid modifying inputs in place; return new values instead.

### Logging

- Use `log` in main only; avoid noisy logging inside request handlers.
- Keep logs actionable and minimal.

### Files and edits

- Avoid introducing new dependencies unless necessary.
- Keep `scripts/` as small, self-contained Go programs.
- If adding tests, place them alongside the package (`*_test.go`).

## Cursor/Copilot rules

- No Cursor rules found in `.cursor/rules/` or `.cursorrules`.
- No Copilot instructions found in `.github/copilot-instructions.md`.

## Tips for agents

- Prefer small, focused diffs; do not reformat unrelated code.
- Be careful with upstream URLs and regex defaults; they are user-configurable.
- Preserve ETag behavior; it is part of the integration contract.
- When editing handlers, verify `check_endpoints.go` still passes.
