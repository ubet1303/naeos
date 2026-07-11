# Project Structure

Struktur repositori NAEOS.

## Root Level

| File/Dir | Deskripsi |
|----------|-----------|
| `cmd/naeos/` | CLI commands (30+ files) |
| `internal/` | Internal packages |
| `pkg/` | Public packages |
| `docs/` | Documentation (40+ files) |
| `governance/` | Governance documents |
| `constitution/` | Engineering constitution |
| `specification/` | Specification documents |
| `wiki/` | Project wiki |
| `templates/` | ADR/RFC templates |
| `examples/` | Example documents |
| `go.mod` | Go module definition |
| `go.sum` | Go module checksum |
| `README.md` | Project readme |
| `CHANGELOG.md` | Version history |
| `LICENSE` | Apache 2.0 license |
| `.github/` | GitHub workflows |
| `.gitignore` | Git ignore rules |

## cmd/naeos/

CLI commands organized by feature:

```
cmd/naeos/
├── main.go              # Root command
├── run_cmd.go           # Pipeline execution
├── validate_cmd.go      # Spec validation
├── lint_cmd.go          # Spec linting
├── create_cmd.go        # Project creation
├── scaffold_cmd.go      # Code scaffolding
├── export_cmd.go        # Artifact export
├── compile_cmd.go       # AI compilation
├── context_cmd.go       # Context bundles
├── profile_cmd.go       # Industry profiles
├── artifacts_cmd.go     # Artifact management
├── migrate_cmd.go       # Schema migration
├── doctor_cmd.go        # System health
├── diff_cmd.go          # Spec comparison
├── watch_cmd.go         # File watching
├── workspace_cmd.go     # Workspace management
├── kernel_cmd.go        # Kernel inspection
├── plugin_cmd.go        # Plugin management
├── template_cmd.go      # Template management
├── lock_cmd.go          # Dependency locking
├── rollback_cmd.go      # Rollback changes
├── audit_cmd.go         # Spec auditing
├── marketplace_cmd.go   # Profile marketplace
├── status_cmd.go        # Pipeline status
├── docs_cmd.go          # Documentation
├── ai_cmd.go            # AI assistance
├── init_cmd.go          # Init config
├── version_cmd.go       # Version info
├── completion_cmd.go    # Shell completion
├── repair_cmd.go        # Spec repair
├── preview_cmd.go       # Preview mode
├── helpers.go           # Shared helpers
├── e2e_test.go          # E2E tests
└── main_test.go         # Unit tests
```

## internal/

```
internal/
├── specification/       # Spec processing
│   ├── parser/          # YAML/JSON parser
│   ├── normalizer/      # Data normalization
│   └── resolver/        # Cross-reference resolution
├── neir/               # NEIR model
│   ├── model/          # Model definitions
│   ├── builder/        # NEIR builder
│   └── validator/      # NEIR validator
├── compiler/           # AI instruction compiler
│   └── adapters/       # Output adapters
├── context/            # Context bundles
│   └── bundle/         # Bundle generator
├── generation/         # Code generation
│   ├── engine/         # Generation engine
│   ├── adapters/       # Language adapters
│   └── renderers/      # Template renderers
├── governance/         # Governance
│   ├── policy/         # Policy evaluator
│   └── review/         # Artifact review
├── artifacts/          # Artifact store
├── profiles/           # Industry profiles
├── migration/          # Schema migration
├── security/           # Security rules
├── events/             # Event bus
├── knowledge/          # Knowledge graph
├── planner/            # Task scheduling
├── runtime/            # Runtime engine
├── lock/               # Dependency locking
├── rollback/           # Rollback management
├── marketplace/        # Profile marketplace
├── templates/          # Template engine
├── workspace/          # Workspace management
├── watch/              # File watcher
├── profiling/          # Performance profiling
├── registry/           # Service registry
├── ai/                 # AI assistance
├── docs/               # Documentation engine
├── diff/               # Diff engine
├── lint/               # Lint rules
├── create/             # Project creation
└── shared/             # Shared utilities
    ├── log/            # Logging
    ├── strutil/        # String utilities
    └── contracts/      # Shared contracts
```

## pkg/

```
pkg/
├── pipeline/           # Main pipeline
├── kernel/             # System kernel
├── config/             # Configuration
└── plugin/             # Plugin system
```

## File Count

| Directory | Files |
|-----------|-------|
| cmd/naeos/ | 35 |
| internal/ | 60+ |
| pkg/ | 4 |
| docs/ | 44 |
| governance/ | 8 |
| Total | 150+ |
