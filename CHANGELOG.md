# Changelog

## v3.0.0 (2026-07-23)

### NEIR v2.0 Specification

- **Conditional Modules**: Modules can now have `when` conditions evaluated against environment variables
  ```yaml
  modules:
    - name: feature-x
      path: ./feature-x
      when:
        env: FEATURE_X_ENABLED
        equals: "true"
  ```
- **Environment Profiles**: Override module/service configuration per environment
  ```yaml
  profiles:
    - name: production
      modules:
        - name: core
          path: ./core-prod
      services:
        - name: api
          port: 8080
  ```
- **Module Inheritance**: Modules can `extend` a base module, inheriting path, dependencies, and description
  ```yaml
  modules:
    - name: base
      path: ./base
      dependencies: [utils]
    - name: app
      extend: base
  ```

### Enterprise Features

- **Audit Log System**: Structured audit logging with JSON and Splunk HEC output
  - Supports filtering by action, user, and severity
  - Built-in `JSONWriter` for file-based audit logs
  - Built-in `SplunkWriter` for Splunk HTTP Event Collector integration
- **RBAC (Role-Based Access Control)**: Built-in roles (admin, developer, viewer, operator) with customizable permissions
  - `spec:read`, `spec:write`, `spec:delete`, `spec:execute`
  - `user:read`, `user:manage`, `team:read`, `team:manage`
  - `audit:read`, `config:read`, `config:write`
- **Compliance Reports**: Generate SOC2 and HIPAA compliance reports with pass/fail/partial scoring

### LSP Server (from v2.1)

- Real-time diagnostics, autocompletion, and hover information for spec YAML files
- `naeos lsp` command to start the language server

### AI & Developer Experience

- Knowledge graph integration with `naeos ai suggest`
- HTML diff output with `naeos diff --format html`
- VS Code extension scaffold with LSP support

### Performance (from v2.1)

- Parallel multi-module generation (17% faster pipeline, 44% faster validation)
- Per-stage incremental caching
- Benchmark suite for small/medium/large projects

### Migration from v2.x

1. NEIR v2.0 is backward-compatible with v1.x specs
2. New fields (`when`, `profiles`, `extend`) are optional
3. Enterprise features require explicit import from `internal/enterprise/`
4. LSP server runs via `naeos lsp` command
