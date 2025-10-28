#! /usr/bin/env nix-shell
{ pkgs ? import <nixpkgs> { 
    config = { 
      allowUnfree = true;
      allowUnsupportedSystem = true;
    };
  }
}:
let
  # Import unstable channel for claude-code
  unstable = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {
    config = {
      allowUnfree = true;
      allowUnsupportedSystem = true;
    };
  };
in
pkgs.mkShell {
  buildInputs = with pkgs; [
    # Go toolchain
    go_1_23

    # Go development tools
    gopls           # Go language server
    gotools         # goimports, godoc, etc.
    go-tools        # staticcheck, etc.
    delve           # Go debugger

    # Database tools
    goose           # Database migration tool
    sqlite          # SQLite database

    # Development utilities
    unstable.claude-code
  ];

  # Allow unsafe/experimental features
  NIX_CONFIG = ''
    experimental-features = nix-command flakes
    allow-import-from-derivation = true
  '';

  # Set up Go environment
  shellHook = ''
    echo "Codex development environment loaded"
    echo "Go version: $(go version)"
    echo "Database: SQLite $(sqlite3 --version | cut -d' ' -f1)"
    echo ""
    echo "Available commands:"
    echo "  go build    - Build the project"
    echo "  go test     - Run tests"
    echo "  goose       - Database migrations"
    echo "  dlv         - Go debugger"
    echo ""
  '';
}
