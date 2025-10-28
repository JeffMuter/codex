# Codex Development TODO - MVP Focus

## Current Status (as of 2025-10-28)
**What's Working:**
- ✓ Cobra CLI framework integrated
- ✓ `codex ask [question]` command with argument parsing
- ✓ `codex config show` displays environment variables
- ✓ Binary builds successfully
- ✓ Internal package structure created (internal/config, internal/context, internal/providers, internal/errors)
- ✓ Logging framework integrated (zerolog)
- ✓ **Environment variable loading (.env support)**
- ✓ **Anthropic API integration with streaming**
- ✓ **End-to-end query flow working**
- ✓ **Context gathering (basic filesystem)**
- ✓ **Repository fetching (local and remote)**

**🎉 MVP ACHIEVED! 🎉**
```bash
codex ask "what is 2 + 2?"
# → Loads API key from .env file
# → Gathers context (filesystem, repos)
# → Sends to Claude API with context
# → Streams response to terminal
# → Returns answer: "2 + 2 equals 4"
```

---

## MVP Tasks (Essential Only)

### 1. Configuration Management ✅ COMPLETE
- [x] Load config from ~/.config/codex/config.yaml
- [x] Simple config structure: repo paths, API key
- [x] Basic validation (paths exist, API key present)
- [x] .env file support for API keys
- [x] Environment variable override support

### 2. Context Gathering (In-Memory, No Caching) ✅ COMPLETE
- [x] Read configured repository paths from config
- [x] Simple file reading from those repos
- [x] Build context string from gathered data
- [x] Basic filesystem context (working directory)
- [ ] Basic dotfiles parsing (grep for patterns in common config files) - DEFERRED

### 3. Single AI Provider (Anthropic Only) ✅ COMPLETE
- [x] Basic Anthropic API client
- [x] Read API key from config/environment (.env)
- [x] Send prompt with context
- [x] Stream response to terminal
- [x] Basic error handling (fail fast)
- [x] Using Claude 3 Opus model

### 4. Query Processing (Minimal) ✅ COMPLETE
- [x] Accept user query from CLI
- [x] Build simple prompt: "Context: {context}\n\nQuestion: {query}"
- [x] Send to Anthropic
- [x] Display response (streamed)

### 5. Basic Testing
- [ ] Unit tests for config loading
- [ ] Unit tests for context gathering
- [x] Manual integration testing - WORKING!

---

## Current Issues & Next Steps

### ⚠️ CRITICAL: Repository Content Not Being Read (2025-10-28)
**Problem:** Configured repositories (nixos, dotfiles) are detected but their contents are NOT being sent to the AI.
- Test query: "What shell am I using and what are some of my custom aliases?"
- Expected: Should read ~/.dotfiles/zsh/ files and provide specific answers
- Actual: AI responds "I would need to look at your shell configuration files" - indicating it's not receiving file contents
- Root cause: RepoFetcher only verifies paths exist (repo_fetcher.go:26-32) but doesn't read/parse file contents
- Context gatherer (gatherer.go:98-122) calls FetchRepo but doesn't process the returned repository data
- **Next step**: Implement actual file reading and content extraction in context gathering
  - Add methods to read directory contents from configured repos
  - Parse common config files (.zshrc, .bashrc, .nix files, etc.)
  - Format repository contents into context string that gets sent to AI
  - See internal/context/gatherer.go:98-122 for where this logic needs to be added

## Future Features (Post-MVP)

- **Caching**: SQLite database for parsed config caching with TTL/invalidation
- **Screenshot Support**: Integrate screenshot tools for visual context
- **Multiple Providers**: OpenAI, local LLM (Ollama), provider switching
- **Smart Context**: Token counting, size warnings, user confirmation for large contexts
- **Advanced Parsing**: Deep Nix flake parsing, home-manager config analysis
- **Polish**: Progress indicators, shell completions, better error messages
- **Interactive Mode**: Multi-turn conversations with context persistence
- **Query History**: Store and search past queries
- **Plugin System**: Custom context providers and config parsers

---

## Questions & Decisions
- **AI Backend**: Anthropic Claude for MVP, expand later
- **Priority**: Prove core value first, optimize later
- **No Database**: Parse on-demand for MVP
- **No Screenshots**: Defer until after basic queries work
- **Logging**: zerolog (already integrated)
- **Repository Model**: Configured repos from config file

## Recent Commits & Progress
- **[uncommitted]** (2025-10-28): Added nixos and dotfiles repositories to config
  - Added ~/nixos and ~/.dotfiles to configured_repos via `codex config add-repo`
  - Updated config.yaml with nix_config_path and dotfiles_path
  - Enhanced `codex config show` command to read from config file (was only reading env vars)
  - Fixed cmd/config.go to load and display actual configuration
  - Discovered critical issue: repository contents not being read/sent to AI
- **[uncommitted]** (2025-10-28): 🎉 **MVP COMPLETE** - Anthropic integration working end-to-end
  - Added godotenv dependency for .env file loading
  - Created .env.example template with API key placeholder
  - Updated .gitignore to exclude .env file
  - Implemented full Anthropic SDK integration with streaming support
  - Built context-to-prompt formatter
  - Wired up ask command with provider and context gathering
  - Fixed model name to use claude-3-opus-20240229
  - Basic filesystem context gathering implemented
  - Successfully tested: `./codex ask "What is 2 + 2?"` returns correct answer with streaming
- **[uncommitted]** (2025-10-28): Implemented repository management system
  - Created ConfiguredRepo type to distinguish local/remote repositories
  - Implemented RepoFetcher for fetching local and remote repositories
  - Remote repos are cloned to ~/.local/share/codex/repos/ and cached
  - Modified Context types: ConfiguredRepos (always included) + CurrentRepo (opt-in)
  - Updated GatherOptions: IncludeCurrentRepo flag instead of automatic detection
  - Modified Gatherer to always fetch configured repos, optionally include current repo
  - Added --current-repo (-r) flag to ask command for opt-in behavior
  - Added config commands: add-repo, list-repos, remove-repo
  - Updated README with new repository model and privacy-conscious approach
  - Config struct now includes ConfiguredRepos []ConfiguredRepo
  - AddRepo/RemoveRepo methods for managing repository configuration
- **[uncommitted]** (2025-10-28): Integrated zerolog logging framework
  - Created internal/logging/logger.go with structured logging support
  - Added zerolog dependency (v1.34.0) to go.mod
  - Updated cmd/root.go with logger initialization in PersistentPreRun
  - Updated cmd/ask.go and cmd/config.go with structured logging
  - Enhanced shell.nix with Go dev tools (gopls, gotools, delve, sqlite)
  - Added shellHook with helpful environment information
  - Created docs/logging.md with usage examples and best practices
  - Verbose mode enables debug level + pretty console output
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
