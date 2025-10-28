package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindGitRepository(t *testing.T) {
	// Test finding git repository from current directory
	// This test assumes we're running in a git repository
	repoPath, err := findGitRepository("")
	if err != nil {
		t.Skipf("Skipping test: not in a git repository: %v", err)
	}

	// Verify .git exists in the found path
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("Expected .git directory at %s, but got error: %v", gitDir, err)
	}

	t.Logf("Found git repository at: %s", repoPath)
}

func TestFindGitRepositoryFromSubdirectory(t *testing.T) {
	// First find the repository root
	repoPath, err := findGitRepository("")
	if err != nil {
		t.Skipf("Skipping test: not in a git repository: %v", err)
	}

	// Try finding from a subdirectory
	subDir := filepath.Join(repoPath, "internal", "context")
	foundPath, err := findGitRepository(subDir)
	if err != nil {
		t.Errorf("Failed to find git repository from subdirectory: %v", err)
	}

	if foundPath != repoPath {
		t.Errorf("Expected repository path %s, but got %s", repoPath, foundPath)
	}
}

func TestFindGitRepositoryNotFound(t *testing.T) {
	// Test with a directory that's definitely not in a git repository
	_, err := findGitRepository("/")
	if err == nil {
		t.Error("Expected error when searching for git repository from root, but got nil")
	}
}

func TestGetGitRemote(t *testing.T) {
	// Find repository first
	repoPath, err := findGitRepository("")
	if err != nil {
		t.Skipf("Skipping test: not in a git repository: %v", err)
	}

	// Get remote - may or may not exist, both are valid
	remote, err := getGitRemote(repoPath)
	if err != nil {
		t.Errorf("Failed to get git remote: %v", err)
	}

	// If remote exists, log it
	if remote != "" {
		t.Logf("Found git remote: %s", remote)
	} else {
		t.Log("No git remote configured (this is OK)")
	}
}

func TestGatherCurrentRepo(t *testing.T) {
	// Test gathering current working directory repository information
	repoPath, err := findGitRepository("")
	if err != nil {
		t.Skipf("Skipping test: not in a git repository: %v", err)
	}

	// Get remote (may or may not exist)
	remote, _ := getGitRemote(repoPath)

	ctx := &RepositoryContext{
		Path:   repoPath,
		Remote: remote,
		Type:   "current",
	}

	// Verify we have a path
	if ctx.Path == "" {
		t.Error("Expected repository path, but got empty string")
	}

	// Verify .git exists at the path
	gitDir := filepath.Join(ctx.Path, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("Expected .git directory at %s, but got error: %v", gitDir, err)
	}

	t.Logf("Repository context: Path=%s, Remote=%s, Type=%s", ctx.Path, ctx.Remote, ctx.Type)
}
