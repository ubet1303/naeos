# Roadmap

## Fase 1: Fondasi ✅

- [x] Refine core documents
- [x] Ensure terminology consistency
- [x] Add contribution/onboarding guides
- [x] Core pipeline implementation
- [x] NEIR model
- [x] CLI with 30+ commands

## Fase 2: Tooling dan Validasi ✅

- [x] ADR/RFC templates
- [x] Review mechanisms
- [x] Document validation rules
- [x] Policy evaluator
- [x] Governance system

## Fase 3: Referensi Implementasi

- [ ] Reference implementations
- [ ] Requirement-to-deployment workflow
- [ ] Industry-specific profiles expansion
- [ ] Advanced migration tools

## Fase 4: Ekosistem

- [ ] AI agent/toolchain interoperability
- [ ] Public documentation portal
- [ ] Cross-organization adoption
- [ ] Plugin marketplace

## Completed Technical Implementation

### v0.1.0 — Foundation
- Core pipeline: parser, normalizer, resolver, NEIR builder, validator
- Planner: DAG graph with topological sort and cycle detection
- Generator engine: Go project code, Dockerfile, CI, documentation
- Policy evaluator with 7 operators and 5 default rules
- Artifact reviewer with governance rules
- Knowledge graph with 14 node types and 13 edge types
- Runtime execution engine with deduplication
- Telemetry event collector

### v0.2.0 — Compiler Foundation
- Compiler transforms NEIR into AI instruction sets
- 6 output adapters: Copilot, Claude Code, Cursor, Gemini CLI, Codex, OpenCode
- Artifact Store with content-hash dedup
- 5 industry profiles: SaaS, AI Agent, FinTech, Healthcare, Government
- CLI: compile, profile, artifacts, migrate
- CI/CD: GitHub Actions, CodeQL

### v0.3.0 — Core Specification
- Spec Language v2: variable interpolation, env resolution, reference resolution
- Validation Kernel: circular dependencies, port conflicts, module boundaries
- Schema versioning: auto-check spec version on parse
- AI Context Bundles: generate LLM-ready context from specs
- CLI: context command
