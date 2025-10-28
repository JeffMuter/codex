package providers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	codexContext "codex/internal/context"
)

// AnthropicProvider implements the Provider interface for Anthropic Claude
type AnthropicProvider struct {
	apiKey string
	model  string
	client *anthropic.Client
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &AnthropicProvider{
		apiKey: apiKey,
		model:  "claude-3-5-haiku-20241022", // Claude 3.5 Haiku - fast and efficient
		client: &client,
	}
}

// NewAnthropicProviderWithModel creates a new Anthropic provider with a specific model
func NewAnthropicProviderWithModel(apiKey string, model string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	if model == "" {
		model = "claude-3-5-haiku-20241022" // Default to Haiku
	}

	return &AnthropicProvider{
		apiKey: apiKey,
		model:  model,
		client: &client,
	}
}

// Name returns the provider name
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// SendQuery sends a query to Anthropic Claude API
func (p *AnthropicProvider) SendQuery(ctx context.Context, query string, contextData *codexContext.Context, writer io.Writer) error {
	// Calculate and print context memory usage
	memoryBytes := calculateContextMemory(contextData)
	fmt.Fprintf(writer, "Context Memory: %s\n\n", formatBytes(memoryBytes))

	// Build the prompt with context
	prompt := p.buildPrompt(query, contextData)

	// Create the message request
	stream := p.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	// Stream the response
	for stream.Next() {
		event := stream.Current()

		// Handle content block delta events - check the type string
		if event.Type == "content_block_delta" {
			// Access the Text field from the Delta
			if event.Delta.Text != "" {
				if _, err := writer.Write([]byte(event.Delta.Text)); err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	// Add newline at the end
	writer.Write([]byte("\n"))

	return nil
}

// buildPrompt constructs the prompt from query and context
func (p *AnthropicProvider) buildPrompt(query string, ctx *codexContext.Context) string {
	var sb strings.Builder

	sb.WriteString("You are Codex, a ruthlessly concise CLI assistant.\n\n")
	sb.WriteString("**ABSOLUTE RULES - VIOLATING THESE IS UNACCEPTABLE**:\n")
	sb.WriteString("1. ANSWER IN 5 WORDS OR LESS when possible.\n")
	sb.WriteString("2. NO introductions. NO explanations. NO context. NO examples. NO code blocks.\n")
	sb.WriteString("3. Format for lookups: `value` (file:line)\n")
	sb.WriteString("4. DO NOT explain what the user will do with the answer.\n")
	sb.WriteString("5. DO NOT restate the question.\n")
	sb.WriteString("6. If you write more than one sentence, you have FAILED.\n\n")

	// Add context sections
	if ctx != nil {
		if len(ctx.ConfiguredRepos) > 0 {
			sb.WriteString("## Configured Repositories\n\n")
			for _, repo := range ctx.ConfiguredRepos {
				sb.WriteString(fmt.Sprintf("### Repository: %s (%s)\n", repo.Source, repo.Type))
				if repo.Remote != "" {
					sb.WriteString(fmt.Sprintf("Remote: %s\n", repo.Remote))
				}
				sb.WriteString(fmt.Sprintf("Path: %s\n", repo.Path))

				// Include file contents if available
				if repo.Contents != nil {
					sb.WriteString(fmt.Sprintf("\n**Files: %d files, %d bytes total**\n\n",
						repo.Contents.TotalFiles, repo.Contents.TotalSize))

					for _, file := range repo.Contents.Files {
						sb.WriteString(fmt.Sprintf("#### File: %s\n", file.RelativePath))
						sb.WriteString("```\n")
						sb.WriteString(file.Content)
						if !strings.HasSuffix(file.Content, "\n") {
							sb.WriteString("\n")
						}
						sb.WriteString("```\n\n")
					}
				}
				sb.WriteString("\n")
			}
		}

		if ctx.CurrentRepo != nil {
			sb.WriteString("## Current Repository\n")
			sb.WriteString(fmt.Sprintf("Path: %s\n", ctx.CurrentRepo.Path))
			if ctx.CurrentRepo.Remote != "" {
				sb.WriteString(fmt.Sprintf("Remote: %s\n", ctx.CurrentRepo.Remote))
			}

			// Include file contents if available
			if ctx.CurrentRepo.Contents != nil {
				sb.WriteString(fmt.Sprintf("\n**Files: %d files, %d bytes total**\n\n",
					ctx.CurrentRepo.Contents.TotalFiles, ctx.CurrentRepo.Contents.TotalSize))

				for _, file := range ctx.CurrentRepo.Contents.Files {
					sb.WriteString(fmt.Sprintf("#### File: %s\n", file.RelativePath))
					sb.WriteString("```\n")
					sb.WriteString(file.Content)
					if !strings.HasSuffix(file.Content, "\n") {
						sb.WriteString("\n")
					}
					sb.WriteString("```\n\n")
				}
			}
			sb.WriteString("\n")
		}

		if ctx.Filesystem != nil {
			sb.WriteString("## Filesystem Context\n")
			sb.WriteString(fmt.Sprintf("Current Directory: %s\n", ctx.Filesystem.CurrentDir))
			sb.WriteString("\n")
		}
	}

	// Add the user's query
	sb.WriteString("## User Query\n")
	sb.WriteString(query)
	sb.WriteString("\n")

	return sb.String()
}

// Validate checks if the provider is properly configured
func (p *AnthropicProvider) Validate() error {
	if p.apiKey == "" {
		return fmt.Errorf("anthropic API key is required")
	}
	return nil
}

// EstimateTokens estimates token count for a query
func (p *AnthropicProvider) EstimateTokens(query string, context *codexContext.Context) (int, error) {
	// Rough estimation: ~4 characters per token
	prompt := p.buildPrompt(query, context)
	estimatedTokens := len(prompt) / 4
	return estimatedTokens, nil
}

// GetCostEstimate returns estimated cost in USD
func (p *AnthropicProvider) GetCostEstimate(tokens int) (float64, error) {
	// Claude 3.5 Sonnet pricing (as of late 2024):
	// Input: $3 per million tokens
	// Output: $15 per million tokens
	// For estimation, assume 1:1 input/output ratio
	const avgCostPerToken = (3.0 + 15.0) / 2.0 / 1_000_000
	return float64(tokens) * avgCostPerToken, nil
}

// SetModel allows changing the model
func (p *AnthropicProvider) SetModel(model string) {
	p.model = model
}

// calculateContextMemory calculates the total memory usage of the context data
func calculateContextMemory(ctx *codexContext.Context) int64 {
	if ctx == nil {
		return 0
	}

	var totalBytes int64

	// Calculate configured repos memory
	for _, repo := range ctx.ConfiguredRepos {
		totalBytes += calculateRepoMemory(repo)
	}

	// Calculate current repo memory
	if ctx.CurrentRepo != nil {
		totalBytes += calculateRepoMemory(ctx.CurrentRepo)
	}

	// Calculate filesystem context
	if ctx.Filesystem != nil {
		totalBytes += int64(len(ctx.Filesystem.CurrentDir))
		for _, dir := range ctx.Filesystem.ParentDirs {
			totalBytes += int64(len(dir))
		}
		for _, file := range ctx.Filesystem.Files {
			totalBytes += int64(len(file))
		}
	}

	// Calculate screenshot memory
	if ctx.Screenshot != nil {
		totalBytes += int64(len(ctx.Screenshot.Path))
		totalBytes += int64(len(ctx.Screenshot.Data))
		totalBytes += int64(len(ctx.Screenshot.MimeType))
		totalBytes += int64(len(ctx.Screenshot.Tool))
	}

	return totalBytes
}

// calculateRepoMemory calculates memory usage for a single repository context
func calculateRepoMemory(repo *codexContext.RepositoryContext) int64 {
	if repo == nil {
		return 0
	}

	var totalBytes int64

	// Basic strings
	totalBytes += int64(len(repo.Path))
	totalBytes += int64(len(repo.Remote))
	totalBytes += int64(len(repo.Source))
	totalBytes += int64(len(repo.Type))

	// Calculate contents memory
	if repo.Contents != nil {
		for _, file := range repo.Contents.Files {
			totalBytes += int64(len(file.Path))
			totalBytes += int64(len(file.RelativePath))
			totalBytes += int64(len(file.Content))
			totalBytes += int64(file.Size) // This counts actual file size
		}
	}

	return totalBytes
}

// formatBytes formats a byte count into a human-readable string
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB (%d bytes)", float64(bytes)/float64(GB), bytes)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB (%d bytes)", float64(bytes)/float64(MB), bytes)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB (%d bytes)", float64(bytes)/float64(KB), bytes)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
