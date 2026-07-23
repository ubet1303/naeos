package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/NAEOS-foundation/naeos/internal/specification/normalizer"
	"github.com/NAEOS-foundation/naeos/internal/specification/parser"
	"github.com/NAEOS-foundation/naeos/internal/specification/resolver"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      any             `json:"id"`
}

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	ID      any    `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type InitializeParams struct {
	ProcessID    int        `json:"processId"`
	RootURI      string     `json:"rootUri"`
	Capabilities ClientCaps `json:"capabilities"`
}

type ClientCaps struct {
	TextDocument *TextDocCaps `json:"textDocument,omitempty"`
}

type TextDocCaps struct {
	Completion *CompletionCaps `json:"completion,omitempty"`
	Hover      *HoverCaps      `json:"hover,omitempty"`
}

type CompletionCaps struct {
	CompletionItem *CompletionItemCaps `json:"completionItem,omitempty"`
}

type CompletionItemCaps struct {
	SnippetSupport bool `json:"snippetSupport"`
}

type HoverCaps struct {
	ContentFormat []string `json:"contentFormat"`
}

type InitializeResult struct {
	Capabilities ServerCaps `json:"capabilities"`
}

type ServerCaps struct {
	TextDocumentSync   int             `json:"textDocumentSync"`
	CompletionProvider *CompletionProv `json:"completionProvider,omitempty"`
	HoverProvider      bool            `json:"hoverProvider"`
	DiagnosticProvider *DiagnosticProv `json:"diagnosticProvider,omitempty"`
}

type CompletionProv struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

type DiagnosticProv struct {
	InterFileDependencies bool `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool `json:"workspaceDiagnostics"`
}

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type DidOpenParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}

type DidChangeParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     int          `json:"version"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type CompletionParams struct {
	TextDocument struct {
		URI string `json:"uri"`
	} `json:"textDocument"`
	Position Position `json:"position"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind,omitempty"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
}

type HoverParams struct {
	TextDocument struct {
		URI string `json:"uri"`
	} `json:"textDocument"`
	Position Position `json:"position"`
}

type HoverResult struct {
	Contents MarkedContent `json:"contents"`
	Range    Range         `json:"range"`
}

type MarkedContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

const (
	SeverityError       = 1
	SeverityWarning     = 2
	SeverityInformation = 3
	SeverityHint        = 4

	CompletionKindKeyword  = 14
	CompletionKindProperty = 10
	CompletionKindValue    = 12
	CompletionKindModule   = 9
	CompletionKindService  = 12

	methodInitialize           = "initialize"
	methodInitialized          = "initialized"
	methodShutdown             = "shutdown"
	methodExit                 = "exit"
	methodTextDocumentDidOpen  = "textDocument/didOpen"
	methodTextDocumentChange   = "textDocument/didChange"
	methodTextDocumentHover    = "textDocument/hover"
	methodTextDocumentComplete = "textDocument/completion"
)

type Server struct {
	documents map[string]string
	mu        sync.RWMutex
	writer    io.Writer
}

func NewServer(w io.Writer) *Server {
	return &Server{
		documents: make(map[string]string),
		writer:    w,
	}
}

func (s *Server) Handle(raw []byte) error {
	var msg struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params,omitempty"`
		ID      json.RawMessage `json:"id,omitempty"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	if msg.ID != nil {
		var id any
		if err := json.Unmarshal(msg.ID, &id); err == nil {
			if id == nil {
				return s.handleNotification(msg.Method, msg.Params)
			}
		}
		return s.handleRequest(msg.Method, msg.Params, id)
	}
	return s.handleNotification(msg.Method, msg.Params)
}

func (s *Server) handleRequest(method string, params json.RawMessage, id any) error {
	var result any
	var rpcErr *Error

	switch method {
	case methodInitialize:
		result = s.handleInitialize(params)
	case methodShutdown:
		result = nil
	case methodTextDocumentHover:
		result = s.handleHover(params)
	case methodTextDocumentComplete:
		result = s.handleCompletion(params)
	default:
		rpcErr = &Error{Code: -32601, Message: fmt.Sprintf("method not found: %s", method)}
	}

	resp := Response{JSONRPC: "2.0", ID: id, Result: result, Error: rpcErr}
	return s.sendResponse(resp)
}

func (s *Server) handleNotification(method string, params json.RawMessage) error {
	switch method {
	case methodInitialized, "":
		return nil
	case methodTextDocumentDidOpen:
		return s.handleDidOpen(params)
	case methodTextDocumentChange:
		return s.handleDidChange(params)
	case methodExit:
		return nil
	}
	return nil
}

func (s *Server) handleInitialize(_ json.RawMessage) InitializeResult {
	return InitializeResult{
		Capabilities: ServerCaps{
			TextDocumentSync: 1,
			CompletionProvider: &CompletionProv{
				TriggerCharacters: []string{":", " ", "-"},
			},
			HoverProvider: true,
			DiagnosticProvider: &DiagnosticProv{
				InterFileDependencies: false,
				WorkspaceDiagnostics:  false,
			},
		},
	}
}

func (s *Server) handleDidOpen(params json.RawMessage) error {
	var p DidOpenParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}
	s.mu.Lock()
	s.documents[p.TextDocument.URI] = p.TextDocument.Text
	s.mu.Unlock()
	return s.validateAndPublish(p.TextDocument.URI, p.TextDocument.Version)
}

func (s *Server) handleDidChange(params json.RawMessage) error {
	var p DidChangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}
	if len(p.ContentChanges) > 0 {
		s.mu.Lock()
		s.documents[p.TextDocument.URI] = p.ContentChanges[len(p.ContentChanges)-1].Text
		s.mu.Unlock()
	}
	return s.validateAndPublish(p.TextDocument.URI, p.TextDocument.Version)
}

func (s *Server) validateAndPublish(uri string, version int) error {
	s.mu.RLock()
	text := s.documents[uri]
	s.mu.RUnlock()

	diagnostics := s.validate(text)

	params := PublishDiagnosticsParams{
		URI:         uri,
		Version:     version,
		Diagnostics: diagnostics,
	}
	return s.sendNotification("textDocument/publishDiagnostics", params)
}

func (s *Server) validate(text string) []Diagnostic {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	var diagnostics []Diagnostic

	p := parser.NewParser("")
	doc, err := p.Parse(text)
	if err != nil {
		diagnostics = append(diagnostics, Diagnostic{
			Range:    Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 10}},
			Severity: SeverityError,
			Source:   "naeos",
			Message:  fmt.Sprintf("Parse error: %v", err),
		})
		return diagnostics
	}

	if doc != nil && doc.Project == "" {
		diagnostics = append(diagnostics, Diagnostic{
			Range:    Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 10}},
			Severity: SeverityWarning,
			Source:   "naeos",
			Message:  "Missing 'project' field — a project name is recommended",
		})
	}

	norm := normalizer.NewNormalizer()
	normalized, err := norm.Normalize(doc)
	if err != nil {
		diagnostics = append(diagnostics, Diagnostic{
			Range:    Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 10}},
			Severity: SeverityError,
			Source:   "naeos",
			Message:  fmt.Sprintf("Normalization error: %v", err),
		})
		return diagnostics
	}

	res := resolver.NewResolver()
	_, err = res.Resolve(normalized)
	if err != nil {
		diagnostics = append(diagnostics, Diagnostic{
			Range:    Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 10}},
			Severity: SeverityError,
			Source:   "naeos",
			Message:  fmt.Sprintf("Resolution error: %v", err),
		})
	}

	return diagnostics
}

func (s *Server) handleHover(params json.RawMessage) any {
	var p HoverParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil
	}

	s.mu.RLock()
	text := s.documents[p.TextDocument.URI]
	s.mu.RUnlock()

	line := s.getLine(text, p.Position.Line)
	trimmed := strings.TrimSpace(line)

	hoverText := s.hoverForLine(trimmed)
	if hoverText == "" {
		return nil
	}

	return HoverResult{
		Contents: MarkedContent{Kind: "markdown", Value: hoverText},
		Range: Range{
			Start: Position{Line: p.Position.Line, Character: 0},
			End:   Position{Line: p.Position.Line, Character: len(line)},
		},
	}
}

func (s *Server) hoverForLine(line string) string {
	hovers := map[string]string{
		"project":      "**project** — The name of the project.\n\nExample: `project: my-app`",
		"modules":      "**modules** — List of application modules. Each module has a name, path, and optional dependencies.",
		"name":         "**name** — The name of the module or service.",
		"path":         "**path** — Filesystem path for the module (e.g., `./internal/core`).",
		"dependencies": "**dependencies** — List of module names this module depends on.",
		"services":     "**services** — List of application services (HTTP, gRPC, worker, etc.).",
		"kind":         "**kind** — Service type: `http`, `grpc`, `worker`, `event`, or `graphql`.",
		"port":         "**port** — Network port for the service (1-65535).",
		"endpoints":    "**endpoints** — List of API endpoints for this service.",
		"method":       "**method** — HTTP method: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`.",
		"action":       "**action** — Action name for this endpoint.",
		"architecture": "**architecture** — Architecture configuration (pattern, principles, description).",
		"pattern":      "**pattern** — Architecture pattern: `hexagonal`, `clean`, `layered`, `event-driven`, `microservices`, `modular-monolith`.",
		"deployment":   "**deployment** — Deployment configuration (strategy, environments).",
		"strategy":     "**strategy** — Deployment strategy: `blue-green`, `canary`, `rolling`, ` recreate`.",
		"testing":      "**testing** — Testing configuration (strategy, coverage target).",
		"generation":   "**generation** — Code generation settings (languages, output directory).",
		"languages":    "**languages** — Programming languages for code generation: `go`, `typescript`, `python`, `java`, `rust`.",
	}

	key := strings.TrimRight(strings.TrimSpace(line), ":")
	if desc, ok := hovers[key]; ok {
		return desc
	}

	for k, desc := range hovers {
		if strings.HasPrefix(key, k) {
			return desc
		}
	}

	return ""
}

func (s *Server) handleCompletion(params json.RawMessage) any {
	var p CompletionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil
	}

	s.mu.RLock()
	text := s.documents[p.TextDocument.URI]
	s.mu.RUnlock()

	line := s.getLine(text, p.Position.Line)
	trimmed := strings.TrimSpace(line)

	items := s.completionForLine(trimmed)

	return CompletionList{
		IsIncomplete: false,
		Items:        items,
	}
}

func (s *Server) completionForLine(line string) []CompletionItem {
	topLevel := []CompletionItem{
		{Label: "project", Kind: CompletionKindProperty, Detail: "Project name", InsertText: "project: ${1:my-project}"},
		{Label: "modules", Kind: CompletionKindProperty, Detail: "Application modules", InsertText: "modules:\n  - name: ${1:core}\n    path: ./internal/core"},
		{Label: "services", Kind: CompletionKindProperty, Detail: "Application services", InsertText: "services:\n  - name: ${1:api}\n    kind: http\n    port: ${2:8080}"},
		{Label: "architecture", Kind: CompletionKindProperty, Detail: "Architecture config", InsertText: "architecture:\n  pattern: ${1|hexagonal,clean,layered,event-driven,microservices,modular-monolith|}"},
		{Label: "deployment", Kind: CompletionKindProperty, Detail: "Deployment config", InsertText: "deployment:\n  strategy: ${1|blue-green,canary,rolling,recreate|}"},
		{Label: "testing", Kind: CompletionKindProperty, Detail: "Testing config", InsertText: "testing:\n  strategy: ${1|unit,integration,e2e|}\n  coverage: \"${2:80%}\""},
		{Label: "generation", Kind: CompletionKindProperty, Detail: "Code generation", InsertText: "generation:\n  languages:\n    - ${1|go,typescript,python,java,rust|}"},
	}

	if strings.HasPrefix(line, "- name:") || strings.HasPrefix(line, "  - name:") {
		return []CompletionItem{
			{Label: "name", Kind: CompletionKindProperty, Detail: "Module name", InsertText: "name: ${1:module-name}"},
			{Label: "path", Kind: CompletionKindProperty, Detail: "Module path", InsertText: "path: ./internal/${1:module-name}"},
			{Label: "description", Kind: CompletionKindProperty, Detail: "Module description"},
			{Label: "dependencies", Kind: CompletionKindProperty, Detail: "Module dependencies", InsertText: "dependencies:\n  - ${1:core}"},
		}
	}
	if strings.HasPrefix(line, "kind:") || strings.HasPrefix(line, "  kind:") {
		return []CompletionItem{
			{Label: "http", Kind: CompletionKindValue, Detail: "HTTP service"},
			{Label: "grpc", Kind: CompletionKindValue, Detail: "gRPC service"},
			{Label: "worker", Kind: CompletionKindValue, Detail: "Background worker"},
			{Label: "event", Kind: CompletionKindValue, Detail: "Event-driven service"},
			{Label: "graphql", Kind: CompletionKindValue, Detail: "GraphQL service"},
		}
	}
	if strings.HasPrefix(line, "method:") || strings.HasPrefix(line, "  method:") {
		return []CompletionItem{
			{Label: "GET", Kind: CompletionKindValue},
			{Label: "POST", Kind: CompletionKindValue},
			{Label: "PUT", Kind: CompletionKindValue},
			{Label: "DELETE", Kind: CompletionKindValue},
			{Label: "PATCH", Kind: CompletionKindValue},
			{Label: "HEAD", Kind: CompletionKindValue},
			{Label: "OPTIONS", Kind: CompletionKindValue},
		}
	}
	if strings.HasPrefix(line, "pattern:") {
		return []CompletionItem{
			{Label: "hexagonal", Kind: CompletionKindValue},
			{Label: "clean", Kind: CompletionKindValue},
			{Label: "layered", Kind: CompletionKindValue},
			{Label: "event-driven", Kind: CompletionKindValue},
			{Label: "microservices", Kind: CompletionKindValue},
			{Label: "modular-monolith", Kind: CompletionKindValue},
		}
	}
	if strings.HasPrefix(line, "strategy:") {
		return []CompletionItem{
			{Label: "blue-green", Kind: CompletionKindValue},
			{Label: "canary", Kind: CompletionKindValue},
			{Label: "rolling", Kind: CompletionKindValue},
			{Label: "recreate", Kind: CompletionKindValue},
		}
	}
	if strings.HasPrefix(line, "- ") && !strings.Contains(line, "name") {
		return []CompletionItem{
			{Label: "go", Kind: CompletionKindValue, Detail: "Go language"},
			{Label: "typescript", Kind: CompletionKindValue, Detail: "TypeScript language"},
			{Label: "python", Kind: CompletionKindValue, Detail: "Python language"},
			{Label: "java", Kind: CompletionKindValue, Detail: "Java language"},
			{Label: "rust", Kind: CompletionKindValue, Detail: "Rust language"},
		}
	}

	return topLevel
}

func (s *Server) getLine(text string, lineNum int) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	current := 0
	for scanner.Scan() {
		if current == lineNum {
			return scanner.Text()
		}
		current++
	}
	return ""
}

func (s *Server) sendResponse(resp Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return s.sendRaw(data)
}

func (s *Server) sendNotification(method string, params any) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	n := Notification{JSONRPC: "2.0", Method: method, Params: data}
	return s.sendRaw(json.RawMessage(mustMarshal(n)))
}

func (s *Server) sendRaw(data []byte) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	_, err := s.writer.Write([]byte(header))
	if err != nil {
		return err
	}
	_, err = s.writer.Write(data)
	return err
}

func mustMarshal(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}

func ParseMessages(r io.Reader) ([]Request, error) {
	reader := bufio.NewReader(r)
	var requests []Request

	for {
		var contentLength int
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return requests, nil
				}
				return requests, err
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			if strings.HasPrefix(line, "Content-Length: ") {
				cl := strings.TrimPrefix(line, "Content-Length: ")
				contentLength, _ = strconv.Atoi(cl)
			}
		}

		if contentLength <= 0 {
			continue
		}

		body := make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return requests, err
		}

		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			continue
		}
		requests = append(requests, req)
	}
}
