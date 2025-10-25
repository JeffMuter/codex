# Codex

A context-aware CLI tool that provides intelligent assistance based on your complete development environment setup.

## Overview

Codex is a CLI assistant that understands your entire development context - from your system configuration and dotfiles to your current working repository. By analyzing your Nix configuration, home-manager setup, and current workspace, codex provides personalized recommendations, keybind information, and tooling suggestions tailored to your specific environment.

## How It Works

When you run `codex ask` from anywhere in your terminal, codex automatically gathers:

1. **System Configuration**: Your Nix configuration repository
2. **Dotfiles/Home Manager**: Your dotfiles or Nix home-manager configuration
3. **Current Repository**: The git repository at or above your current directory (optional context)
4. **Current Filesystem Position**: Your current working directory and its parent hierarchy
5. **Screenshots**: Optional visual context via flag

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

Codex needs to know about your configuration repositories. Set these environment variables or configure via the CLI:

```bash
# Your Nix configuration repository path
export CODEX_NIX_CONFIG="/path/to/nix-config"

# Your dotfiles or home-manager repository path
export CODEX_DOTFILES="/path/to/dotfiles"
```

## Usage

### Basic Query

```bash
codex ask "What's my tmux prefix key?"
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

## Features

- **Context-Aware**: Understands your complete environment setup
- **Repository Detection**: Automatically finds and includes current git repository context
- **Filesystem Traversal**: Analyzes current position and parent directories
- **Visual Context**: Optional screenshot support for UI-related questions
- **Personalized Recommendations**: Suggestions based on your actual tooling
- **Keybind Discovery**: Find keybindings across all your configured tools
- **Configuration Analysis**: Deep understanding of your Nix, dotfiles, and tool configs

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

1. **Discovery**: Locate configured repositories and current workspace
2. **Analysis**: Parse relevant configuration files
3. **Context Building**: Assemble comprehensive environment snapshot
4. **Query Processing**: Send context + query to AI assistant
5. **Response**: Receive personalized, context-aware answers

## License

[Specify your license here]
