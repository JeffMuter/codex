# Logging in Codex

Codex uses [zerolog](https://github.com/rs/zerolog) for structured, high-performance logging.

## Features

- **Zero-allocation**: Minimal memory footprint for CLI performance
- **Structured logging**: JSON output with type-safe fields
- **Pretty console output**: Human-readable format in verbose mode
- **Log levels**: Debug, Info, Warn, Error

## Usage in Code

### Basic Logging

```go
import "codex/internal/logging"

// Info level
logging.Logger.Info().Msg("Operation completed")

// Warning
logging.Logger.Warn().Msg("Potential issue detected")

// Error
logging.Logger.Error().Err(err).Msg("Operation failed")
```

### Structured Fields

```go
logging.Logger.Debug().
    Str("key", "value").
    Int("count", 42).
    Bool("success", true).
    Msg("Processing completed")
```

### With Context

```go
contextLogger := logging.WithContext(map[string]interface{}{
    "user_id": "123",
    "request_id": "abc",
})
contextLogger.Info().Msg("User action performed")
```

## CLI Usage

### Normal Mode (INFO level, JSON output)
```bash
codex ask "What's my tmux prefix?"
# Output: User-facing messages only
```

### Verbose Mode (DEBUG level, Pretty output)
```bash
codex -v ask "What's my tmux prefix?"
# Output: User messages + colored debug logs
# Example:
# 2025-10-28T15:13:22Z DBG Processing ask command question="What's my tmux prefix?" screenshot=false
```

## Log Levels

- **Debug**: Detailed information for debugging (only in verbose mode)
- **Info**: General informational messages
- **Warn**: Warning messages for potential issues
- **Error**: Error messages for failures

## Output Format

### JSON (default)
```json
{"level":"debug","time":"2025-10-28T15:13:22Z","message":"Processing ask command","question":"test","screenshot":false}
```

### Pretty Console (verbose mode)
```
2025-10-28T15:13:22Z DBG Processing ask command question="test" screenshot=false
2025-10-28T15:13:22Z WRN Ask command not yet fully implemented
```

## Integration Points

The logger is initialized in `cmd/root.go` via `PersistentPreRun`, which runs before all commands:

```go
PersistentPreRun: func(cmd *cobra.Command, args []string) {
    // Initialize logger based on verbose flag
    logging.InitWithVerbose(verbose)
}
```

This ensures all commands have access to properly configured logging.

## Best Practices

1. **Use structured fields**: Always use `.Str()`, `.Int()`, etc. for values
2. **Meaningful messages**: Make the `Msg()` descriptive
3. **Appropriate levels**:
   - Debug: Internal state, detailed flow
   - Info: Successful operations, milestones
   - Warn: Recoverable issues, deprecations
   - Error: Failures, exceptions
4. **Don't log secrets**: Never log API keys, tokens, or sensitive data
5. **Context is king**: Add relevant context fields (command, flags, paths)

## Performance

Zerolog is the fastest Go logging library:
- Zero allocation in hot paths
- Minimal CPU overhead
- Perfect for CLI tools where startup time matters
