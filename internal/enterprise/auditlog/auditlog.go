package auditlog

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Action string

const (
	ActionCreate  Action = "create"
	ActionRead    Action = "read"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionExecute Action = "execute"
	ActionExport  Action = "export"
	ActionAuth    Action = "auth"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

type Entry struct {
	Timestamp  time.Time         `json:"timestamp"`
	UserID     string            `json:"user_id"`
	TenantID   string            `json:"tenant_id,omitempty"`
	Action     Action            `json:"action"`
	Resource   string            `json:"resource"`
	ResourceID string            `json:"resource_id,omitempty"`
	Severity   Severity          `json:"severity"`
	Details    string            `json:"details,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	IPAddress  string            `json:"ip_address,omitempty"`
	UserAgent  string            `json:"user_agent,omitempty"`
	Success    bool              `json:"success"`
	ErrorCode  string            `json:"error_code,omitempty"`
}

type Writer interface {
	Write(entry Entry) error
	Close() error
}

type Logger struct {
	writers []Writer
	mu      sync.RWMutex
}

func NewLogger(writers ...Writer) *Logger {
	return &Logger{writers: writers}
}

func (l *Logger) AddWriter(w Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, w)
}

func (l *Logger) Log(entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if entry.Severity == "" {
		entry.Severity = SeverityInfo
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	var firstErr error
	for _, w := range l.writers {
		if err := w.Write(entry); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	var firstErr error
	for _, w := range l.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// JSONWriter writes audit entries as JSON lines to a file.
type JSONWriter struct {
	file *os.File
	enc  *json.Encoder
	mu   sync.Mutex
}

func NewJSONWriter(path string) (*JSONWriter, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open audit log: %w", err)
	}
	return &JSONWriter{
		file: f,
		enc:  json.NewEncoder(f),
	}, nil
}

func (w *JSONWriter) Write(entry Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.enc.Encode(entry)
}

func (w *JSONWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}

// MemoryWriter stores entries in memory (for testing).
type MemoryWriter struct {
	entries []Entry
	mu      sync.Mutex
}

func NewMemoryWriter() *MemoryWriter {
	return &MemoryWriter{}
}

func (w *MemoryWriter) Write(entry Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = append(w.entries, entry)
	return nil
}

func (w *MemoryWriter) Close() error { return nil }

func (w *MemoryWriter) Entries() []Entry {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]Entry, len(w.entries))
	copy(result, w.entries)
	return result
}

// SplunkWriter writes audit entries in HEC (HTTP Event Collector) format.
type SplunkWriter struct {
	entries []map[string]any
	mu      sync.Mutex
}

func NewSplunkWriter() *SplunkWriter {
	return &SplunkWriter{}
}

func (w *SplunkWriter) Write(entry Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	hevent := map[string]any{
		"time":       entry.Timestamp.Unix(),
		"host":       entry.Resource,
		"source":     "naeos-audit",
		"sourcetype": "naeos:audit",
		"event": map[string]any{
			"user_id":     entry.UserID,
			"tenant_id":   entry.TenantID,
			"action":      string(entry.Action),
			"resource":    entry.Resource,
			"resource_id": entry.ResourceID,
			"severity":    string(entry.Severity),
			"details":     entry.Details,
			"success":     entry.Success,
			"error_code":  entry.ErrorCode,
			"ip_address":  entry.IPAddress,
			"user_agent":  entry.UserAgent,
		},
	}
	w.entries = append(w.entries, hevent)
	return nil
}

func (w *SplunkWriter) Close() error { return nil }

func (w *SplunkWriter) Events() []map[string]any {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]map[string]any, len(w.entries))
	copy(result, w.entries)
	return result
}

// FilterByAction returns entries matching the given action.
func FilterByAction(entries []Entry, action Action) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.Action == action {
			result = append(result, e)
		}
	}
	return result
}

// FilterByUser returns entries for the given user ID.
func FilterByUser(entries []Entry, userID string) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result
}

// FilterBySeverity returns entries at or above the given severity level.
func FilterBySeverity(entries []Entry, min Severity) []Entry {
	levels := map[Severity]int{
		SeverityInfo: 0, SeverityWarning: 1, SeverityError: 2, SeverityCritical: 3,
	}
	minLevel := levels[min]
	var result []Entry
	for _, e := range entries {
		if levels[e.Severity] >= minLevel {
			result = append(result, e)
		}
	}
	return result
}
