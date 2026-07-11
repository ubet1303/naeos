# Pipeline

Pipeline adalah rantai pemrosesan utama NAEOS yang mentransformasikan spesifikasi menjadi artifacts.

## Tahapan Pipeline

### 1. Parse
Mengurai YAML/JSON spesifikasi ke `SpecDocument`.

```bash
naeos run --input-file spec.yaml
```

**Output:** `parser.SpecDocument` — struct dengan Project, Modules, Services, Architecture, Deployment, Testing, Generation.

### 2. Normalize
Menormalisasi data spesifikasi ke format standar.

**Output:** `normalizer.NormalizedSpec`

### 3. Resolve
Resolve cross-references antar modul dan service.

**Output:** `resolver.ResolvedSpec`

### 4. Build NEIR
Membangun model NEIR dari spesifikasi yang sudah di-resolve.

**Output:** `model.NEIR` — model representasi intermediasi.

### 5. Validate
Memvalidasi NEIR:
- Module dependency validation
- Port conflict detection
- Architecture pattern validation
- Security rule validation

### 6. Schedule
Menjadwalkan tugas berdasarkan dependency graph.

**Output:** `[]scheduler.Task`

### 7. Generate
Menghasilkan artifacts untuk target bahasa/framework.

**Output:** `[]engine.Artifact`

### 8. Review
Review artifacts terhadap governance rules.

**Output:** `[]review.ReviewResult`

## Pipeline Execution

```go
// Basic usage
cfg := pipeline.Config{
    Name:      "my-pipeline",
    Mode:      "production",
    Verbose:   true,
    OutputDir: "./output",
    Languages: []string{"go", "typescript"},
}

p, err := pipeline.New(cfg)
result, err := p.Run(specInput)

// Dengan context
ctx := context.Background()
result, err := p.RunContext(ctx, specInput)

// Dry run
cfg.DryRun = true
result, err := p.Run(specInput)
```

## Hooks

Pipeline mendukung hooks untuk customisasi:

```go
cfg := pipeline.Config{
    Hooks: &pipeline.Hooks{
        BeforeParse:    []pipeline.HookFunc{logHook},
        AfterParse:     []pipeline.HookFunc{validateHook},
        BeforeRun:      []pipeline.HookFunc{preRunHook},
        AfterRun:       []pipeline.HookFunc{postRunHook},
        BeforeGenerate: []pipeline.HookFunc{preGenHook},
        AfterGenerate:  []pipeline.HookFunc{postGenHook},
    },
}
```

## Kernel Integration

Pipeline terintegrasi dengan kernel untuk:
- Service registry
- Event bus (pub/sub)
- Telemetry collection

```go
// Kernel metrics
metrics := p.KernelMetrics()

// Subscribe to events
p.Subscribe("pipeline.run", func(payload any) {
    // Handle event
})

// Publish events
p.Publish("custom.event", map[string]any{"key": "value"})
```

## Pipeline Stages Detail

| Stage | Input | Output | Error Handling |
|-------|-------|--------|----------------|
| Parse | string | SpecDocument | Return error |
| Normalize | SpecDocument | NormalizedSpec | Return error |
| Resolve | NormalizedSpec | ResolvedSpec | Return error |
| Build | ResolvedSpec | NEIR | Return error |
| Validate | NEIR | Validated NEIR | Return error |
| Schedule | NEIR | Tasks | Return error |
| Generate | NEIR | Artifacts | Return error |
| Review | Artifacts | Reviews | Log warnings |
| Write | Artifacts | Files on disk | Return error |
