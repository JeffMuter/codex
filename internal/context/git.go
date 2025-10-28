package context

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// findGitRepository traverses up from the working directory to find a .git directory
func findGitRepository(startDir string) (string, error) {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	// Make path absolute
	absPath, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	currentDir := absPath
	for {
		gitDir := filepath.Join(currentDir, ".git")
		info, err := os.Stat(gitDir)
		if err == nil {
			// Found .git - could be directory or file (for worktrees/submodules)
			if info.IsDir() {
				return currentDir, nil
			}
			// If .git is a file, it's a worktree or submodule reference
			// Read the file to find the actual git directory
			content, err := os.ReadFile(gitDir)
			if err == nil && strings.HasPrefix(string(content), "gitdir:") {
				return currentDir, nil
			}
		}

		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached filesystem root without finding .git
			return "", fmt.Errorf("not a git repository (or any parent up to mount point)")
		}
		currentDir = parentDir
	}
}

// getGitRemote retrieves the remote URL for the repository
func getGitRemote(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		// No origin remote is not necessarily an error
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

