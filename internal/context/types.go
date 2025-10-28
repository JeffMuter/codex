package context

import "time"

// Context represents all gathered context information
type Context struct {
	Timestamp    time.Time          `json:"timestamp"`
	Repository   *RepositoryContext `json:"repository,omitempty"`
	Filesystem   *FilesystemContext `json:"filesystem,omitempty"`
	NixConfig    *NixContext        `json:"nix_config,omitempty"`
	Dotfiles     *DotfilesContext   `json:"dotfiles,omitempty"`
	Screenshot   *Screenshot        `json:"screenshot,omitempty"`
}

// RepositoryContext contains git repository information
type RepositoryContext struct {
	Path          string   `json:"path"`
	Remote        string   `json:"remote,omitempty"`
	Branch        string   `json:"branch"`
	RecentCommits []Commit `json:"recent_commits,omitempty"`
	Status        string   `json:"status,omitempty"`
}

// Commit represents a git commit
type Commit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
}

// FilesystemContext contains current directory and file information
type FilesystemContext struct {
	CurrentDir string   `json:"current_dir"`
	ParentDirs []string `json:"parent_dirs,omitempty"`
	Files      []string `json:"files,omitempty"`
}

// NixContext contains parsed Nix configuration
type NixContext struct {
	ConfigPath       string              `json:"config_path"`
	IsFlake          bool                `json:"is_flake"`
	Packages         []string            `json:"packages,omitempty"`
	SystemConfig     map[string]any      `json:"system_config,omitempty"`
	LastParsed       time.Time           `json:"last_parsed"`
	CacheKey         string              `json:"cache_key"`
}

// DotfilesContext contains parsed dotfiles configuration
type DotfilesContext struct {
	DotfilesPath     string              `json:"dotfiles_path"`
	IsHomeManager    bool                `json:"is_home_manager"`
	Configs          map[string]any      `json:"configs,omitempty"`
	Keybindings      map[string][]Keybind `json:"keybindings,omitempty"`
	LastParsed       time.Time           `json:"last_parsed"`
	CacheKey         string              `json:"cache_key"`
}

// Keybind represents a keyboard binding
type Keybind struct {
	Key         string `json:"key"`
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
	Mode        string `json:"mode,omitempty"`
}

// Screenshot contains screenshot data
type Screenshot struct {
	Path     string `json:"path"`
	Data     []byte `json:"data,omitempty"`
	MimeType string `json:"mime_type"`
	Tool     string `json:"tool"` // Which tool was used to capture
}

// GatherOptions configures what context to gather
type GatherOptions struct {
	IncludeRepository bool
	IncludeFilesystem bool
	IncludeNixConfig  bool
	IncludeDotfiles   bool
	CaptureScreenshot bool
	WorkingDir        string // If empty, uses current directory
}
