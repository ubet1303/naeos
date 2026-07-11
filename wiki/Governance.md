# Governance

Model tata kelola NAEOS untuk mengatur kebijakan, validasi, dan review.

## Struktur Governance

```
┌─────────────────────────────────────────┐
│            Governance Model             │
├─────────────┬─────────────┬─────────────┤
│  Charter    │  Policies   │  Reviews    │
│  Vision     │  Rules      │  Audit      │
│  Mission    │  Evaluator  │  Trail      │
│  Manifesto  │  Enforcer   │  Reports    │
└─────────────┴─────────────┴─────────────┘
```

## Document Hierarchy

1. **NAEOS-CON-001** — Engineering Constitution
2. **NAEOS-GOV-001** — Project Charter
3. **NAEOS-GOV-002** — Vision
4. **NAEOS-GOV-003** — Mission
5. **NAEOS-GOV-004** — Manifesto
6. **NAEOS-GOV-005** — Core Principles
7. **NAEOS-GOV-006** — Governance Model
8. **NAEOS-GOV-007** — Roadmap
9. **NAEOS-GOV-008** — Versioning Policy

## Policy System

### Default Rules

```go
policy.DefaultRules() []Rule
```

| Rule ID | Condition | Priority | Action |
|---------|-----------|----------|--------|
| project-required | exists:project | 1 | block |
| modules-required | exists:modules | 1 | block |
| architecture-pattern-valid | in:architecture_pattern,... | 2 | warn |
| deployment-strategy-valid | in:deployment_strategy,... | 2 | warn |
| service-port-positive | gt:service_port,0 | 3 | warn |

### Condition Operators

| Operator | Contoh | Deskripsi |
|----------|--------|-----------|
| `exists` | `exists:project` | Cek apakah key ada |
| `not_empty` | `not_empty:project` | Cek apakah key tidak kosong |
| `contains` | `contains:project,app` | Cek apakah value mengandung substring |
| `gt` | `gt:port,0` | Greater than |
| `lt` | `lt:port,65535` | Less than |
| `in` | `in:kind,http,grpc` | Value ada di list |

### Custom Rules

```go
rules := []policy.Rule{
    {
        RuleID:    "require-testing",
        Condition: "exists:testing",
        Priority:  1,
        Action:    "block",
        Scope:     "spec",
        Enabled:   true,
    },
}
```

## Review System

```go
reviewer := review.NewReviewer()

result, err := reviewer.ReviewArtifact(
    "main.go",
    string(content),
    []string{"no-todo", "no-placeholder"},
)
```

## Audit

```bash
naeos audit --input-file spec.yaml
```

## Enforcement

### CI/CD Integration

```yaml
# .github/workflows/ci.yml
- name: Validate Spec
  run: naeos validate --input-file spec.yaml

- name: Lint Spec
  run: naeos lint --input-file spec.yaml

- name: Audit
  run: naeos audit --input-file spec.yaml
```

### Pipeline Integration

```go
cfg := pipeline.Config{
    Policies: policy.DefaultRules(),
}
```
