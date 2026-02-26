# Contributing

## Setting up the environment

```sh
go mod download
go build ./...
```

Requires [Go 1.22+](https://go.dev/doc/install).

## Project Structure

This is a hand-crafted SDK — no code generation. All source files are manually maintained.

Key files:
- `client.go` — Client + functional options
- `errors.go` — Typed errors
- `session.go`, `event.go`, `agent.go`, etc. — Service types and methods

See [REDESIGN.md](REDESIGN.md) for architecture decisions.

## Running Tests

```sh
go test ./...
```

Some tests may require a running OpenCode server or mock server.

## Formatting

```sh
gofmt -w .
```

## Submitting Changes

1. Fork the repo
2. Create a feature branch
3. Make your changes
4. Run tests and formatting
5. Open a PR with a clear description
