# Codex Development TODO

## 🎉 MVP COMPLETE! 🎉

**What Works:**
- ✅ Cobra CLI framework with commands: `ask`, `config show/add-repo/list-repos/remove-repo`
- ✅ Configuration management (YAML + .env support)
- ✅ Anthropic API integration with streaming responses
- ✅ Context gathering (filesystem + repository fetching)
- ✅ End-to-end query flow: `codex ask "question"` → Claude API → streamed answer
- ✅ Repository management (local and remote with caching)
- ✅ Logging framework (zerolog)
- ✅ Privacy-conscious design (opt-in current repo with `--current-repo` flag)

**MVP Demo:**
```bash
codex ask "what is 2 + 2?"
# → Loads API key from .env
# → Gathers context (filesystem, configured repos)
# → Sends to Claude API
# → Streams response to terminal
```

---

## Future Features (Post-MVP)

### Enhancement: Repository Content Reading
**Priority: High** - Currently repositories are detected but file contents are not read
- Implement actual file reading in RepoFetcher
- Parse common config files (.zshrc, .bashrc, .nix files, etc.)
- Format repository contents into context string for AI
- See internal/context/gatherer.go:98-122 and repo_fetcher.go:26-32

### Testing & Quality
- Unit tests for config loading
- Unit tests for context gathering
- Integration tests for end-to-end flow
- Error handling edge cases

### Caching & Performance
- SQLite database for parsed config caching
- TTL and cache invalidation strategies
- Context size optimization and chunking
- Token counting and warnings for large contexts

### Additional AI Providers
- OpenAI integration
- Local LLM support (Ollama)
- Provider switching mechanism
- Model selection per provider

### Advanced Context Features
- Screenshot support for visual context
- Deep Nix flake parsing
- Home-manager config analysis
- Dotfiles smart parsing (patterns in common config files)
- User confirmation for large contexts (API cost awareness)

### User Experience Polish
- Progress indicators during context gathering
- Shell completions (Cobra auto-generation)
- Better error messages and diagnostics
- Interactive mode with multi-turn conversations
- Context persistence across queries

### Extended Features
- Query history storage and search
- Plugin system for custom context providers
- Custom config parsers
- Export/import configuration
- Template system for common queries

---

## Architecture Decisions (Established)
- **AI Backend**: Anthropic Claude (MVP complete), will expand to other providers
- **Philosophy**: Prove core value first, optimize later
- **Database**: No caching in MVP - parse on-demand, SQLite planned for future
- **Screenshots**: Deferred to post-MVP features
- **Logging**: zerolog with structured logging
- **Repository Model**: Configured repos from config file + opt-in current repo
- **Privacy**: Explicit user control over what context is shared

## Development History

### MVP Development (2025-10-28)
Successfully built working CLI tool with core functionality:

**Key Milestones:**
1. **Infrastructure Setup**
   - Cobra CLI framework with command structure
   - Internal package architecture (config, context, providers, errors)
   - Zerolog logging with verbose mode support
   - Nix development environment with Go 1.23

2. **Configuration System**
   - YAML config file support (~/.config/codex/config.yaml)
   - .env file loading for API keys (godotenv)
   - Environment variable overrides
   - Repository management commands (add-repo, list-repos, remove-repo)

3. **Context Gathering**
   - Filesystem context collection
   - Repository fetching (local and remote with caching)
   - Privacy-conscious opt-in model for current repo (--current-repo flag)
   - ConfiguredRepo type distinguishing local/remote sources

4. **AI Integration** ✅
   - Full Anthropic SDK integration
   - Streaming response support
   - Context-to-prompt formatting
   - claude-3-opus-20240229 model

5. **Working Demo**
   - `codex ask "question"` end-to-end flow
   - Successful manual integration testing
   - Real-time streamed responses

**Commits:**
- b2b1991: Added Cobra framework and core CLI structure
- (Multiple uncommitted): Full MVP implementation and testing

## Technical Notes
- **Language**: Go 1.23 with Nix development environment
- **Dependencies**: Cobra (CLI), zerolog (logging), godotenv (env), Anthropic SDK
- **Database**: goose migrations ready, SQLite backend planned for caching (not yet implemented)
- **Philosophy**: Context-aware without being intrusive, privacy-first design
- **Future**: Context caching as JSON in SQLite, token counting for API cost awareness
