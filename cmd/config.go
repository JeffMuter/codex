package cmd

import (
	"fmt"
	"os"

	"codex/internal/config"
	"codex/internal/logging"

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
		logging.Logger.Debug().Msg("Displaying configuration")

		fmt.Println("Codex Configuration")
		fmt.Println("===================")
		fmt.Println()

		// Check environment variables
		nixConfig := os.Getenv("CODEX_NIX_CONFIG")
		dotfiles := os.Getenv("CODEX_DOTFILES")

		logging.Logger.Debug().
			Str("nix_config", nixConfig).
			Str("dotfiles", dotfiles).
			Msg("Current configuration values")

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

		logging.Logger.Debug().
			Str("key", key).
			Str("value", value).
			Msg("Setting configuration value")

		// TODO: Implement actual config file writing
		// For now, just provide instructions
		switch key {
		case "nix-config":
			logging.Logger.Info().
				Str("key", "CODEX_NIX_CONFIG").
				Str("value", value).
				Msg("Configuration set instruction provided")
			fmt.Printf("To set nix-config path, add to your shell profile:\n")
			fmt.Printf("  export CODEX_NIX_CONFIG=\"%s\"\n", value)
			fmt.Println("\nNote: Config file persistence coming soon")
		case "dotfiles":
			logging.Logger.Info().
				Str("key", "CODEX_DOTFILES").
				Str("value", value).
				Msg("Configuration set instruction provided")
			fmt.Printf("To set dotfiles path, add to your shell profile:\n")
			fmt.Printf("  export CODEX_DOTFILES=\"%s\"\n", value)
			fmt.Println("\nNote: Config file persistence coming soon")
		default:
			logging.Logger.Error().
				Str("key", key).
				Msg("Unknown configuration key")
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
	configCmd.AddCommand(configAddRepoCmd)
	configCmd.AddCommand(configListReposCmd)
	configCmd.AddCommand(configRemoveRepoCmd)
}

// configAddRepoCmd adds a repository to the configured repos list
var configAddRepoCmd = &cobra.Command{
	Use:   "add-repo [path-or-url]",
	Short: "Add a repository to always include as context",
	Long: `Add a repository to be included as context in all queries.
Can be a local path or a remote git URL.

Examples:
  codex config add-repo /path/to/local/repo
  codex config add-repo https://github.com/user/repo
  codex config add-repo git@github.com:user/repo.git`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]

		logging.Logger.Debug().
			Str("source", source).
			Msg("Adding repository to configuration")

		// Load config
		cfg, err := config.Load()
		if err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to load configuration")
			fmt.Fprintf(os.Stderr, "Error: failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		// Add repo
		if err := cfg.AddRepo(source); err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to add repository")
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Ensure directories exist
		if err := config.EnsureDirectories(); err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to create directories")
			fmt.Fprintf(os.Stderr, "Error: failed to create directories: %v\n", err)
			os.Exit(1)
		}

		// Save config
		if err := cfg.Save(); err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to save configuration")
			fmt.Fprintf(os.Stderr, "Error: failed to save configuration: %v\n", err)
			os.Exit(1)
		}

		logging.Logger.Info().
			Str("source", source).
			Msg("Repository added successfully")
		fmt.Printf("✓ Repository added: %s\n", source)
	},
}

// configListReposCmd lists all configured repositories
var configListReposCmd = &cobra.Command{
	Use:   "list-repos",
	Short: "List all configured repositories",
	Long:  `Display all repositories that are configured to be included as context.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.Logger.Debug().Msg("Listing configured repositories")

		// Load config
		cfg, err := config.Load()
		if err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to load configuration")
			fmt.Fprintf(os.Stderr, "Error: failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configured Repositories")
		fmt.Println("=======================")
		fmt.Println()

		if len(cfg.ConfiguredRepos) == 0 {
			fmt.Println("No repositories configured.")
			fmt.Println()
			fmt.Println("Add a repository with: codex config add-repo [path-or-url]")
			return
		}

		for i, repo := range cfg.ConfiguredRepos {
			fmt.Printf("%d. %s (%s)\n", i+1, repo.Source, repo.Type)
			if repo.Type == "remote" && repo.CachePath != "" {
				fmt.Printf("   Cached at: %s\n", repo.CachePath)
			}
		}
	},
}

// configRemoveRepoCmd removes a repository from the configured repos list
var configRemoveRepoCmd = &cobra.Command{
	Use:   "remove-repo [path-or-url]",
	Short: "Remove a repository from configured repos",
	Long: `Remove a repository from being included as context.

Examples:
  codex config remove-repo /path/to/local/repo
  codex config remove-repo https://github.com/user/repo`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]

		logging.Logger.Debug().
			Str("source", source).
			Msg("Removing repository from configuration")

		// Load config
		cfg, err := config.Load()
		if err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to load configuration")
			fmt.Fprintf(os.Stderr, "Error: failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		// Remove repo
		if err := cfg.RemoveRepo(source); err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to remove repository")
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Save config
		if err := cfg.Save(); err != nil {
			logging.Logger.Error().Err(err).Msg("Failed to save configuration")
			fmt.Fprintf(os.Stderr, "Error: failed to save configuration: %v\n", err)
			os.Exit(1)
		}

		logging.Logger.Info().
			Str("source", source).
			Msg("Repository removed successfully")
		fmt.Printf("✓ Repository removed: %s\n", source)
	},
}

// Helper function to return value or default
func getConfigValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
