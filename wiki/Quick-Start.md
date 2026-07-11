# Quick Start

## Prasyarat

- Go 1.22+
- Git

## Instalasi

```bash
git clone https://github.com/NAEOS-foundation/naeos.git
cd naeos
go build ./cmd/naeos/
```

## Mulai

### 1. Buat spesifikasi

Buat file `spec.yaml`:

```yaml
project: my-app
modules:
  - name: auth
    path: ./auth
    description: Authentication module
  - name: api
    path: ./api
    description: REST API
    dependencies: [auth]
services:
  - name: gateway
    kind: http
    port: 8080
    endpoints:
      - method: GET
        path: /health
        action: healthCheck
architecture:
  pattern: hexagonal
deployment:
  strategy: rolling
  environments: [staging, production]
testing:
  strategy: unit
  coverage: "80"
generation:
  languages: [go, typescript]
```

### 2. Jalankan pipeline

```bash
naeos run --config config.yaml --input-file spec.yaml
```

### 3. Generate AI context bundle

```bash
naeos context --input-file spec.yaml
```

### 4. Compile ke AI tool

```bash
naeos compile --target copilot --input-file spec.yaml
naeos compile --all --input-file spec.yaml
```

### 5. Cek kesehatan

```bash
naeos doctor
```

## Contoh lengkap

```bash
# Buat project baru
naeos init

# Jalankan pipeline dengan verbose
naeos run --config config.yaml --input-file spec.yaml --verbose

# Preview tanpa tulis ke disk
naeos run --config config.yaml --input-file spec.yaml --dry-run

# Profil industri
naeos profile list
naeos profile apply saas

# Artifact management
naeos artifacts list
naeos artifacts summary

# Schema migration
naeos migrate plan
naeos migrate run
```
