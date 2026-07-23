package lsp

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestServerInitialize(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	initReq := `{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}`
	if err := s.Handle([]byte(initReq)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(buf.Bytes()[findBodyStart(buf.Bytes()):], &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected result to be map")
	}

	caps, ok := result["capabilities"].(map[string]any)
	if !ok {
		t.Fatal("expected capabilities")
	}

	if caps["hoverProvider"] != true {
		t.Error("expected hoverProvider = true")
	}
	if caps["textDocumentSync"] != float64(1) {
		t.Errorf("textDocumentSync = %v, want 1", caps["textDocumentSync"])
	}
}

func TestServerDidOpen(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	didOpen := `{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///test.yaml","text":"project: test-proj"}}}`
	if err := s.Handle([]byte(didOpen)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if _, ok := s.documents["file:///test.yaml"]; !ok {
		t.Error("document not stored")
	}
}

func TestServerDidChange(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	s.documents["file:///test.yaml"] = "project: old"
	change := `{"jsonrpc":"2.0","method":"textDocument/didChange","params":{"textDocument":{"uri":"file:///test.yaml","version":2},"contentChanges":[{"text":"project: new"}]}}`
	if err := s.Handle([]byte(change)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if s.documents["file:///test.yaml"] != "project: new" {
		t.Errorf("document not updated: %q", s.documents["file:///test.yaml"])
	}
}

func TestServerHover(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	s.documents["file:///test.yaml"] = "project: test\nmodules:\n  - name: core"

	hoverReq := `{"jsonrpc":"2.0","method":"textDocument/hover","params":{"textDocument":{"uri":"file:///test.yaml"},"position":{"line":0,"character":2}},"id":2}`
	if err := s.Handle([]byte(hoverReq)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(buf.Bytes()[findBodyStart(buf.Bytes()):], &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected hover result")
	}

	contents, ok := result["contents"].(map[string]any)
	if !ok {
		t.Fatal("expected contents")
	}

	if !strings.Contains(contents["value"].(string), "project") {
		t.Error("expected hover to mention 'project'")
	}
}

func TestServerCompletion(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	s.documents["file:///test.yaml"] = "project: test\n"

	completeReq := `{"jsonrpc":"2.0","method":"textDocument/completion","params":{"textDocument":{"uri":"file:///test.yaml"},"position":{"line":1,"character":0}},"id":3}`
	if err := s.Handle([]byte(completeReq)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(buf.Bytes()[findBodyStart(buf.Bytes()):], &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected completion result")
	}

	items, ok := result["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatal("expected completion items")
	}

	first := items[0].(map[string]any)
	if first["label"] == nil {
		t.Error("expected label in completion item")
	}
}

func TestServerCompletionKindValues(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	s.documents["file:///test.yaml"] = "project: test\nservices:\n  - name: api\n    kind: "

	completeReq := `{"jsonrpc":"2.0","method":"textDocument/completion","params":{"textDocument":{"uri":"file:///test.yaml"},"position":{"line":3,"character":10}},"id":3}`
	if err := s.Handle([]byte(completeReq)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(buf.Bytes()[findBodyStart(buf.Bytes()):], &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected completion result")
	}

	items, ok := result["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatal("expected kind completion items")
	}

	found := false
	for _, item := range items {
		m := item.(map[string]any)
		if m["label"] == "http" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'http' in kind completions")
	}
}

func TestServerValidateEmpty(t *testing.T) {
	diagnostics := (&Server{}).validate("")
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics for empty input, got %d", len(diagnostics))
	}
}

func TestServerValidateMissingProject(t *testing.T) {
	diagnostics := (&Server{}).validate("modules:\n  - name: core")
	found := false
	for _, d := range diagnostics {
		if d.Severity == SeverityWarning && strings.Contains(d.Message, "project") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning about missing project")
	}
}

func TestServerValidateValid(t *testing.T) {
	diagnostics := (&Server{}).validate("project: test\nmodules:\n  - name: core\n    path: ./core")
	errCount := 0
	for _, d := range diagnostics {
		if d.Severity == SeverityError {
			errCount++
		}
	}
	if errCount != 0 {
		t.Errorf("expected 0 errors for valid spec, got %d", errCount)
	}
}

func TestServerHoverForLine(t *testing.T) {
	s := &Server{}
	tests := []struct {
		line   string
		expect bool
	}{
		{"project:", true},
		{"modules:", true},
		{"services:", true},
		{"kind:", true},
		{"port:", true},
		{"pattern:", true},
		{"strategy:", true},
		{"unknown_key:", false},
	}

	for _, tt := range tests {
		result := s.hoverForLine(tt.line)
		if tt.expect && result == "" {
			t.Errorf("hoverForLine(%q) returned empty, expected content", tt.line)
		}
		if !tt.expect && result != "" {
			t.Errorf("hoverForLine(%q) returned content, expected empty", tt.line)
		}
	}
}

func TestServerShutdown(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	shutdownReq := `{"jsonrpc":"2.0","method":"shutdown","id":4}`
	if err := s.Handle([]byte(shutdownReq)); err != nil {
		t.Fatalf("Handle: %v", err)
	}
}

func TestServerUnknownMethod(t *testing.T) {
	var buf bytes.Buffer
	s := NewServer(&buf)

	req := `{"jsonrpc":"2.0","method":"unknown/method","id":5}`
	if err := s.Handle([]byte(req)); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(buf.Bytes()[findBodyStart(buf.Bytes()):], &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("error code = %d, want -32601", resp.Error.Code)
	}
}

func TestGetLine(t *testing.T) {
	s := &Server{}
	text := "line1\nline2\nline3"
	tests := []struct {
		lineNum int
		expect  string
	}{
		{0, "line1"},
		{1, "line2"},
		{2, "line3"},
		{5, ""},
	}
	for _, tt := range tests {
		got := s.getLine(text, tt.lineNum)
		if got != tt.expect {
			t.Errorf("getLine(%d) = %q, want %q", tt.lineNum, got, tt.expect)
		}
	}
}

func findBodyStart(data []byte) int {
	idx := bytes.Index(data, []byte("\r\n\r\n"))
	if idx == -1 {
		idx = bytes.Index(data, []byte("\n\n"))
		if idx != -1 {
			return idx + 2
		}
		return 0
	}
	return idx + 4
}
