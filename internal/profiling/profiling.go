package profiling

import (
	"fmt"
	"strings"
	"time"
)

type StageMetrics struct {
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration
	Metadata  map[string]any
}

type PipelineProfile struct {
	Stages      []StageMetrics
	TotalStart  time.Time
	TotalEnd    time.Time
	TotalTime   time.Duration
}

func NewProfile() *PipelineProfile {
	return &PipelineProfile{}
}

func (p *PipelineProfile) StartStage(name string) *StageMetrics {
	stage := &StageMetrics{
		Name:      name,
		StartedAt: time.Now(),
		Metadata:  make(map[string]any),
	}
	p.Stages = append(p.Stages, *stage)
	return &p.Stages[len(p.Stages)-1]
}

func (p *PipelineProfile) EndStage(name string) {
	for i := range p.Stages {
		if p.Stages[i].Name == name && p.Stages[i].EndedAt.IsZero() {
			p.Stages[i].EndedAt = time.Now()
			p.Stages[i].Duration = p.Stages[i].EndedAt.Sub(p.Stages[i].StartedAt)
			return
		}
	}
}

func (p *PipelineProfile) Finish() {
	p.TotalEnd = time.Now()
	p.TotalTime = p.TotalEnd.Sub(p.TotalStart)
}

func (p *PipelineProfile) Start() {
	p.TotalStart = time.Now()
}

func (p *PipelineProfile) Summary() string {
	var sb strings.Builder
	sb.WriteString("Pipeline Performance Profile\n")
	sb.WriteString("============================\n\n")
	sb.WriteString(fmt.Sprintf("Total time: %s\n\n", p.TotalTime.Round(time.Microsecond)))
	sb.WriteString("Stage Breakdown:\n")
	sb.WriteString(fmt.Sprintf("  %-20s %15s %8s\n", "Stage", "Duration", "%"))
	sb.WriteString(fmt.Sprintf("  %-20s %15s %8s\n", "-----", "--------", "--"))
	for _, stage := range p.Stages {
		pct := 0.0
		if p.TotalTime > 0 {
			pct = float64(stage.Duration) / float64(p.TotalTime) * 100
		}
		sb.WriteString(fmt.Sprintf("  %-20s %15s %7.1f%%\n",
			stage.Name,
			stage.Duration.Round(time.Microsecond),
			pct))
	}
	return sb.String()
}

func (p *PipelineProfile) SlowestStage() *StageMetrics {
	if len(p.Stages) == 0 {
		return nil
	}
	slowest := &p.Stages[0]
	for i := range p.Stages {
		if p.Stages[i].Duration > slowest.Duration {
			slowest = &p.Stages[i]
		}
	}
	return slowest
}

func (p *PipelineProfile) FastestStage() *StageMetrics {
	if len(p.Stages) == 0 {
		return nil
	}
	fastest := &p.Stages[0]
	for i := range p.Stages {
		if p.Stages[i].Duration < fastest.Duration && p.Stages[i].Duration > 0 {
			fastest = &p.Stages[i]
		}
	}
	return fastest
}

type BenchmarkResult struct {
	Name       string
	Iterations int
	AvgTime    time.Duration
	MinTime    time.Duration
	MaxTime    time.Duration
	TotalTime  time.Duration
}

func Benchmark(name string, iterations int, fn func()) *BenchmarkResult {
	result := &BenchmarkResult{
		Name:       name,
		Iterations: iterations,
		MinTime:    time.Hour,
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		fn()
		duration := time.Since(start)
		result.TotalTime += duration
		if duration < result.MinTime {
			result.MinTime = duration
		}
		if duration > result.MaxTime {
			result.MaxTime = duration
		}
	}

	result.AvgTime = result.TotalTime / time.Duration(iterations)
	return result
}

func (r *BenchmarkResult) String() string {
	return fmt.Sprintf("Benchmark %s: %d iterations, avg=%s, min=%s, max=%s, total=%s",
		r.Name, r.Iterations,
		r.AvgTime.Round(time.Microsecond),
		r.MinTime.Round(time.Microsecond),
		r.MaxTime.Round(time.Microsecond),
		r.TotalTime.Round(time.Microsecond))
}
