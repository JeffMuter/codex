package cmd

import (
	"fmt"
	"strings"

	"codex/internal/logging"

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
	Run: func(cmd *cobra.Command, args []string) {
		question := strings.Join(args, " ")

		logging.Logger.Debug().
			Str("question", question).
			Bool("screenshot", screenshot).
			Bool("current_repo", currentRepo).
			Msg("Processing ask command")

		// TODO: Implement the actual ask logic
		// This will involve:
		// 1. Gathering context from all sources (configured repos + optional current repo)
		// 2. Building the complete prompt
		// 3. Sending to AI provider
		// 4. Streaming the response

		fmt.Printf("Processing question: %s\n", question)
		if screenshot {
			logging.Logger.Debug().Msg("Screenshot context will be captured")
			fmt.Println("Screenshot context will be captured")
		}
		if currentRepo {
			logging.Logger.Debug().Msg("Current repository will be included as context")
			fmt.Println("Current repository will be included as context")
		}

		// Placeholder response
		logging.Logger.Warn().Msg("Ask command not yet fully implemented")
		fmt.Println("\n[Not yet implemented - context gathering and AI provider integration coming soon]")
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Local flags for the ask command
	askCmd.Flags().BoolVarP(&screenshot, "screenshot", "s", false, "capture a screenshot for visual context")
	askCmd.Flags().BoolVarP(&currentRepo, "current-repo", "r", false, "include current working directory repository as context")
}
