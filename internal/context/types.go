package context

import "time"

// Context represents all gathered context information
type Context struct {
	Timestamp         time.Time            `json:"timestamp"`
	ConfiguredRepos   []*RepositoryContext `json:"configured_repos,omitempty"`
	CurrentRepo       *RepositoryContext   `json:"current_repo,omitempty"`
	Filesystem        *FilesystemContext   `json:"filesystem,omitempty"`
	NixConfig         *NixContext          `json:"nix_config,omitempty"`
	Dotfiles          *DotfilesContext     `json:"dotfiles,omitempty"`
	Screenshot        *Screenshot          `json:"screenshot,omitempty"`
}

// RepositoryContext contains git repository information
type RepositoryContext struct {
	Path   string `json:"path"`
	Remote string `json:"remote,omitempty"`
	Source string `json:"source,omitempty"` // Original source (URL or path)
	Type   string `json:"type,omitempty"`   // "local" or "remote" or "current"
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
	IncludeCurrentRepo bool   // Include current working directory repo (opt-in)
	IncludeFilesystem  bool
	IncludeNixConfig   bool
	IncludeDotfiles    bool
	CaptureScreenshot  bool
	WorkingDir         string // If empty, uses current directory
}
