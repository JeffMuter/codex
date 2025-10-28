# Codex Development TODO

## Current Status (as of 2025-10-28)
**What's Working:**
- ✓ Cobra CLI framework integrated
- ✓ `codex ask [question]` command with `--screenshot` and `--verbose` flags
- ✓ `codex config show` displays environment variables
- ✓ `codex config set [key] [value]` provides setup instructions
- ✓ All help text and command structure in place
- ✓ Binary builds successfully
- ✓ Internal package structure created (internal/config, internal/context, internal/providers, internal/errors)

**What's Next:**
- Choose and integrate logging framework (zerolog or zap)
- Update cmd/ files to use new internal packages
- Configuration file support implementation (~/.config/codex/config.yaml)
- Git repository detection and context gathering
- Database setup with goose migrations

---

## Phase 1: Core CLI Structure ✓ (Partially Complete)
- [x] Set up CLI framework (cobra)
  - [x] Implement `codex ask` command with argument parsing
  - [x] Implement `codex config` command group (show/set subcommands)
  - [x] Add `--screenshot` flag support (parsing ready, capture logic pending)
  - [x] Add `--verbose` flag for debugging
- [ ] Configuration management
  - [x] Read from environment variables (CODEX_NIX_CONFIG, CODEX_DOTFILES)
  - [ ] Support config file (~/.config/codex/config.yaml or similar)
  - [ ] Implement config validation
  - [x] Add `codex config show` to display current configuration
  - [ ] Implement `codex config set` persistence (currently shows env var instructions)
- [x] Basic project structure
  - [x] Create cmd/ package with root, ask, config commands
  - [x] Create internal/ package structure (config/, context/, providers/, errors/)
  - [ ] Create pkg/ if needed for reusable libraries
  - [x] Set up error handling patterns (internal/errors package)
  - [ ] Add logging framework (structured logging, e.g., zerolog or zap) - PENDING USER CHOICE

## Phase 2: Context Discovery & Analysis
- [ ] Repository detection
  - [ ] Implement git repository finder (traverse up from cwd)
  - [ ] Parse basic repo metadata (remote, branch, recent commits)
  - [ ] Handle non-git directories gracefully
- [ ] Nix configuration parser
  - [ ] Locate Nix config repository
  - [ ] Parse .nix files for installed packages
  - [ ] Extract system configuration
  - [ ] Handle flakes vs traditional Nix configs
  - [ ] Convert parsed data to JSON format for storage
- [ ] Dotfiles/Home-manager parser
  - [ ] Locate dotfiles repository
  - [ ] Parse common config files (tmux, neovim, shell, etc.)
  - [ ] Extract keybindings from various tools
  - [ ] Handle both plain dotfiles and home-manager setups
  - [ ] Convert parsed data to JSON format for storage
- [ ] Filesystem context
  - [ ] Capture current working directory
  - [ ] Build parent directory hierarchy
  - [ ] Include relevant files from current directory
  - [ ] Serialize as JSON

## Phase 3: Database & Caching (SQLite)
- [ ] Database setup
  - [ ] Set up goose migrations structure for SQLite
  - [ ] Create initial schema for context cache
  - [ ] Add tables for parsed configs with timestamps (store as JSON)
  - [ ] Add indexes for efficient lookups
  - [ ] Set SQLite pragmas for performance (WAL mode, etc.)
  - [ ] Handle database location (~/.local/share/codex/codex.db or similar)
- [ ] Context caching
  - [ ] Implement cache key generation (config repo + hash)
  - [ ] Cache parsed Nix configuration as JSON
  - [ ] Cache parsed dotfiles/home-manager as JSON
  - [ ] Add TTL or invalidation strategy (check file mtimes)
  - [ ] Implement cache refresh logic
  - [ ] Add `codex cache clear` command

## Phase 4: Screenshot Support
- [ ] Screenshot tool detection
  - [ ] Detect available screenshot tools (scrot, maim, grim, etc.)
  - [ ] Handle different display servers (X11, Wayland)
  - [ ] Add fallback options if no tool found
- [ ] Screenshot capture
  - [ ] Implement screenshot capture via shell command
  - [ ] Save to temporary location
  - [ ] Clean up temp files after use
  - [ ] Handle screenshot errors gracefully

## Phase 5: AI Provider Integration
- [x] Provider interface
  - [x] Define common interface for AI providers (internal/providers/provider.go)
  - [x] Implement provider factory/registry
  - [ ] Add configuration for provider selection
- [ ] Anthropic Claude integration
  - [x] Provider skeleton created (internal/providers/anthropic.go)
  - [ ] API client implementation
  - [ ] Handle API key from environment/config
  - [ ] Implement rate limiting
  - [ ] Handle errors and retries
- [ ] OpenAI integration
  - [x] Provider skeleton created (internal/providers/openai.go)
  - [ ] API client implementation
  - [ ] Support multiple models
  - [ ] Handle API key from environment/config
- [ ] Local LLM support
  - [x] Provider skeleton created (internal/providers/ollama.go)
  - [ ] Ollama integration
  - [ ] Support for other local providers
  - [ ] Handle local model availability checks
- [ ] Provider configuration
  - [ ] Add `codex config set-provider` command
  - [ ] Store provider preference
  - [ ] Allow per-query provider override

## Phase 6: Query Processing
- [ ] Context builder
  - [ ] Combine all context sources into structured JSON format
  - [ ] Calculate total context size (tokens/characters)
  - [ ] Optimize context size (summarize, prioritize relevant parts)
  - [ ] Handle screenshots (base64 encode, attach to request)
  - [ ] Create system prompt with context
- [ ] Context size validation
  - [ ] Define context size thresholds (warn at X tokens, error at Y)
  - [ ] Implement token estimation for different providers
  - [ ] Show warning with size estimate and cost estimate if large
  - [ ] Prompt user for confirmation before sending large context
  - [ ] Add `--force` flag to skip confirmation
  - [ ] Add `--no-context` flag to send query without context
- [ ] Query handler
  - [ ] Accept user query
  - [ ] Build complete prompt with context
  - [ ] Validate context size and get user confirmation if needed
  - [ ] Send to selected AI provider
  - [ ] Stream response to terminal
  - [ ] Handle errors and provide fallback

## Phase 7: Testing & Polish
- [ ] Unit tests
  - [ ] Test config parsing
  - [ ] Test repository detection
  - [ ] Test context building
  - [ ] Test provider interface
- [ ] Integration tests
  - [ ] Test full query flow
  - [ ] Test with sample config repositories
  - [ ] Test error scenarios
- [ ] Documentation
  - [ ] Add code comments
  - [ ] Update README with actual usage
  - [ ] Add example configurations
  - [ ] Document troubleshooting common issues
- [ ] Polish
  - [ ] Add progress indicators for slow operations
  - [ ] Improve error messages
  - [ ] Add shell completions
  - [ ] Optimize performance

## Phase 8: Advanced Features (Future)
- [ ] Interactive mode
  - [ ] Multi-turn conversations
  - [ ] Context persistence across queries
- [ ] Query history
  - [ ] Store past queries (optional, privacy-conscious)
  - [ ] Search through history
  - [ ] Rerun past queries
- [ ] Smart context selection
  - [ ] Analyze query to determine relevant context
  - [ ] Only include necessary config files
  - [ ] Reduce token usage
- [ ] Plugin system
  - [ ] Support custom context providers
  - [ ] Allow custom parsers for config formats

## Questions & Decisions Made
- **AI Backend**: Multiple providers (Anthropic, OpenAI, Local LLM)
- **Priority**: Core CLI first, then context gathering
- **Database**: SQLite for context caching (embedded, no server required)
- **Config Format**: JSON for storing parsed configurations
- **Context Size**: Warn user and require confirmation before sending large context
- **Screenshots**: Shell out to existing tools (scrot, maim, grim)

## Recent Commits & Progress
- **[uncommitted]** (2025-10-28): Created internal package structure
  - Created internal/config/config.go with full config management (load, save, validate)
  - Created internal/context/types.go with context data structures
  - Created internal/context/gatherer.go with context gathering interface
  - Created internal/providers/provider.go with AI provider interface
  - Created internal/providers/anthropic.go, openai.go, ollama.go with provider skeletons
  - Created internal/errors/errors.go with standardized error handling
  - Added support for YAML config files, environment variables, and multiple AI providers
- **b2b1991** (2025-10-28): Added Cobra framework and implemented core CLI structure
  - Created cmd/root.go, cmd/ask.go, cmd/config.go
  - Implemented `ask`, `config show`, and `config set` commands
  - Added --screenshot and --verbose flags
  - Updated main.go to use Cobra's Execute()
- **89e4753**: Initial main func and TODO setup
- **a6e2d4d**: First push with project structure

## Notes
- Project uses Go 1.24.4 with Nix development environment
- Database migrations via goose (SQLite backend) - not yet implemented
- Focus on being context-aware without being intrusive
- Privacy-conscious design (local processing where possible)
- All parsed context will be stored as JSON in SQLite
- User confirmation required for large context to avoid unexpected API costs
- Cobra provides auto-generated shell completion support
