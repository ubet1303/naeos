package profiling

import (
	"testing"
	"time"
)

func TestNewProfile(t *testing.T) {
	p := NewProfile()
	if p == nil {
		t.Fatal("expected non-nil profile")
	}
}

func TestProfileStartFinish(t *testing.T) {
	p := NewProfile()
	p.Start()
	time.Sleep(time.Millisecond)
	p.Finish()
	if p.TotalTime <= 0 {
		t.Error("expected positive total time")
	}
}

func TestProfileStages(t *testing.T) {
	p := NewProfile()
	p.Start()
	p.StartStage("parse")
	time.Sleep(time.Millisecond)
	p.EndStage("parse")
	p.StartStage("generate")
	time.Sleep(time.Millisecond)
	p.EndStage("generate")
	p.Finish()

	if len(p.Stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(p.Stages))
	}
	if p.Stages[0].Duration <= 0 {
		t.Error("expected positive duration for parse stage")
	}
}

func TestSlowestStage(t *testing.T) {
	p := NewProfile()
	p.StartStage("fast")
	time.Sleep(time.Millisecond)
	p.EndStage("fast")
	p.StartStage("slow")
	time.Sleep(10 * time.Millisecond)
	p.EndStage("slow")

	slowest := p.SlowestStage()
	if slowest.Name != "slow" {
		t.Errorf("expected slow, got %s", slowest.Name)
	}
}

func TestFastestStage(t *testing.T) {
	p := NewProfile()
	p.StartStage("slow")
	time.Sleep(10 * time.Millisecond)
	p.EndStage("slow")
	p.StartStage("fast")
	time.Sleep(time.Millisecond)
	p.EndStage("fast")

	fastest := p.FastestStage()
	if fastest.Name != "fast" {
		t.Errorf("expected fast, got %s", fastest.Name)
	}
}

func TestSummary(t *testing.T) {
	p := NewProfile()
	p.Start()
	p.StartStage("parse")
	time.Sleep(time.Millisecond)
	p.EndStage("parse")
	p.Finish()

	summary := p.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestBenchmark(t *testing.T) {
	result := Benchmark("test", 10, func() {
		time.Sleep(time.Millisecond)
	})
	if result.Iterations != 10 {
		t.Errorf("expected 10 iterations, got %d", result.Iterations)
	}
	if result.AvgTime <= 0 {
		t.Error("expected positive avg time")
	}
	if result.String() == "" {
		t.Error("expected non-empty string")
	}
}

func TestSlowestStageEmpty(t *testing.T) {
	p := NewProfile()
	if p.SlowestStage() != nil {
		t.Error("expected nil for empty stages")
	}
}

func TestFastestStageEmpty(t *testing.T) {
	p := NewProfile()
	if p.FastestStage() != nil {
		t.Error("expected nil for empty stages")
	}
}
