# Context Bundles

Context Bundles menghasilkan ringkasan proyek dalam format yang optimal untuk konsumsi LLM.

## Penggunaan CLI

```bash
# Generate markdown
naeos context --input-file spec.yaml

# Generate plain text
naeos context --input-file spec.yaml --output plain

# Generate JSON
naeos context --input-file spec.yaml --output json

# Simpan ke file
naeos context --input-file spec.yaml --output-file context.md
```

## Output Format

### Markdown

```markdown
# my-app — AI Context Bundle

## Summary
Project: my-app; Modules: auth, api; Services: 1

## Modules

- **auth** (`./auth`)
  Dependencies: core
- **api** (`./api`)
  Dependencies: auth

## Services

- **gateway** (kind=http, port=8080)
  - GET /users → listUsers
  - POST /users → createUser

## Languages
go, typescript
```

### Plain Text

```
Project: my-app
Modules: 2, Services: 1
Languages: go, typescript
  Module: auth (./auth)
    deps: core
  Module: api (./api)
    deps: auth
  Service: gateway kind=http port=8080
```

## Penggunaan Go

```go
import (
    "github.com/NAEOS-foundation/naeos/internal/context/bundle"
    "github.com/NAEOS-foundation/naeos/internal/specification/parser"
)

// Dari spesifikasi
p := parser.NewParser()
doc, err := p.Parse(specContent)

gen := bundle.NewGenerator(nil)
contextBundle := gen.GenerateFromSpec(doc)

// Output markdown
markdown := contextBundle.ToMarkdown()

// Output plain text
plainText := contextBundle.ToPlainText()

// Supported targets
targets := contextBundle.SupportedTargets()
// ["markdown", "plain"]
```

### Dari NEIR

```go
contextBundle := gen.GenerateFromNEIR(neir)
```

## Struktur Bundle

```go
type Bundle struct {
    Project   string
    Summary   string
    Modules   []ModuleContext
    Services  []ServiceContext
    Languages []string
    Targets   []string
    NEIR      string
    Raw       string
    Metadata  map[string]string
}

type ModuleContext struct {
    Name         string
    Path         string
    Description  string
    Dependencies []string
}

type ServiceContext struct {
    Name      string
    Kind      string
    Port      int
    Endpoints []EndpointContext
}

type EndpointContext struct {
    Method string
    Path   string
    Action string
}
```

## Metadata

Bundle menyertakan metadata:

```json
{
  "generated_by": "naeos-context-bundle",
  "module_count": "2",
  "service_count": "1"
}
```

## Integration dengan AI Tools

```go
// Generate context untuk Copilot
contextBundle := gen.GenerateFromSpec(doc)
copilotInstructions := contextBundle.ToMarkdown()

// Simpan ke .github/copilot-context.md
os.WriteFile(".github/copilot-context.md", []byte(copilotInstructions), 0644)
```
