# Compiler

Compiler mentransformasikan NEIR model ke AI instruction sets untuk berbagai target tool.

## Target

| Target | Output Files |
|--------|-------------|
| **Copilot** | `.github/copilot-instructions.md`, `.github/copilot-context.md`, `.github/copilot-rules.md` |
| **Claude Code** | `CLAUDE.md`, `.claude/context.md`, `.claude/rules.md` |
| **Cursor** | `.cursorrules`, `.cursor/context.md` |
| **Gemini CLI** | `.gemini/CONFIG.md`, `.gemini/context.md` |
| **Codex** | `AGENTS.md`, `.codex/context.md` |
| **OpenCode** | `AGENTS.md`, `.opencode/context.md`, `.opencode/rules.md` |

## Penggunaan CLI

```bash
# Compile ke target tertentu
naeos compile --target copilot --input-file spec.yaml

# Compile ke semua target
naeos compile --all --input-file spec.yaml

# Simpan output ke direktori
naeos compile --target claude --input-file spec.yaml --output ./ai-config
```

## Penggunaan Go

```go
import (
    "github.com/NAEOS-foundation/naeos/internal/compiler"
    "github.com/NAEOS-foundation/naeos/internal/compiler/adapters"
)

// Buat compiler
c := compiler.New()

// Register adapter
c.Register(adapters.NewCopilotAdapter())
c.Register(adapters.NewClaudeAdapter())

// Compile
output, err := c.Compile(neir, compiler.TargetCopilot)
if err != nil {
    log.Fatal(err)
}

// Output files
for _, file := range output.Files {
    fmt.Printf("%s (%s)\n", file.Path, file.Kind)
    fmt.Println(file.Content)
}
```

## Compile All

```go
outputs, err := c.CompileAll(neir)
for _, output := range outputs {
    fmt.Printf("Target: %s, Files: %d\n", output.Target, len(output.Files))
}
```

## Output Structure

```go
type CompiledOutput struct {
    Target      Target
    Files       []OutputFile
    Summary     string
    CompiledAt  time.Time
    NEIRVersion string
}

type OutputFile struct {
    Path    string
    Content string
    Kind    string  // "instructions", "context", "rules", "config"
}
```

## Adapter Interface

```go
type Adapter interface {
    Target() Target
    Compile(neir *model.NEIR) (*CompiledOutput, error)
}
```

### Custom Adapter

```go
type MyAdapter struct{}

func (a MyAdapter) Target() compiler.Target {
    return "my-tool"
}

func (a MyAdapter) Compile(neir *model.NEIR) (*compiler.CompiledOutput, error) {
    // Custom compilation logic
    return &compiler.CompiledOutput{
        Target: a.Target(),
        Files:  []compiler.OutputFile{...},
    }, nil
}
```
