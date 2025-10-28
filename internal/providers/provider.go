package providers

import (
	"context"
	"fmt"
	"io"

	codexContext "codex/internal/context"
)

// Provider defines the interface for AI providers
type Provider interface {
	// Name returns the provider name
	Name() string

	// SendQuery sends a query with context and streams the response
	SendQuery(ctx context.Context, query string, context *codexContext.Context, writer io.Writer) error

	// Validate checks if the provider is properly configured
	Validate() error

	// EstimateTokens estimates token count for a query
	EstimateTokens(query string, context *codexContext.Context) (int, error)

	// GetCostEstimate returns estimated cost in USD for a query
	GetCostEstimate(tokens int) (float64, error)
}

// Config holds provider configuration
type Config struct {
	Provider     string
	AnthropicKey string
	OpenAIKey    string
	OllamaURL    string
}

// NewProvider creates a provider based on configuration
func NewProvider(cfg *Config) (Provider, error) {
	switch cfg.Provider {
	case "anthropic":
		if cfg.AnthropicKey == "" {
			return nil, fmt.Errorf("anthropic provider requires API key")
		}
		return NewAnthropicProvider(cfg.AnthropicKey), nil

	case "openai":
		if cfg.OpenAIKey == "" {
			return nil, fmt.Errorf("openai provider requires API key")
		}
		return NewOpenAIProvider(cfg.OpenAIKey), nil

	case "ollama":
		if cfg.OllamaURL == "" {
			cfg.OllamaURL = "http://localhost:11434"
		}
		return NewOllamaProvider(cfg.OllamaURL), nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

// Response represents a streaming response chunk
type Response struct {
	Content string
	Done    bool
	Error   error
}
