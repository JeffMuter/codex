package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileContent represents a single file's content
type FileContent struct {
	Path         string `json:"path"`
	RelativePath string `json:"relative_path"`
	Content      string `json:"content"`
	Size         int    `json:"size"`
}

// RepoContents represents all files from a repository
type RepoContents struct {
	Files      []FileContent `json:"files"`
	TotalSize  int           `json:"total_size"`
	TotalFiles int           `json:"total_files"`
}

// ContentReader reads repository contents with intelligent filtering
type ContentReader struct {
	maxFileSize   int64 // Maximum file size to read (in bytes)
	maxTotalSize  int64 // Maximum total content size
	includeHidden bool  // Include hidden files
}

// NewContentReader creates a new content reader with sensible defaults
func NewContentReader() *ContentReader {
	return &ContentReader{
		maxFileSize:   100 * 1024,      // 100KB per file
		maxTotalSize:  2 * 1024 * 1024, // 2MB total (~500K tokens)
		includeHidden: true,
	}
}

// ReadRepoContents reads all relevant files from a repository
func (cr *ContentReader) ReadRepoContents(repoPath string) (*RepoContents, error) {
	contents := &RepoContents{
		Files: make([]FileContent, 0),
	}

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't read
		}

		// Skip if we've exceeded total size limit
		if contents.TotalSize > int(cr.maxTotalSize) {
			return filepath.SkipAll
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories unless includeHidden is true
			if !cr.includeHidden && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			// Always skip common build/cache directories
			if cr.shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}

			return nil
		}

		// Skip files that should be ignored
		if cr.shouldSkipFile(path, info) {
			return nil
		}

		// Read the file
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			relPath = path
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		// Skip binary files
		if isBinary(content) {
			return nil
		}

		fileContent := FileContent{
			Path:         path,
			RelativePath: relPath,
			Content:      string(content),
			Size:         len(content),
		}

		contents.Files = append(contents.Files, fileContent)
		contents.TotalSize += fileContent.Size
		contents.TotalFiles++

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read repository contents: %w", err)
	}

	return contents, nil
}

// shouldSkipDirectory returns true if the directory should be skipped
func (cr *ContentReader) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		// Version control
		".git",

		// Dependencies
		"node_modules",
		"vendor",

		// Build/Dist
		"dist",
		"build",
		"target",
		".next",
		".nuxt",
		"out",

		// Caches
		".cache",
		".npm",
		".yarn",
		".gradle",
		".m2",
		".turbo",
		"tmp",
		"temp",
		".tmp",
		"Cache",           // Electron/browser cache
		"GPUCache",        // GPU cache
		"Code Cache",      // Code cache
		"DawnCache",       // WebGPU cache
		"IndexedDB",       // Browser IndexedDB
		"LocalStorage",    // Browser local storage
		"SessionStorage",  // Browser session storage
		"Service Worker",  // Service worker cache
		"adblock",         // Adblock lists

		// Python
		"__pycache__",
		".pytest_cache",
		".mypy_cache",
		".venv",
		"venv",
		".tox",

		// Coverage/Testing
		"coverage",
		".coverage",
		".nyc_output",

		// IDE/Editor
		".idea",
		".vscode",
		".DS_Store",

		// Browser/Electron specific
		"Dictionaries",    // Spell check dictionaries
		"WebStorage",      // Web storage
		"Partitions",      // Browser partitions
		"DawnGraphiteCache",
		"DawnWebGPUCache",

		// Nix
		"result", // Nix result symlinks
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// shouldSkipFile returns true if the file should be skipped
func (cr *ContentReader) shouldSkipFile(path string, info os.FileInfo) bool {
	// Skip if file is too large
	if info.Size() > cr.maxFileSize {
		return true
	}

	// Skip hidden files unless includeHidden is true
	if !cr.includeHidden && strings.HasPrefix(info.Name(), ".") {
		return true
	}

	// Skip by extension
	skipExtensions := []string{
		// Images
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
		".webp", ".bmp", ".tiff", ".tif",

		// Documents
		".pdf", ".doc", ".docx",
		".md",  // Markdown files
		".txt", // Text files (often logs, dumps, or data files)

		// Archives
		".zip", ".tar", ".gz", ".bz2", ".xz", ".rar", ".7z",

		// Executables/Libraries
		".exe", ".dll", ".so", ".dylib",
		".o", ".a", ".pyc", ".class",

		// Lock files (usually too verbose)
		".lock",
		".sum", // checksums

		// Compiled/Generated
		".wasm",      // WebAssembly binaries
		".min.js",    // Minified JavaScript
		".min.css",   // Minified CSS
		".map",       // Source maps
		".bundle.js", // Bundled JS
		".chunk.js",  // Code-split chunks
		".asar",      // Electron app archives
		".pak",       // Chrome/Electron package files

		// Databases
		".db", ".sqlite", ".sqlite3",
		".mdb", ".accdb",
		".ldb", ".leveldb", // LevelDB files (used by browsers/Electron)
		".indexeddb",       // IndexedDB files

		// Media
		".mp4", ".mov", ".avi", ".mkv", ".webm", // Video
		".mp3", ".wav", ".flac", ".ogg", ".m4a", // Audio

		// Fonts
		".woff", ".woff2", ".ttf", ".eot", ".otf",

		// Design files
		".psd", ".ai", ".sketch", ".fig", ".xd",

		// IDE files
		".iml", // IntelliJ
		".swp", ".swo", // Vim swap

		// Logs and data
		".log",
		".bdic",  // Binary dictionary files
		".dat",   // Generic data files
		".bin",   // Binary files
		".data",  // Data files
		".cache", // Cache files
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}

	// Skip specific filenames
	skipFiles := []string{
		// Lock files
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"Cargo.lock",
		"Gemfile.lock",
		"composer.lock",
		"poetry.lock",
		"go.sum",

		// OS files
		".DS_Store",
		"Thumbs.db",
		"desktop.ini",

		// Large config/generated files
		"*.min.js",
		"*.min.css",

		// Browser/Electron cache files
		"Cookies",
		"Cookies-journal",
		"History",
		"History-journal",
		"Preferences",
		"Current Session",
		"Current Tabs",
		"Last Session",
		"Last Tabs",
		"Network Persistent State",
		"TransportSecurity",
		"Web Data",
		"Web Data-journal",
		"DIPS",
		"DIPS-wal",
		"Trust Tokens",
		"Shared Dictionary",
		"QuotaManager",
		"QuotaManager-journal",
	}

	base := filepath.Base(path)
	for _, skipFile := range skipFiles {
		if base == skipFile {
			return true
		}
	}

	return false
}

// isBinary checks if content appears to be binary
func isBinary(content []byte) bool {
	// Check first 512 bytes for null bytes
	checkLen := 512
	if len(content) < checkLen {
		checkLen = len(content)
	}

	for i := 0; i < checkLen; i++ {
		if content[i] == 0 {
			return true
		}
	}

	return false
}

// SetMaxFileSize sets the maximum file size to read
func (cr *ContentReader) SetMaxFileSize(size int64) {
	cr.maxFileSize = size
}

// SetMaxTotalSize sets the maximum total content size
func (cr *ContentReader) SetMaxTotalSize(size int64) {
	cr.maxTotalSize = size
}

// SetIncludeHidden sets whether to include hidden files
func (cr *ContentReader) SetIncludeHidden(include bool) {
	cr.includeHidden = include
}
