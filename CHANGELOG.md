# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-07-11

### Added
- **Compiler Foundation** (`internal/compiler/`): Transforms NEIR into AI instruction sets for 6 target tools.
- **AI Output Adapters** (`internal/compiler/adapters/`):
  - GitHub Copilot — `.github/copilot-instructions.md`, `.github/copilot-context.md`, `.github/copilot-rules.md`
  - Claude Code — `CLAUDE.md`, `.claude/context.md`, `.claude/rules.md`
  - Cursor — `.cursorrules`, `.cursor/context.md`
  - Gemini CLI — `.gemini/CONFIG.md`, `.gemini/context.md`
  - Codex — `AGENTS.md`, `.codex/context.md`
  - OpenCode — `AGENTS.md`, `.opencode/context.md`, `.opencode/rules.md`
- **Artifact Store** (`internal/artifacts/`): Manages generated outputs with content-hash dedup, kind detection, metadata, and disk persistence.
- **Profile Registry** (`internal/profiles/`): 5 industry-specific profiles (SaaS, AI Agent, FinTech, Healthcare, Government) with modules, services, architecture, security, deployment, and testing templates.
- **Migration constants**: `CurrentVersion` (0.1.0) and `TargetVersion` (0.3.0) exported for version-aware tooling.
- **CLI commands**:
  - `naeos compile` — compile spec into AI instruction sets (per-target or `--all`)
  - `naeos profile list|show|search|apply` — browse and apply industry profiles
  - `naeos artifacts list|info|dedup|summary` — manage generated artifact store
  - `naeos migrate run|plan|versions` — manage schema migrations with dry-run support
- Comprehensive test suites: compiler (6 tests), adapters (14 tests), artifacts (14 tests), profiles (9 tests)

### Changed
- All 63 packages pass, `go vet` clean, `go build` clean.

### Added
- Documentation index with recommended reading orders (beginner, policy, profile, CLI, testing).
- NES-028 CLI Reference — comprehensive CLI command documentation.
- NES-029 Configuration — pipeline configuration reference.
- NES-030 Specification Language — NAEOS specification language docs.
- NES-031 Errors — exhaustive error catalog.
- NES-032 Telemetry — telemetry and metrics reference.
- NES-033 Testing Guide — test guide with coverage requirements.
- NES-034 Event Bus — internal pub/sub event bus documentation.
- NES-035 Version Management — SemVer management documentation.
- NES-036 Template Renderer — template rendering engine documentation.
- NES-037 Knowledge Graph & Provenance — knowledge graph and lineage documentation.
- NES-038 Shared Types & Contracts — shared types and contracts documentation.
- NAEOS-GOV-002 Vision — long-term vision document.
- NAEOS-GOV-005 Core Principles — 8 core engineering principles.
- Expanded 18 NES stub documents (NES-003 through NES-022) with full API references and examples.
- `status` command — display current pipeline and project status.
- Auto-detection of config files (`config.yaml`, `config.yml`, `config.json`, `naeos.yaml`, `naeos.yml`, `naeos.json`, `.naeos/config.*`) in working directory.
- Global `--dry-run` flag for preview mode across all commands.
- Per-command `--dry-run` flag for `run`, `export`, and `preview` commands.
- Language-aware scaffold — `--language` flag now generates correct files for Go, TypeScript, Python, Java, and Rust.
- E2E test suite with comprehensive pipeline integration tests.
- Additional benchmarks for dry-run, full-spec, and verbose pipeline runs.
- Fixed GoAdapter `cleanModulePath` to correctly handle relative paths (e.g., `./internal/core`).

### Changed
- NES-001 Repository — updated repository structure to match actual codebase paths.
- DOCUMENTATION-INDEX.md — added NES-028 through NES-038, Go package reference section, CLI and testing reading orders.
- **Refactored `cmd/naeos/main.go`**: split 1876-line monolith into 28 separate command files for better maintainability.
- All CLI commands now support `--config` auto-detection (no longer required to specify explicitly).
- Improved CLI help text with usage examples for all commands.
- Pipeline `Config` struct now includes `DryRun` field for preview mode.
- `preview` command now uses dry-run mode by default.
- Removed unused `hashContent()` function from CLI.
- Consistent error handling across all CLI commands.
- Go adapter `GenerateProject` now generates a complete runnable main.go with HTTP server setup, health check, and API endpoints.

## [0.1.0] - 2026-01-01

### Added
- Initial project structure.
- CLI with 11 subcommands: init, run, validate, inspect, doctor, repair, scaffold, export, preview, kernel, version.
- Core pipeline: parser, normalizer, resolver, NEIR builder, validator.
- Planner: DAG graph with topological sort and cycle detection.
- Generator engine: Go project code, Dockerfile, CI, documentation.
- Policy evaluator with 7 operators and 5 default rules.
- Artifact reviewer with governance rules.
- Knowledge graph with 14 node types and 13 edge types.
- Provenance tracking store.
- Runtime execution engine with deduplication.
- Telemetry event collector.
- 34 modular design documents (NES-000 through NES-033).
- 10 specification documents (NAEOS-SPEC-001 through 010).
- 8 constitutional documents (NAEOS-CON-001 through 008).
- 8 governance documents (NAEOS-GOV-001 through 008).
- 4 kernel specification documents (NAEOS-KER-001 through 004).
- 7 policy system documents (NAEOS-POL-001 through 007).
- 7 profile system documents (NAEOS-PRO-001 through 007).
- 1 reference architecture document (NAEOS-NRA-001).
- ADR and RFC templates with examples.
- Example specifications (minimal and full).
