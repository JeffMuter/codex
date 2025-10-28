package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"codex/internal/config"
	codexContext "codex/internal/context"
	"codex/internal/logging"
	"codex/internal/providers"

	"github.com/spf13/cobra"
)

var (
	screenshot  bool
	currentRepo bool
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a question with full context awareness",
	Long: `Ask a question and receive an answer based on your complete development environment.
Codex will automatically gather context from:
  - Your configured repositories (always included)
  - Your Nix configuration (if configured)
  - Your dotfiles/home-manager setup (if configured)
  - Current git repository (opt-in with --current-repo flag)
  - Current filesystem location
  - Screenshot (if --screenshot flag is used)

Examples:
  codex ask "What's my tmux prefix key?"
  codex ask --screenshot "How do I achieve this layout?"
  codex ask --current-repo "Explain this codebase structure"
  codex ask "What CLI tools do I have for JSON processing?"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")

		logging.Logger.Debug().
			Str("question", question).
			Bool("screenshot", screenshot).
			Bool("current_repo", currentRepo).
			Msg("Processing ask command")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		logging.Logger.Debug().
			Str("provider", cfg.Provider).
			Int("configured_repos", len(cfg.ConfiguredRepos)).
			Msg("Configuration loaded")

		// Create AI provider
		var provider providers.Provider
		switch cfg.Provider {
		case "anthropic":
			if cfg.Model != "" {
				provider = providers.NewAnthropicProviderWithModel(cfg.AnthropicKey, cfg.Model)
			} else {
				provider = providers.NewAnthropicProvider(cfg.AnthropicKey)
			}
		default:
			return fmt.Errorf("unsupported provider: %s", cfg.Provider)
		}

		// Validate provider
		if err := provider.Validate(); err != nil {
			return fmt.Errorf("provider validation failed: %w", err)
		}

		logging.Logger.Debug().Str("provider", provider.Name()).Msg("Provider initialized")

		// Gather context
		gatherer := codexContext.NewGatherer(cfg)

		workingDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		gatherOpts := codexContext.GatherOptions{
			IncludeCurrentRepo: currentRepo,
			IncludeFilesystem:  true,
			IncludeNixConfig:   cfg.NixConfigPath != "",
			IncludeDotfiles:    cfg.DotfilesPath != "",
			CaptureScreenshot:  screenshot,
			WorkingDir:         workingDir,
		}

		logging.Logger.Debug().Msg("Gathering context...")
		ctx, err := gatherer.Gather(gatherOpts)
		if err != nil {
			return fmt.Errorf("failed to gather context: %w", err)
		}

		logging.Logger.Debug().
			Int("configured_repos", len(ctx.ConfiguredRepos)).
			Bool("has_current_repo", ctx.CurrentRepo != nil).
			Msg("Context gathered")

		// Send query to provider
		logging.Logger.Debug().Msg("Sending query to provider...")

		apiCtx := context.Background()
		if err := provider.SendQuery(apiCtx, question, ctx, os.Stdout); err != nil {
			return fmt.Errorf("failed to get response: %w", err)
		}

		logging.Logger.Debug().Msg("Query completed successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Local flags for the ask command
	askCmd.Flags().BoolVarP(&screenshot, "screenshot", "s", false, "capture a screenshot for visual context")
	askCmd.Flags().BoolVarP(&currentRepo, "current-repo", "r", false, "include current working directory repository as context")
}
