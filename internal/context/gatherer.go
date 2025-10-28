package context

import (
	"fmt"
	"time"

	"codex/internal/config"
)

// Gatherer collects context information from various sources
type Gatherer struct {
	nixConfigPath   string
	dotfilesPath    string
	configuredRepos []config.ConfiguredRepo
	repoFetcher     *RepoFetcher
}

// NewGatherer creates a new context gatherer
func NewGatherer(cfg *config.Config) *Gatherer {
	return &Gatherer{
		nixConfigPath:   cfg.NixConfigPath,
		dotfilesPath:    cfg.DotfilesPath,
		configuredRepos: cfg.ConfiguredRepos,
		repoFetcher:     NewRepoFetcher(),
	}
}

// Gather collects all requested context information
func (g *Gatherer) Gather(opts GatherOptions) (*Context, error) {
	ctx := &Context{
		Timestamp: time.Now(),
	}

	// Always gather configured repositories
	configuredRepos, err := g.gatherConfiguredRepos()
	if err != nil {
		// Log warning but continue
		// TODO: Add proper logging
	} else {
		ctx.ConfiguredRepos = configuredRepos
	}

	// Optionally gather current repository (opt-in via flag)
	if opts.IncludeCurrentRepo {
		currentRepo, err := g.gatherCurrentRepo(opts.WorkingDir)
		if err != nil {
			// Don't fail if we're not in a git repository
			// Just log and continue
			// TODO: Add proper logging
		} else {
			ctx.CurrentRepo = currentRepo
		}
	}

	// Gather filesystem context
	if opts.IncludeFilesystem {
		fsCtx, err := g.gatherFilesystem(opts.WorkingDir)
		if err != nil {
			return nil, fmt.Errorf("failed to gather filesystem context: %w", err)
		}
		ctx.Filesystem = fsCtx
	}

	// Gather Nix configuration
	if opts.IncludeNixConfig && g.nixConfigPath != "" {
		nixCtx, err := g.gatherNixConfig()
		if err != nil {
			// Log warning but continue
			// TODO: Add proper logging
		} else {
			ctx.NixConfig = nixCtx
		}
	}

	// Gather dotfiles configuration
	if opts.IncludeDotfiles && g.dotfilesPath != "" {
		dotfilesCtx, err := g.gatherDotfiles()
		if err != nil {
			// Log warning but continue
			// TODO: Add proper logging
		} else {
			ctx.Dotfiles = dotfilesCtx
		}
	}

	// Capture screenshot
	if opts.CaptureScreenshot {
		screenshot, err := g.captureScreenshot()
		if err != nil {
			return nil, fmt.Errorf("failed to capture screenshot: %w", err)
		}
		ctx.Screenshot = screenshot
	}

	return ctx, nil
}

// gatherConfiguredRepos fetches all configured repositories
func (g *Gatherer) gatherConfiguredRepos() ([]*RepositoryContext, error) {
	var repos []*RepositoryContext

	for _, configuredRepo := range g.configuredRepos {
		repoPath, err := g.repoFetcher.FetchRepo(configuredRepo)
		if err != nil {
			// Log error but continue with other repos
			// TODO: Add proper logging
			continue
		}

		// Get git remote if available
		remote, _ := getGitRemote(repoPath)

		repos = append(repos, &RepositoryContext{
			Path:   repoPath,
			Remote: remote,
			Source: configuredRepo.Source,
			Type:   configuredRepo.Type,
		})
	}

	return repos, nil
}

// gatherCurrentRepo collects current working directory repository information
func (g *Gatherer) gatherCurrentRepo(workingDir string) (*RepositoryContext, error) {
	repoPath, err := findGitRepository(workingDir)
	if err != nil {
		return nil, err
	}

	remote, _ := getGitRemote(repoPath)

	return &RepositoryContext{
		Path:   repoPath,
		Remote: remote,
		Type:   "current",
	}, nil
}

// gatherFilesystem collects current directory information
func (g *Gatherer) gatherFilesystem(workingDir string) (*FilesystemContext, error) {
	// TODO: Implement filesystem context gathering
	// - Get current directory
	// - Build parent directory hierarchy
	// - List files in current directory
	return nil, fmt.Errorf("not yet implemented")
}

// gatherNixConfig parses Nix configuration
func (g *Gatherer) gatherNixConfig() (*NixContext, error) {
	// TODO: Implement Nix config parsing
	// - Detect if flake or traditional config
	// - Parse .nix files
	// - Extract installed packages
	// - Extract system configuration
	return nil, fmt.Errorf("not yet implemented")
}

// gatherDotfiles parses dotfiles configuration
func (g *Gatherer) gatherDotfiles() (*DotfilesContext, error) {
	// TODO: Implement dotfiles parsing
	// - Detect if home-manager or plain dotfiles
	// - Parse common config files
	// - Extract keybindings
	return nil, fmt.Errorf("not yet implemented")
}

// captureScreenshot captures a screenshot
func (g *Gatherer) captureScreenshot() (*Screenshot, error) {
	// TODO: Implement screenshot capture
	// - Detect available screenshot tools
	// - Capture screenshot
	// - Save to temp location
	return nil, fmt.Errorf("not yet implemented")
}
