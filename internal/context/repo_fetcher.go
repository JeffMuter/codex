package context

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"codex/internal/config"
)

// RepoFetcher handles fetching and caching of repositories
type RepoFetcher struct {
	cacheDir string
}

// NewRepoFetcher creates a new repository fetcher
func NewRepoFetcher() *RepoFetcher {
	return &RepoFetcher{
		cacheDir: config.GetRepoCachePath(),
	}
}

// FetchRepo fetches a repository (local or remote) and returns its path
func (rf *RepoFetcher) FetchRepo(repo config.ConfiguredRepo) (string, error) {
	if repo.Type == "local" {
		// For local repos, just verify the path exists
		if _, err := os.Stat(repo.Source); err != nil {
			return "", fmt.Errorf("local repository not found: %s: %w", repo.Source, err)
		}
		return repo.Source, nil
	}

	// For remote repos, clone or update
	return rf.fetchRemoteRepo(repo)
}

// fetchRemoteRepo clones or updates a remote repository
func (rf *RepoFetcher) fetchRemoteRepo(repo config.ConfiguredRepo) (string, error) {
	cachePath := repo.CachePath

	// Ensure cache directory exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Check if already cloned
	gitDir := filepath.Join(cachePath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		// Repository already exists, try to update it
		if err := rf.updateRepo(cachePath); err != nil {
			// If update fails, log but continue with existing cached version
			// TODO: Add proper logging
			return cachePath, nil
		}
		return cachePath, nil
	}

	// Repository not cached, clone it
	return cachePath, rf.cloneRepo(repo.Source, cachePath)
}

// cloneRepo clones a remote repository
func (rf *RepoFetcher) cloneRepo(url, destPath string) error {
	cmd := exec.Command("git", "clone", url, destPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository %s: %w\nOutput: %s", url, err, output)
	}
	return nil
}

// updateRepo updates an existing cloned repository
func (rf *RepoFetcher) updateRepo(repoPath string) error {
	// Check if we should update (don't update too frequently)
	if !rf.shouldUpdate(repoPath) {
		return nil
	}

	// Fetch latest changes
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch updates: %w", err)
	}

	// Reset to origin/main or origin/master
	// First, try to determine the default branch
	branch, err := rf.getDefaultBranch(repoPath)
	if err != nil {
		// If we can't determine, try common defaults
		branch = "main"
	}

	cmd = exec.Command("git", "reset", "--hard", fmt.Sprintf("origin/%s", branch))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		// Try master if main fails
		cmd = exec.Command("git", "reset", "--hard", "origin/master")
		cmd.Dir = repoPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to reset repository: %w", err)
		}
	}

	// Update the timestamp file
	rf.updateTimestamp(repoPath)

	return nil
}

// shouldUpdate checks if repository should be updated based on last update time
func (rf *RepoFetcher) shouldUpdate(repoPath string) bool {
	timestampFile := filepath.Join(repoPath, ".codex_last_update")
	info, err := os.Stat(timestampFile)
	if err != nil {
		// No timestamp file, should update
		return true
	}

	// Update if last update was more than 1 hour ago
	return time.Since(info.ModTime()) > time.Hour
}

// updateTimestamp updates the timestamp file to track last update
func (rf *RepoFetcher) updateTimestamp(repoPath string) {
	timestampFile := filepath.Join(repoPath, ".codex_last_update")
	os.WriteFile(timestampFile, []byte(time.Now().Format(time.RFC3339)), 0644)
}

// getDefaultBranch attempts to determine the default branch of a repository
func (rf *RepoFetcher) getDefaultBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Output is like "refs/remotes/origin/main"
	branch := filepath.Base(string(output))
	return branch, nil
}
