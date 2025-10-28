package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage codex configuration",
	Long: `Configure codex settings including paths to your Nix configuration
and dotfiles repositories. Configuration can be set via environment
variables, config file, or using these commands.`,
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display the current codex configuration including all paths and settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Verbose mode enabled")
		}

		fmt.Println("Codex Configuration")
		fmt.Println("===================")
		fmt.Println()

		// Check environment variables
		nixConfig := os.Getenv("CODEX_NIX_CONFIG")
		dotfiles := os.Getenv("CODEX_DOTFILES")

		fmt.Printf("Nix Config Path:  %s\n", getConfigValue(nixConfig, "not set"))
		fmt.Printf("Dotfiles Path:    %s\n", getConfigValue(dotfiles, "not set"))
		fmt.Println()

		// TODO: Also read from config file when implemented
		fmt.Println("Note: Configuration file support coming soon (~/.config/codex/config.yaml)")
	},
}

// configSetCmd sets configuration values
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Currently supported keys:
  nix-config  - Path to your Nix configuration repository
  dotfiles    - Path to your dotfiles or home-manager repository

Examples:
  codex config set nix-config /path/to/nix-config
  codex config set dotfiles /path/to/dotfiles`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if verbose {
			fmt.Printf("Setting %s = %s\n", key, value)
		}

		// TODO: Implement actual config file writing
		// For now, just provide instructions
		switch key {
		case "nix-config":
			fmt.Printf("To set nix-config path, add to your shell profile:\n")
			fmt.Printf("  export CODEX_NIX_CONFIG=\"%s\"\n", value)
			fmt.Println("\nNote: Config file persistence coming soon")
		case "dotfiles":
			fmt.Printf("To set dotfiles path, add to your shell profile:\n")
			fmt.Printf("  export CODEX_DOTFILES=\"%s\"\n", value)
			fmt.Println("\nNote: Config file persistence coming soon")
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown configuration key '%s'\n", key)
			fmt.Fprintln(os.Stderr, "Supported keys: nix-config, dotfiles")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

// Helper function to return value or default
func getConfigValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
