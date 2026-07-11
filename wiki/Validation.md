# Validation

NAEOS menyediakan validasi komprehensif pada multiple level.

## Spec Validation

### Basic Validation

```go
import "github.com/NAEOS-foundation/naeos/internal/lint"

result, err := lint.ValidateSpec(specContent)
if err != nil {
    log.Fatal(err)
}

if !result.Valid {
    for _, issue := range result.Issues {
        fmt.Printf("[%s] %s\n", issue.Severity, issue.Message)
    }
}
```

### CLI

```bash
naeos validate --input-file spec.yaml
naeos lint --input-file spec.yaml
```

## Module Validation

### Duplicate Detection

```go
validator := parser.NewSpecValidator()
issues := validator.ValidateModules(modules)

// issues akan berisi "module-duplicate" jika ada module yang sama
```

### Circular Dependency Detection

```go
modules := []parser.Module{
    {Name: "a", Dependencies: []string{"b"}},
    {Name: "b", Dependencies: []string{"a"}},
}

issues := validator.ValidateModules(modules)
// issues akan berisi "circular-dependency"
```

### Dangling Dependency Detection

```go
modules := []parser.Module{
    {Name: "a", Dependencies: []string{"nonexistent"}},
}

issues := validator.ValidateModules(modules)
// issues akan berisi "dependency-not-found"
```

## Service Validation

### Port Conflict Detection

```go
services := []parser.Service{
    {Name: "api", Port: 8080},
    {Name: "web", Port: 8080},
}

issues := validator.ValidateServices(services)
// issues akan berisi "service-port-conflict"
```

### Port Range Validation

```go
services := []parser.Service{
    {Name: "bad", Port: 99999},
}

issues := validator.ValidateServices(services)
// issues akan berisi "service-port-range"
```

## Schema Version Validation

### Auto-check on Parse

Parser otomatis memeriksa versi saat parsing:

```yaml
version: 0.3.0  # Auto-checked
```

### Manual Check

```go
result := parser.CheckSpecVersion("0.3.0")
if !result.Valid {
    fmt.Println(result.Message)
}
```

## Policy Validation

### Default Rules

```go
import "github.com/NAEOS-foundation/naeos/internal/governance/policy"

evaluator := policy.NewEvaluator()
rules := policy.DefaultRules()

ctx := map[string]any{
    "project":  "my-app",
    "modules":  3,
    "services": 2,
}

results, err := evaluator.EvaluateRules(rules, ctx)
```

### Custom Rules

```go
rules := []policy.Rule{
    {
        RuleID:    "require-architecture",
        Condition: "exists:architecture",
        Priority:  1,
        Action:    "block",
        Scope:     "spec",
        Enabled:   true,
    },
}
```

## Validation Result

```go
type ValidationResult struct {
    Valid        bool
    Issues       []ValidationIssue
    Warnings     []ValidationIssue
    ModuleCount  int
    ServiceCount int
}

type ValidationIssue struct {
    Severity string  // "error", "warning", "info"
    Rule     string
    Message  string
    Line     int
}
```
