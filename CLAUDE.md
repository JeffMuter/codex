# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

This project uses Nix for dependency management. Enter the development shell:

```bash
nix-shell shell.nix
```

The environment includes:
- Go 1.23
- goose (database migration tool)
- claude-code (from nixos-unstable)

## Technology Stack

- **Language**: Go 1.23
- **Database Migrations**: goose

## Build and Run

Standard Go commands:

```bash
# Build
go build -o codex main.go

# Run
go run main.go

# Test
go test ./...

# Test with verbose output
go test -v ./...

# Test a specific package
go test ./path/to/package

# Run a specific test
go test -run TestName ./...
```

## Database Migrations

This project uses goose for database migrations:

```bash
# Create a new migration
goose create migration_name sql

# Run migrations
goose up

# Rollback last migration
goose down

# Check migration status
goose status
```
