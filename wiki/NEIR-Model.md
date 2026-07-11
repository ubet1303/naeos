# NEIR Model

**Nusantara Engineering Intermediate Representation** — model representasi intermediasi yang merepresentasikan proyek secara unified.

## Struktur NEIR

```go
type NEIR struct {
    Project        *project.Project
    Architecture   *architecture.Architecture
    Domain         *domain.Domain
    Modules        []module.Module
    Components     []component.Component
    Services       []service.Service
    APIs           []api.API
    Storage        []storage.Storage
    Infrastructure *infrastructure.Infrastructure
    Security       *security.Security
    AI             *ai.AI
    Documentation  *docs.Documentation
    Deployment     *deployment.Deployment
    Testing        *testing.Testing
    Metadata       *metadata.Metadata
    Generation     *generation.GenerationConfig
}
```

## Komponen NEIR

### Project
```go
type Project struct {
    Name        string
    Version     string
    Description string
    License     string
    Authors     []string
    Repository  string
    Tags        []string
    Attributes  map[string]string
}
```

### Module
```go
type Module struct {
    Name         string
    Path         string
    Description  string
    Packages     []string
    Dependencies []string
    Attributes   map[string]string
}
```

### Service
```go
type Service struct {
    Name        string
    Kind        ServiceKind  // http, grpc, worker, cli, job
    Port        int
    Description string
    Endpoints   []Endpoint
    Middleware   []string
    Attributes  map[string]string
}
```

### Endpoint
```go
type Endpoint struct {
    Method string  // GET, POST, PUT, DELETE, PATCH
    Path   string
    Action string
}
```

### GenerationConfig
```go
type GenerationConfig struct {
    Languages []language.Language  // go, typescript, python, java, rust
    OutputDir string
    ModuleDir string
}
```

## NEIR Builder

```go
builder := builder.NewBuilder()
neir, err := builder.Build(resolvedSpec)
```

## NEIR Validator

```go
validator := validator.NewValidator()
err := validator.Validate(neir)
```

## NEIR to AI Instruction Sets

```go
compiler := compiler.New()
compiler.Register(copilotAdapter)
compiler.Register(claudeAdapter)

output, err := compiler.Compile(neir, compiler.TargetCopilot)
```

## NEIR Version

```go
import "github.com/NAEOS-foundation/naeos/internal/neir/version"

currentVersion := version.CurrentVersion()
```
