package cmd

import (
	"os"

	"codex/internal/logging"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codex",
	Short: "A context-aware CLI assistant for your development environment",
	Long: `Codex is a CLI assistant that understands your entire development context.
It analyzes your Nix configuration, dotfiles, and current workspace to provide
personalized recommendations, keybind information, and tooling suggestions
tailored to your specific environment.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger based on verbose flag
		logging.InitWithVerbose(verbose)

		if verbose {
			logging.Logger.Debug().
				Str("command", cmd.Name()).
				Strs("args", args).
				Msg("Command execution started")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logging.Logger.Error().Err(err).Msg("Command execution failed")
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output for debugging")
}
