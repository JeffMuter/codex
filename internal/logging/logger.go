package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Global logger instance
var Logger zerolog.Logger

// LogLevel represents the logging level
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Config holds logging configuration
type Config struct {
	Level      LogLevel
	Pretty     bool // Use console writer for human-readable output
	TimeFormat string
}

// DefaultConfig returns default logging configuration
func DefaultConfig() Config {
	return Config{
		Level:      LevelInfo,
		Pretty:     false,
		TimeFormat: time.RFC3339,
	}
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) {
	// Set global log level
	switch cfg.Level {
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Configure time format
	zerolog.TimeFieldFormat = cfg.TimeFormat

	var output io.Writer = os.Stderr

	// Use console writer for pretty output (development/verbose mode)
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
	}

	// Create logger
	Logger = zerolog.New(output).With().Timestamp().Logger()

	// Set as global logger
	log.Logger = Logger
}

// InitWithVerbose is a convenience function for CLI verbose flag
// verbose=true enables debug level + pretty printing
func InitWithVerbose(verbose bool) {
	cfg := DefaultConfig()
	if verbose {
		cfg.Level = LevelDebug
		cfg.Pretty = true
	}
	Init(cfg)
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
	return &Logger
}

// WithContext returns a logger with additional context fields
func WithContext(fields map[string]interface{}) zerolog.Logger {
	ctx := Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx.Logger()
}
