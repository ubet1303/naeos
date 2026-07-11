package diff

import (
	"testing"
)

func TestComputeDiffAdded(t *testing.T) {
	d := ComputeDiff("", "hello world", "new.txt")
	if d.Type != ChangeAdded {
		t.Errorf("expected ChangeAdded, got %s", d.Type)
	}
	if d.NewSize != 11 {
		t.Errorf("expected size 11, got %d", d.NewSize)
	}
}

func TestComputeDiffRemoved(t *testing.T) {
	d := ComputeDiff("hello world", "", "deleted.txt")
	if d.Type != ChangeRemoved {
		t.Errorf("expected ChangeRemoved, got %s", d.Type)
	}
	if d.OldSize != 11 {
		t.Errorf("expected size 11, got %d", d.OldSize)
	}
}

func TestComputeDiffUnchanged(t *testing.T) {
	d := ComputeDiff("same content", "same content", "same.txt")
	if d.Type != ChangeUnchanged {
		t.Errorf("expected ChangeUnchanged, got %s", d.Type)
	}
}

func TestComputeDiffModified(t *testing.T) {
	d := ComputeDiff("old content", "new content", "modified.txt")
	if d.Type != ChangeModified {
		t.Errorf("expected ChangeModified, got %s", d.Type)
	}
	if d.OldSize != 11 || d.NewSize != 11 {
		t.Errorf("expected sizes 11/11, got %d/%d", d.OldSize, d.NewSize)
	}
}

func TestSummary(t *testing.T) {
	diffs := []*FileDiff{
		{Type: ChangeAdded},
		{Type: ChangeRemoved},
		{Type: ChangeModified},
		{Type: ChangeUnchanged},
		{Type: ChangeModified},
	}
	added, removed, modified, unchanged := Summary(diffs)
	if added != 1 || removed != 1 || modified != 2 || unchanged != 1 {
		t.Errorf("expected 1,1,2,1 got %d,%d,%d,%d", added, removed, modified, unchanged)
	}
}

func TestFormatDiffAdded(t *testing.T) {
	d := &FileDiff{Path: "new.txt", Type: ChangeAdded}
	output := FormatDiff(d)
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatDiffUnchanged(t *testing.T) {
	d := &FileDiff{Path: "same.txt", Type: ChangeUnchanged}
	output := FormatDiff(d)
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatDiffNil(t *testing.T) {
	output := FormatDiff(nil)
	if output != "" {
		t.Error("expected empty output for nil")
	}
}

func TestComputeLineDiff(t *testing.T) {
	old := []string{"line1", "line2", "line3"}
	new := []string{"line1", "modified", "line3", "line4"}
	lines := computeLineDiff(old, new)
	if len(lines) == 0 {
		t.Error("expected non-empty diff lines")
	}
}

func TestAddedLines(t *testing.T) {
	lines := addedLines("a\nb\nc")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if l.Type != ChangeAdded {
			t.Errorf("expected ChangeAdded, got %s", l.Type)
		}
	}
}

func TestRemovedLines(t *testing.T) {
	lines := removedLines("a\nb")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if l.Type != ChangeRemoved {
			t.Errorf("expected ChangeRemoved, got %s", l.Type)
		}
	}
}
