# Codex

A context-aware CLI tool that provides intelligent assistance based on your complete development environment setup.

## Overview

Codex is a CLI assistant that understands your entire development context - from your system configuration and dotfiles to your current working repository. By analyzing your Nix configuration, home-manager setup, and current workspace, codex provides personalized recommendations, keybind information, and tooling suggestions tailored to your specific environment.

## How It Works

When you run `codex ask` from anywhere in your terminal, codex automatically gathers:

1. **Configured Repositories**: Repositories you've explicitly configured (always included as context)
2. **System Configuration**: Your Nix configuration repository (parsed for packages and settings)
3. **Dotfiles/Home Manager**: Your dotfiles or Nix home-manager configuration (parsed for keybinds and configs)
4. **Current Repository**: The git repository at or above your current directory (opt-in via `--current-repo` flag)
5. **Current Filesystem Position**: Your current working directory and its parent hierarchy
6. **Screenshots**: Optional visual context via flag

This rich context enables codex to provide answers about:
- Your configured keybinds across all tools
- Available CLI tools and their configurations
- Better alternatives based on your existing setup
- Environment-specific recommendations
- Tool integration possibilities

## Installation

```bash
# Build from source
go build -o codex main.go

# Install to your PATH
go install
```

## Configuration

### Configured Repositories

Configure repositories that should always be included as context:

```bash
# Add a local repository
codex config add-repo /path/to/my-templates

# Add a remote repository (will be cloned and cached)
codex config add-repo https://github.com/user/dotfiles

# List all configured repositories
codex config list-repos

# Remove a repository
codex config remove-repo /path/to/my-templates
```

### Nix and Dotfiles Configuration

Set paths to your Nix and dotfiles for configuration parsing:

```bash
# Your Nix configuration repository path
export CODEX_NIX_CONFIG="/path/to/nix-config"

# Your dotfiles or home-manager repository path
export CODEX_DOTFILES="/path/to/dotfiles"
```

These are parsed for packages, keybindings, and tool configurations (different from repository context).

## Usage

### Basic Query

```bash
codex ask "What's my tmux prefix key?"
```

### Include Current Repository

```bash
# Include the current working directory's git repository as context
codex ask --current-repo "Explain the architecture of this codebase"
codex ask -r "What patterns are used in this project?"
```

### With Screenshot Context

```bash
codex ask --screenshot "How do I achieve this layout in my window manager?"
```

### Example Queries

- "What keybind do I use for fuzzy file search in neovim?"
- "What CLI tools do I have for working with JSON?"
- "Is there a better way to do X given my current setup?"
- "How do I configure Y to work with my existing Z?"
- "What's my current shell configuration?"
- `codex ask -r "What testing framework is this project using?"`

## Features

- **Context-Aware**: Understands your complete environment setup
- **Configured Repository Context**: Always includes repositories you've explicitly configured
- **Current Repository (Opt-in)**: Include current working directory repo with `--current-repo` flag
- **Remote Repository Support**: Clone and cache remote repositories (GitHub, GitLab, etc.)
- **Filesystem Traversal**: Analyzes current position and parent directories
- **Visual Context**: Optional screenshot support for UI-related questions
- **Personalized Recommendations**: Suggestions based on your actual tooling
- **Keybind Discovery**: Find keybindings across all your configured tools
- **Configuration Analysis**: Deep understanding of your Nix, dotfiles, and tool configs
- **Privacy-Conscious**: No automatic inclusion of current repository without explicit flag

## Development

This project uses Nix for development environment management.

```bash
# Enter development shell
nix-shell shell.nix

# Run tests
go test ./...

# Build
go build -o codex main.go
```

## Architecture

Codex operates in phases:

1. **Repository Fetching**: Fetch all configured repositories (local or remote)
2. **Optional Current Repo**: If `--current-repo` flag is used, include working directory repo
3. **Configuration Parsing**: Parse Nix configs and dotfiles for settings and keybindings
4. **Context Building**: Assemble comprehensive environment snapshot
5. **Query Processing**: Send context + query to AI assistant
6. **Response**: Receive personalized, context-aware answers

### Repository Types

- **Configured Repos**: Explicitly added via `codex config add-repo`, always included
  - Local paths: Direct access to repository contents
  - Remote URLs: Cloned to `~/.local/share/codex/repos/` and kept updated
- **Current Repo**: Working directory repository, only included with `--current-repo` flag
- **Config Repos**: Nix/dotfiles paths parsed for configuration (not included as raw context)

## License

[Specify your license here]
