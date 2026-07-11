# Architecture

## Gambaran Sistem

```
┌─────────────────────────────────────────────────────────┐
│                    NAEOS Architecture                     │
├─────────────┬──────────────┬──────────────┬─────────────┤
│    Input    │  Core Layer  │  Generation  │   Output    │
├─────────────┼──────────────┼──────────────┼─────────────┤
│  Spec YAML  │   Parser     │   Generator  │  Code Files │
│  CLI cmds   │   Normalizer │   Adapters   │  Configs    │
│  Profiles   │   Resolver   │   Renderers  │  Docs       │
│  Context    │   Validator  │   Compiler   │  AI Context │
│             │   Scheduler  │   Profiles   │  Artifacts  │
│             │   Kernel     │              │             │
│             │   Policy     │              │             │
│             │   Review     │              │             │
└─────────────┴──────────────┴──────────────┴─────────────┘
```

## Layer Arsitektur

### 1. Input Layer
Menerima spesifikasi dari berbagai sumber:
- **Spec YAML/JSON** — dokumen spesifikasi utama
- **CLI Commands** — perintah interaktif
- **Profiles** — template industri
- **Context Bundles** — bundle untuk AI

### 2. Core Layer
Memproses spesifikasi melalui pipeline:
- **Parser** — parsing YAML/JSON ke struct
- **Normalizer** — normalisasi data
- **Resolver** — resolve cross-references
- **Validator** — validasi NEIR
- **Scheduler** — scheduling tugas
- **Kernel** — service registry, event bus
- **Policy** — evaluasi kebijakan
- **Review** — governance review

### 3. Generation Layer
Menghasilkan output untuk berbagai target:
- **Generator** — generasi kode multi-bahasa
- **Adapters** — output untuk AI tools (Copilot, Claude, Cursor, Gemini, Codex, OpenCode)
- **Renderers** — rendering template
- **Compiler** — kompilasi ke AI instruction sets
- **Profiles** — template industri

### 4. Output Layer
Output final:
- **Code Files** — kode sumber
- **Config Files** — konfigurasi
- **Docs** — dokumentasi
- **AI Context** — bundle konteks AI
- **Artifacts** — semua output yang dihasilkan

## Data Flow

```
Spec YAML
    │
    ▼
┌─────────┐
│ Parser  │ → SpecDocument
└────┬────┘
     │
     ▼
┌────────────┐
│ Normalizer │ → NormalizedSpec
└────┬───────┘
     │
     ▼
┌──────────┐
│ Resolver │ → ResolvedSpec
└────┬─────┘
     │
     ▼
┌──────────┐
│ Builder  │ → NEIR Model
└────┬─────┘
     │
     ▼
┌───────────┐
│ Validator │ → Validated NEIR
└────┬──────┘
     │
     ▼
┌───────────┐
│ Scheduler │ → Task Graph
└────┬──────┘
     │
     ▼
┌───────────┐
│ Generator │ → Artifacts
└────┬──────┘
     │
     ▼
┌───────────┐
│ Compiler  │ → AI Instruction Sets
└───────────┘
```

## Technology Stack

| Komponen | Teknologi |
|----------|-----------|
| Bahasa | Go 1.22+ |
| CLI Framework | Cobra |
| Serialization | YAML, JSON |
| Config | YAML/JSON auto-detect |
| File Watcher | fsnotify |
| Logging | slog (structured) |
| Testing | go test, race detector |
| CI/CD | GitHub Actions |
| Security | CodeQL |
| License | Apache 2.0 |
