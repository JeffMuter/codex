package context

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ContextSummarizer provides intelligent context summarization
type ContextSummarizer struct {
	maxContextSize int // Maximum context size in bytes
}

// NewContextSummarizer creates a new context summarizer
func NewContextSummarizer(maxSize int) *ContextSummarizer {
	return &ContextSummarizer{
		maxContextSize: maxSize,
	}
}

// SummarizeRepoContents creates a summarized version of repository contents
func (cs *ContextSummarizer) SummarizeRepoContents(contents *RepoContents) *RepoContents {
	if contents == nil || contents.TotalSize <= cs.maxContextSize {
		return contents // No need to summarize
	}

	summarized := &RepoContents{
		Files:      make([]FileContent, 0),
		TotalFiles: contents.TotalFiles,
		TotalSize:  0,
	}

	// Prioritize files by importance
	prioritized := cs.prioritizeFiles(contents.Files)

	// Add files until we hit the size limit
	for _, file := range prioritized {
		if summarized.TotalSize+file.Size > cs.maxContextSize {
			// Try to add a truncated version
			remaining := cs.maxContextSize - summarized.TotalSize
			if remaining > 500 { // Only add if we have at least 500 bytes left
				truncated := cs.truncateFile(file, remaining)
				summarized.Files = append(summarized.Files, truncated)
				summarized.TotalSize += truncated.Size
			}
			break
		}

		summarized.Files = append(summarized.Files, file)
		summarized.TotalSize += file.Size
	}

	return summarized
}

// prioritizeFiles orders files by importance for context
func (cs *ContextSummarizer) prioritizeFiles(files []FileContent) []FileContent {
	// Create a copy to avoid modifying original
	prioritized := make([]FileContent, len(files))
	copy(prioritized, files)

	// Sort by priority (higher priority first)
	// Priority order:
	// 1. Configuration files (highest priority)
	// 2. Source code
	// 3. Documentation
	// 4. Other files

	type scoredFile struct {
		file  FileContent
		score int
	}

	scored := make([]scoredFile, len(prioritized))
	for i, file := range prioritized {
		scored[i] = scoredFile{
			file:  file,
			score: cs.calculatePriority(file),
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	result := make([]FileContent, len(scored))
	for i, sf := range scored {
		result[i] = sf.file
	}

	return result
}

// calculatePriority assigns a priority score to a file
func (cs *ContextSummarizer) calculatePriority(file FileContent) int {
	filename := filepath.Base(file.RelativePath)
	ext := strings.ToLower(filepath.Ext(filename))
	dir := filepath.Dir(file.RelativePath)

	score := 0

	// High priority: Configuration files
	configFiles := []string{
		"flake.nix", "configuration.nix", "home.nix", "hardware-configuration.nix",
		".bashrc", ".zshrc", ".vimrc", ".nvimrc", "init.vim", "init.lua",
		"config.toml", "config.yaml", "config.yml", "config.json",
		".gitconfig", ".tmux.conf", "alacritty.yml", "kitty.conf",
	}
	for _, cf := range configFiles {
		if filename == cf {
			score += 1000
			break
		}
	}

	// High priority: Nix files
	if ext == ".nix" {
		score += 800
	}

	// Medium-high priority: Shell scripts and configs
	shellExts := []string{".sh", ".bash", ".zsh", ".fish"}
	for _, se := range shellExts {
		if ext == se {
			score += 600
			break
		}
	}

	// Medium priority: Source code
	sourceExts := []string{".go", ".rs", ".py", ".js", ".ts", ".c", ".cpp", ".h", ".hpp"}
	for _, se := range sourceExts {
		if ext == se {
			score += 400
			break
		}
	}

	// Lower priority: Documentation
	docExts := []string{".md", ".txt", ".rst"}
	for _, de := range docExts {
		if ext == de {
			score += 200
			break
		}
	}

	// Boost files in root directory
	if dir == "." {
		score += 100
	}

	// Boost README files
	if strings.HasPrefix(strings.ToLower(filename), "readme") {
		score += 150
	}

	// Reduce priority for test files (but don't exclude them)
	if strings.Contains(strings.ToLower(filename), "test") {
		score -= 50
	}

	// Smaller files are slightly preferred (easier to include more files)
	if file.Size < 5000 {
		score += 50
	} else if file.Size > 50000 {
		score -= 50
	}

	return score
}

// truncateFile creates a truncated version of a file
func (cs *ContextSummarizer) truncateFile(file FileContent, maxSize int) FileContent {
	if file.Size <= maxSize {
		return file
	}

	// Reserve space for truncation message
	truncMsg := "\n\n... [truncated] ...\n"
	availableSize := maxSize - len(truncMsg)
	if availableSize <= 0 {
		return FileContent{
			Path:         file.Path,
			RelativePath: file.RelativePath,
			Content:      "[file too large to include]",
			Size:         len("[file too large to include]"),
		}
	}

	// Take first part of file
	truncated := file.Content[:availableSize] + truncMsg

	return FileContent{
		Path:         file.Path,
		RelativePath: file.RelativePath,
		Content:      truncated,
		Size:         len(truncated),
	}
}

// SummarizeContext applies summarization to an entire context
func (cs *ContextSummarizer) SummarizeContext(ctx *Context) *Context {
	if ctx == nil {
		return ctx
	}

	summarized := &Context{
		Timestamp:   ctx.Timestamp,
		Filesystem:  ctx.Filesystem,
		NixConfig:   ctx.NixConfig,
		Dotfiles:    ctx.Dotfiles,
		Screenshot:  ctx.Screenshot,
	}

	// Calculate how much space to allocate per repo
	totalRepos := len(ctx.ConfiguredRepos)
	if ctx.CurrentRepo != nil {
		totalRepos++
	}

	if totalRepos == 0 {
		return summarized
	}

	sizePerRepo := cs.maxContextSize / totalRepos

	// Summarize configured repos
	if len(ctx.ConfiguredRepos) > 0 {
		repoSummarizer := NewContextSummarizer(sizePerRepo)
		summarized.ConfiguredRepos = make([]*RepositoryContext, len(ctx.ConfiguredRepos))
		for i, repo := range ctx.ConfiguredRepos {
			summarizedRepo := &RepositoryContext{
				Path:     repo.Path,
				Remote:   repo.Remote,
				Source:   repo.Source,
				Type:     repo.Type,
				Contents: repoSummarizer.SummarizeRepoContents(repo.Contents),
			}
			summarized.ConfiguredRepos[i] = summarizedRepo
		}
	}

	// Summarize current repo
	if ctx.CurrentRepo != nil {
		repoSummarizer := NewContextSummarizer(sizePerRepo)
		summarized.CurrentRepo = &RepositoryContext{
			Path:     ctx.CurrentRepo.Path,
			Remote:   ctx.CurrentRepo.Remote,
			Source:   ctx.CurrentRepo.Source,
			Type:     ctx.CurrentRepo.Type,
			Contents: repoSummarizer.SummarizeRepoContents(ctx.CurrentRepo.Contents),
		}
	}

	return summarized
}

// EstimateTokens provides a rough estimate of token count
func (cs *ContextSummarizer) EstimateTokens(contents *RepoContents) int {
	if contents == nil {
		return 0
	}
	// Rough estimate: ~4 characters per token
	return contents.TotalSize / 4
}

// FormatSummaryStats returns a human-readable summary of what was included
func (cs *ContextSummarizer) FormatSummaryStats(original, summarized *RepoContents) string {
	if original == nil {
		return "No content"
	}
	if summarized == nil || original.TotalSize == summarized.TotalSize {
		return fmt.Sprintf("%d files (%d bytes, ~%d tokens)",
			original.TotalFiles, original.TotalSize, cs.EstimateTokens(original))
	}

	return fmt.Sprintf("%d/%d files (%d/%d bytes, ~%d tokens) - summarized",
		len(summarized.Files), original.TotalFiles,
		summarized.TotalSize, original.TotalSize,
		cs.EstimateTokens(summarized))
}
