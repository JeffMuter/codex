package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfiguredRepo represents a repository that should always be included as context
type ConfiguredRepo struct {
	Source    string `yaml:"source"`     // Local path or remote URL
	Type      string `yaml:"type"`       // "local" or "remote"
	CachePath string `yaml:"cache_path"` // For remote repos, where they're cached locally
}

// Config represents the application configuration
type Config struct {
	// Paths to configuration repositories for parsing (Nix/dotfiles)
	NixConfigPath string `yaml:"nix_config_path"`
	DotfilesPath  string `yaml:"dotfiles_path"`

	// Repositories to always include as context
	ConfiguredRepos []ConfiguredRepo `yaml:"configured_repos"`

	// AI Provider settings
	Provider     string `yaml:"provider"` // "anthropic", "openai", "ollama"
	AnthropicKey string `yaml:"anthropic_key,omitempty"`
	OpenAIKey    string `yaml:"openai_key,omitempty"`
	OllamaURL    string `yaml:"ollama_url,omitempty"`

	// Database settings
	DatabasePath string `yaml:"database_path"`

	// Cache settings
	CacheTTL int `yaml:"cache_ttl"` // in hours
}

// Default configuration values
const (
	DefaultProvider    = "anthropic"
	DefaultCacheTTL    = 24 // 24 hours
	DefaultOllamaURL   = "http://localhost:11434"
	DefaultConfigDir   = ".config/codex"
	DefaultDataDir     = ".local/share/codex"
	DefaultConfigFile  = "config.yaml"
	DefaultDatabaseFile = "codex.db"
)

// Load reads configuration from file and environment variables
// Environment variables take precedence over config file values
func Load() (*Config, error) {
	cfg := &Config{
		Provider:  DefaultProvider,
		CacheTTL:  DefaultCacheTTL,
		OllamaURL: DefaultOllamaURL,
	}

	// Set default paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cfg.DatabasePath = filepath.Join(homeDir, DefaultDataDir, DefaultDatabaseFile)

	// Try to load from config file
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		if err := loadFromFile(cfg, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromFile reads configuration from YAML file
func loadFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// loadFromEnv overrides configuration with environment variables
func loadFromEnv(cfg *Config) {
	if val := os.Getenv("CODEX_NIX_CONFIG"); val != "" {
		cfg.NixConfigPath = val
	}
	if val := os.Getenv("CODEX_DOTFILES"); val != "" {
		cfg.DotfilesPath = val
	}
	if val := os.Getenv("CODEX_PROVIDER"); val != "" {
		cfg.Provider = val
	}
	if val := os.Getenv("ANTHROPIC_API_KEY"); val != "" {
		cfg.AnthropicKey = val
	}
	if val := os.Getenv("OPENAI_API_KEY"); val != "" {
		cfg.OpenAIKey = val
	}
	if val := os.Getenv("CODEX_OLLAMA_URL"); val != "" {
		cfg.OllamaURL = val
	}
	if val := os.Getenv("CODEX_DATABASE_PATH"); val != "" {
		cfg.DatabasePath = val
	}
}

// Save writes configuration to file
func (cfg *Config) Save() error {
	configPath := GetConfigPath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (cfg *Config) Validate() error {
	// Check if paths exist (if set)
	if cfg.NixConfigPath != "" {
		if _, err := os.Stat(cfg.NixConfigPath); os.IsNotExist(err) {
			return fmt.Errorf("nix config path does not exist: %s", cfg.NixConfigPath)
		}
	}

	if cfg.DotfilesPath != "" {
		if _, err := os.Stat(cfg.DotfilesPath); os.IsNotExist(err) {
			return fmt.Errorf("dotfiles path does not exist: %s", cfg.DotfilesPath)
		}
	}

	// Validate provider
	switch cfg.Provider {
	case "anthropic":
		if cfg.AnthropicKey == "" {
			return fmt.Errorf("anthropic provider requires ANTHROPIC_API_KEY")
		}
	case "openai":
		if cfg.OpenAIKey == "" {
			return fmt.Errorf("openai provider requires OPENAI_API_KEY")
		}
	case "ollama":
		// Ollama doesn't require API key, but URL should be valid
		if cfg.OllamaURL == "" {
			cfg.OllamaURL = DefaultOllamaURL
		}
	default:
		return fmt.Errorf("unknown provider: %s (must be anthropic, openai, or ollama)", cfg.Provider)
	}

	return nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return DefaultConfigFile
	}
	return filepath.Join(homeDir, DefaultConfigDir, DefaultConfigFile)
}

// GetDatabasePath returns the full path to the database file
func GetDatabasePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return DefaultDatabaseFile
	}
	return filepath.Join(homeDir, DefaultDataDir, DefaultDatabaseFile)
}

// EnsureDirectories creates necessary directories for config and data
func EnsureDirectories() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, DefaultConfigDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, DefaultDataDir)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create cache directory for remote repos
	cacheDir := filepath.Join(homeDir, DefaultDataDir, "repos")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	return nil
}

// GetRepoCachePath returns the directory for caching remote repositories
func GetRepoCachePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "repos"
	}
	return filepath.Join(homeDir, DefaultDataDir, "repos")
}

// AddRepo adds a repository to the configured repos list
func (cfg *Config) AddRepo(source string) error {
	// Check if already exists
	for _, repo := range cfg.ConfiguredRepos {
		if repo.Source == source {
			return fmt.Errorf("repository already configured: %s", source)
		}
	}

	// Determine type (local or remote)
	repoType := "local"
	cachePath := ""

	if isRemoteURL(source) {
		repoType = "remote"
		// Generate cache path from URL
		cachePath = generateCachePath(source)
	} else {
		// Make path absolute
		absPath, err := filepath.Abs(source)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}
		source = absPath

		// Verify path exists
		if _, err := os.Stat(source); os.IsNotExist(err) {
			return fmt.Errorf("local path does not exist: %s", source)
		}
	}

	cfg.ConfiguredRepos = append(cfg.ConfiguredRepos, ConfiguredRepo{
		Source:    source,
		Type:      repoType,
		CachePath: cachePath,
	})

	return nil
}

// RemoveRepo removes a repository from the configured repos list
func (cfg *Config) RemoveRepo(source string) error {
	// Make path absolute if it's a local path
	if !isRemoteURL(source) {
		absPath, err := filepath.Abs(source)
		if err == nil {
			source = absPath
		}
	}

	found := false
	newRepos := []ConfiguredRepo{}
	for _, repo := range cfg.ConfiguredRepos {
		if repo.Source != source {
			newRepos = append(newRepos, repo)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("repository not found in configuration: %s", source)
	}

	cfg.ConfiguredRepos = newRepos
	return nil
}

// isRemoteURL checks if a source string is a remote URL
func isRemoteURL(source string) bool {
	return len(source) > 4 && (source[:4] == "http" || source[:4] == "git@" || source[:3] == "ssh")
}

// generateCachePath creates a cache directory name from a URL
func generateCachePath(url string) string {
	// Simple implementation: hash or sanitize the URL
	// For now, just extract the repo name
	parts := filepath.Base(url)
	if len(parts) > 4 && parts[len(parts)-4:] == ".git" {
		parts = parts[:len(parts)-4]
	}
	return filepath.Join(GetRepoCachePath(), parts)
}
