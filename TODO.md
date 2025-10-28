# Codex Development TODO - MVP Focus

## Current Status (as of 2025-10-28)
**What's Working:**
- ✓ Cobra CLI framework integrated
- ✓ `codex ask [question]` command with argument parsing
- ✓ `codex config show` displays environment variables
- ✓ Binary builds successfully
- ✓ Internal package structure created (internal/config, internal/context, internal/providers, internal/errors)
- ✓ Logging framework integrated (zerolog)

**MVP Goal:**
```bash
codex ask "what's my tmux prefix?"
# → Reads config for repo paths
# → Parses those repos for tmux config
# → Sends to Claude API with simple prompt
# → Returns answer
```

---

## MVP Tasks (Essential Only)

### 1. Configuration Management
- [ ] Load config from ~/.config/codex/config.yaml
- [ ] Simple config structure: repo paths, API key
- [ ] Basic validation (paths exist, API key present)

### 2. Context Gathering (In-Memory, No Caching)
- [ ] Read configured repository paths from config
- [ ] Simple file reading from those repos
- [ ] Basic dotfiles parsing (grep for patterns in common config files)
- [ ] Build context string from gathered data

### 3. Single AI Provider (Anthropic Only)
- [ ] Basic Anthropic API client
- [ ] Read API key from config/environment
- [ ] Send prompt with context
- [ ] Stream response to terminal
- [ ] Basic error handling (fail fast)

### 4. Query Processing (Minimal)
- [ ] Accept user query from CLI
- [ ] Build simple prompt: "Context: {context}\n\nQuestion: {query}"
- [ ] Send to Anthropic
- [ ] Display response

### 5. Basic Testing
- [ ] Unit tests for config loading
- [ ] Unit tests for context gathering
- [ ] Manual integration testing

---

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
