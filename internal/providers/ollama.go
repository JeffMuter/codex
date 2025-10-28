package providers

import (
	"context"
	"fmt"
	"io"

	codexContext "codex/internal/context"
)

// OllamaProvider implements the Provider interface for local Ollama
type OllamaProvider struct {
	baseURL string
	model   string
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(baseURL string) *OllamaProvider {
	return &OllamaProvider{
		baseURL: baseURL,
		model:   "llama2", // Default model
	}
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// SendQuery sends a query to Ollama
func (p *OllamaProvider) SendQuery(ctx context.Context, query string, context *codexContext.Context, writer io.Writer) error {
	// TODO: Implement actual API call
	// - Build prompt with context
	// - Make request to Ollama API
	// - Stream response to writer
	// - Handle errors
	return fmt.Errorf("not yet implemented")
}

// Validate checks if Ollama is accessible
func (p *OllamaProvider) Validate() error {
	// TODO: Ping Ollama to check if it's running
	// - Make a simple request to /api/tags or similar
	// - Check if the model is available
	return nil
}

// EstimateTokens estimates token count for a query
func (p *OllamaProvider) EstimateTokens(query string, context *codexContext.Context) (int, error) {
	// For local models, we can be more lenient with estimation
	// TODO: Implement basic token estimation
	return 0, fmt.Errorf("not yet implemented")
}

// GetCostEstimate returns estimated cost (always 0 for local)
func (p *OllamaProvider) GetCostEstimate(tokens int) (float64, error) {
	// Local models are free
	return 0.0, nil
}

// SetModel allows changing the model
func (p *OllamaProvider) SetModel(model string) {
	p.model = model
}

// ListModels returns available Ollama models
func (p *OllamaProvider) ListModels() ([]string, error) {
	// TODO: Query Ollama for available models
	// - GET /api/tags
	return nil, fmt.Errorf("not yet implemented")
}
