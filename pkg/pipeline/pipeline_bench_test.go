package pipeline

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

const benchSpec = `project:
  name: benchapp
  version: "1.0.0"
services:
  - name: api
    port: 8080
  - name: worker
    port: 9090
`

func BenchmarkPipelineRun(b *testing.B) {
	cfg := Config{
		Name:      "benchapp",
		OutputDir: b.TempDir(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, _ := New(cfg)
		p.Run(strings.TrimSpace(benchSpec))
	}
}

func BenchmarkPipelineValidate(b *testing.B) {
	cfg := Config{
		Name:      "benchapp",
		OutputDir: b.TempDir(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, _ := New(cfg)
		p.Validate(strings.TrimSpace(benchSpec))
	}
}

func BenchmarkPipelineNew(b *testing.B) {
	cfg := Config{
		Name:      "benchapp",
		OutputDir: b.TempDir(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(cfg)
	}
}

func buildSpec(modules, services int) string {
	var sb strings.Builder
	sb.WriteString("project: bench-scale\n")
	if modules > 0 {
		sb.WriteString("modules:\n")
		for i := 0; i < modules; i++ {
			fmt.Fprintf(&sb, "  - name: mod%d\n    path: ./internal/mod%d\n", i, i)
		}
	}
	if services > 0 {
		sb.WriteString("services:\n")
		for i := 0; i < services; i++ {
			fmt.Fprintf(&sb, "  - name: svc%d\n    kind: http\n    port: %d\n", i, 8080+i)
		}
	}
	return sb.String()
}

func BenchmarkPipelineRunSmall(b *testing.B) {
	spec := buildSpec(5, 2)
	cfg := Config{Name: "bench-small", OutputDir: b.TempDir()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, _ := New(cfg)
		p.Run(spec)
	}
}

func BenchmarkPipelineRunMedium(b *testing.B) {
	spec := buildSpec(50, 10)
	cfg := Config{Name: "bench-medium", OutputDir: b.TempDir()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, _ := New(cfg)
		p.Run(spec)
	}
}

func BenchmarkPipelineRunLarge(b *testing.B) {
	spec := buildSpec(500, 50)
	cfg := Config{Name: "bench-large", OutputDir: b.TempDir()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, _ := New(cfg)
		p.Run(spec)
	}
}

func BenchmarkPipelineMemoryGrowth(b *testing.B) {
	spec := buildSpec(50, 10)
	cfg := Config{Name: "bench-mem", OutputDir: b.TempDir()}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)

		p, _ := New(cfg)
		p.Run(spec)

		runtime.GC()
		runtime.ReadMemStats(&after)

		b.ReportMetric(float64(after.TotalAlloc-before.TotalAlloc), "bytes_total")
	}
}
