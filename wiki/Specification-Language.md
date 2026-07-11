# Specification Language v2

Bahasa spesifikasi NAEOS mendukung fitur lanjutan untuk interoperabilitas dan fleksibilitas.

## Dasar

```yaml
project: my-app
version: 0.3.0
modules:
  - name: auth
    path: ./auth
    description: Authentication module
    dependencies: []
services:
  - name: gateway
    kind: http
    port: 8080
architecture:
  pattern: hexagonal
  principles: [loose-coupling, high-cohesion]
deployment:
  strategy: rolling
  environments: [staging, production]
testing:
  strategy: unit
  coverage: "80"
generation:
  languages: [go, typescript]
```

## Fitur v2

### Variable Interpolation

Gunakan `${var}` untuk variabel kustom:

```yaml
project: ${project_name}
modules:
  - name: ${module_name}
    path: ./${module_name}
```

**Resolver:**
```go
resolver := parser.NewVariableResolver()
resolver.SetVar("project_name", "my-app")
resolver.SetVar("module_name", "auth")

result, err := resolver.Resolve(input)
```

### Environment Variables

Gunakan `$env{VAR}` untuk membaca environment variables:

```yaml
project: my-app
services:
  - name: api
    kind: http
    port: $env{API_PORT}
```

**Resolver:**
```go
resolver := parser.NewVariableResolver()
result, err := resolver.Resolve("$env{API_PORT}")
// Reads from os.Getenv("API_PORT")
```

### Reference Resolution

Gunakan `$ref{path}` untuk cross-reference:

```yaml
project: my-app
modules:
  - name: auth
    path: ./auth
services:
  - name: api
    dependencies:
      - $ref{modules.auth.name}
```

**Resolver:**
```go
resolver := parser.NewVariableResolver()
resolver.SetRef("modules.auth.name", "auth")

result, err := resolver.Resolve("$ref{modules.auth.name}")
// Output: "auth"
```

### Recursive Resolution

Resolver bekerja secara recursive pada maps, slices, dan nested structures:

```go
input := map[string]any{
    "project": "${name}",
    "services": []any{
        map[string]any{
            "name": "${name}-service",
            "port": 8080,
        },
    },
}

result, err := resolver.ResolveMap(input)
// project: "my-app"
// services[0].name: "my-app-service"
// services[0].port: 8080
```

## Schema Versioning

### Auto-check

Parser otomatis memeriksa field `version`:

```yaml
version: 0.3.0  # Auto-checked on parse
```

### Version Check

```go
result := parser.CheckSpecVersion("0.3.0")
// result.Valid = true
// result.Current = SchemaVersion{0, 3, 0}
// result.Required = SchemaVersion{0, 1, 0}
```

### Constants

```go
const (
    MinSpecVersion    = "0.1.0"
    CurrentSpecVersion = "0.3.0"
)
```

### Version Comparison

```go
v1, _ := parser.ParseSchemaVersion("0.1.0")
v2, _ := parser.ParseSchemaVersion("0.3.0")

v2.GreaterThan(v1)   // true
v1.LessThan(v2)      // true
v1.CompatibleWith(v1) // true
v2.CompatibleWith(v1) // true
```

## Nested Structures

```yaml
project: my-app
architecture:
  pattern: hexagonal
  principles:
    - loose-coupling
    - high-cohesion
  details:
    layers:
      - name: domain
        description: Business logic
      - name: infrastructure
        description: External integrations
```

## Multi-document

```yaml
---
project: core
modules:
  - name: foundation
    path: ./foundation
---
project: extensions
modules:
  - name: plugin-system
    path: ./plugin-system
    dependencies: [foundation]
```
