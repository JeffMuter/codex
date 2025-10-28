package providers

import (
	"context"
	"fmt"
	"io"

	codexContext "codex/internal/context"
)

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey string
	model  string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: apiKey,
		model:  "gpt-4-turbo-preview", // Default model
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// SendQuery sends a query to OpenAI API
func (p *OpenAIProvider) SendQuery(ctx context.Context, query string, context *codexContext.Context, writer io.Writer) error {
	// TODO: Implement actual API call
	// - Build messages with context
	// - Make API request using OpenAI SDK or HTTP client
	// - Stream response to writer
	// - Handle errors and retries
	return fmt.Errorf("not yet implemented")
}

// Validate checks if the provider is properly configured
func (p *OpenAIProvider) Validate() error {
	if p.apiKey == "" {
		return fmt.Errorf("openai API key is required")
	}
	// TODO: Optionally test API key with a test request
	return nil
}

// EstimateTokens estimates token count for a query
func (p *OpenAIProvider) EstimateTokens(query string, context *codexContext.Context) (int, error) {
	// TODO: Implement token estimation
	// - Use tiktoken or approximation
	// - Count tokens in query and context
	return 0, fmt.Errorf("not yet implemented")
}

// GetCostEstimate returns estimated cost in USD
func (p *OpenAIProvider) GetCostEstimate(tokens int) (float64, error) {
	// GPT-4 Turbo pricing (as of late 2024):
	// Input: $10 per million tokens
	// Output: $30 per million tokens
	// For estimation, assume 1:1 input/output ratio
	const avgCostPerToken = (10.0 + 30.0) / 2.0 / 1_000_000
	return float64(tokens) * avgCostPerToken, nil
}

// SetModel allows changing the model
func (p *OpenAIProvider) SetModel(model string) {
	p.model = model
}
