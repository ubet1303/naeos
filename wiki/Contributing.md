# Contributing

Panduan untuk berkontribusi ke NAEOS.

## Development Setup

```bash
# Clone repo
git clone https://github.com/NAEOS-foundation/naeos.git
cd naeos

# Build
go build ./cmd/naeos/

# Test
go test ./...

# Lint
go vet ./...
```

## Project Structure

```
naeos/
├── cmd/naeos/          # CLI commands
│   ├── main.go         # Root command
│   ├── run_cmd.go      # Run pipeline
│   ├── compile_cmd.go  # Compile to AI tools
│   ├── context_cmd.go  # Generate context bundles
│   └── ...             # 30+ command files
├── internal/
│   ├── specification/  # Parser, normalizer, resolver
│   ├── neir/          # NEIR model and builder
│   ├── compiler/      # AI instruction compiler
│   ├── context/       # Context bundle generator
│   ├── generation/    # Code generation
│   ├── governance/    # Policy and review
│   ├── artifacts/     # Artifact store
│   ├── profiles/      # Industry profiles
│   ├── migration/     # Schema migration
│   ├── security/      # Security rules
│   └── shared/        # Shared utilities
├── pkg/
│   ├── pipeline/      # Main pipeline
│   ├── kernel/        # System kernel
│   ├── config/        # Configuration
│   └── plugin/        # Plugin system
├── docs/              # Documentation (40+ files)
├── governance/        # Governance documents
├── constitution/      # Engineering constitution
├── specification/     # Specification documents
└── wiki/              # This wiki
```

## Coding Standards

### Go Style
- Follow `gofmt` and `go vet`
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

### Testing
- Write unit tests for new features
- Maintain test coverage
- Use table-driven tests
- Mock external dependencies

### Documentation
- Update documentation for new features
- Add examples for CLI commands
- Keep CHANGELOG updated

## Pull Request Process

1. Fork the repository
2. Create feature branch (`git checkout -b feature/my-feature`)
3. Write tests
4. Run `go test ./...` and `go vet ./...`
5. Update documentation
6. Commit with descriptive message
7. Push and create PR

## Commit Messages

```
feat: add new feature
fix: fix bug in parser
docs: update documentation
test: add unit tests
refactor: refactor code
chore: maintenance tasks
```

## Code Review

- All PRs require review
- Tests must pass
- Documentation must be updated
- Follow coding standards
