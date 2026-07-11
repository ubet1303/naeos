package pipeline

import (
	"testing"
)

func BenchmarkPipelineRunMinimal(b *testing.B) {
	spec := `
project: bench-api
modules:
  - name: core
    path: ./internal/core
services:
  - name: api
    kind: http
    port: 8080
`
	p, err := New(Config{
		Name: "bench",
		Mode: "development",
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Run(spec)
	}
}

func BenchmarkPipelineRunMultiLanguage(b *testing.B) {
	spec := `
project: bench-multi
modules:
  - name: auth
    path: ./internal/auth
  - name: user
    path: ./internal/user
services:
  - name: gateway
    kind: http
    port: 8080
  - name: worker
    kind: worker
    port: 9090
generation:
  languages:
    - go
    - typescript
    - python
`
	p, err := New(Config{
		Name: "bench-multi",
		Mode: "development",
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Run(spec)
	}
}

func BenchmarkPipelineValidate(b *testing.B) {
	spec := `
project: bench-validate
modules:
  - name: core
    path: ./internal/core
`
	p, err := New(Config{
		Name: "bench-validate",
		Mode: "development",
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Validate(spec)
	}
}

func BenchmarkPipelineNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = New(Config{
			Name: "bench-new",
			Mode: "development",
		})
	}
}

func BenchmarkPipelineRunDryRun(b *testing.B) {
	spec := `
project: bench-dryrun
modules:
  - name: core
    path: ./internal/core
services:
  - name: api
    kind: http
    port: 8080
`
	p, err := New(Config{
		Name:    "bench-dryrun",
		Mode:    "development",
		DryRun:  true,
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Run(spec)
	}
}

func BenchmarkPipelineRunFullSpec(b *testing.B) {
	spec := `
project: bench-full
modules:
  - name: auth
    path: ./internal/auth
    dependencies:
      - core
  - name: core
    path: ./internal/core
  - name: user
    path: ./internal/user
    dependencies:
      - core
  - name: order
    path: ./internal/order
    dependencies:
      - core
      - user
services:
  - name: api-gateway
    kind: http
    port: 8080
  - name: worker
    kind: worker
    port: 9090
architecture:
  pattern: hexagonal
deployment:
  strategy: blue-green
testing:
  strategy: unit
generation:
  languages:
    - go
    - typescript
`
	p, err := New(Config{
		Name: "bench-full",
		Mode: "development",
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Run(spec)
	}
}

func BenchmarkPipelineRunVerbose(b *testing.B) {
	spec := `
project: bench-verbose
modules:
  - name: core
    path: ./internal/core
`
	p, err := New(Config{
		Name:    "bench-verbose",
		Mode:    "development",
		Verbose: true,
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Run(spec)
	}
}
