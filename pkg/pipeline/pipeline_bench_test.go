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
