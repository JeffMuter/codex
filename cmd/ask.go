package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	screenshot bool
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a question with full context awareness",
	Long: `Ask a question and receive an answer based on your complete development environment.
Codex will automatically gather context from:
  - Your Nix configuration
  - Your dotfiles/home-manager setup
  - Current git repository (if available)
  - Current filesystem location
  - Screenshot (if --screenshot flag is used)

Examples:
  codex ask "What's my tmux prefix key?"
  codex ask --screenshot "How do I achieve this layout?"
  codex ask "What CLI tools do I have for JSON processing?"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := strings.Join(args, " ")

		if verbose {
			fmt.Println("Verbose mode enabled")
			fmt.Printf("Question: %s\n", question)
			fmt.Printf("Screenshot requested: %v\n", screenshot)
		}

		// TODO: Implement the actual ask logic
		// This will involve:
		// 1. Gathering context from all sources
		// 2. Building the complete prompt
		// 3. Sending to AI provider
		// 4. Streaming the response

		fmt.Printf("Processing question: %s\n", question)
		if screenshot {
			fmt.Println("Screenshot context will be captured")
		}

		// Placeholder response
		fmt.Println("\n[Not yet implemented - context gathering and AI provider integration coming soon]")
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Local flags for the ask command
	askCmd.Flags().BoolVarP(&screenshot, "screenshot", "s", false, "capture a screenshot for visual context")
}
