package context

import (
	"fmt"
	"time"
)

// Gatherer collects context information from various sources
type Gatherer struct {
	nixConfigPath string
	dotfilesPath  string
}

// NewGatherer creates a new context gatherer
func NewGatherer(nixConfigPath, dotfilesPath string) *Gatherer {
	return &Gatherer{
		nixConfigPath: nixConfigPath,
		dotfilesPath:  dotfilesPath,
	}
}

// Gather collects all requested context information
func (g *Gatherer) Gather(opts GatherOptions) (*Context, error) {
	ctx := &Context{
		Timestamp: time.Now(),
	}

	// Gather repository context
	if opts.IncludeRepository {
		repoCtx, err := g.gatherRepository(opts.WorkingDir)
		if err != nil {
			// Don't fail if we're not in a git repository
			// Just log and continue
		} else {
			ctx.Repository = repoCtx
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

// gatherRepository collects git repository information
func (g *Gatherer) gatherRepository(workingDir string) (*RepositoryContext, error) {
	// TODO: Implement git repository detection and parsing
	// - Find .git directory by traversing up
	// - Parse git config for remote
	// - Get current branch
	// - Get recent commits
	// - Get git status
	return nil, fmt.Errorf("not yet implemented")
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
