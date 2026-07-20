package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/NAEOS-foundation/naeos/internal/promptlib"
)

// LLMProvider identifies a supported LLM provider.
type LLMProvider string

const (
	// ProviderOpenAI is the OpenAI LLM provider.
	ProviderOpenAI LLMProvider = "openai"
	// ProviderAnthropic is the Anthropic LLM provider.
	ProviderAnthropic LLMProvider = "anthropic"
	// ProviderOllama is the Ollama local LLM provider (OpenAI-compatible).
	ProviderOllama LLMProvider = "ollama"
)

// LLMConfig holds configuration for an LLM service.
type LLMConfig struct {
	Provider  LLMProvider
	APIKey    string
	Model     string
	MaxTokens int
	Timeout   time.Duration
	BaseURL   string
}

// LLMService communicates with external LLM APIs to enrich specifications.
type LLMService struct {
	config     LLMConfig
	httpClient *http.Client
	library    *promptlib.Library
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type anthropicRequest struct {
	Model     string        `json:"model"`
	MaxTokens int           `json:"max_tokens"`
	Messages  []chatMessage `json:"messages"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewLLMService creates an LLM service with the given configuration.
func NewLLMService(config LLMConfig, lib ...*promptlib.Library) *LLMService {
	if config.Model == "" {
		switch config.Provider {
		case ProviderOpenAI:
			config.Model = "gpt-4o-mini"
		case ProviderAnthropic:
			config.Model = "claude-3-haiku-20240307"
		case ProviderOllama:
			config.Model = "llama3.2"
		}
	}
	if config.BaseURL == "" && config.Provider == ProviderOllama {
		config.BaseURL = "http://localhost:11434"
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 1024
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	s := &LLMService{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
	if len(lib) > 0 && lib[0] != nil {
		s.library = lib[0]
	}
	return s
}

// EnrichSpec sends a specification to the LLM for enhancement with best practices.
func (s *LLMService) EnrichSpec(specContent string) (string, error) {
	prompt := s.buildEnrichPrompt(specContent)
	return s.callLLM(prompt)
}

func (s *LLMService) buildEnrichPrompt(specContent string) string {
	if s.library != nil {
		rendered, err := s.library.RenderLLM("enrich-spec", map[string]any{
			"SpecContent": specContent,
		})
		if err == nil && rendered.User != "" {
			return rendered.User
		}
	}

	return fmt.Sprintf(`You are a platform engineering expert. Analyze this NAEOS specification and enrich it with best practices.
Add any missing sections that would improve the specification. Keep the existing content intact.
Only output the enriched YAML specification, no explanations.

Specification:
%s`, specContent)
}

// GenerateSuggestions asks the LLM to produce improvement suggestions for a specification.
func (s *LLMService) GenerateSuggestions(specContent string) ([]Suggestion, error) {
	prompt := s.buildSuggestionsPrompt(specContent)

	response, err := s.callLLM(prompt)
	if err != nil {
		return nil, err
	}

	var suggestions []Suggestion
	if err := json.Unmarshal([]byte(cleanJSON(response)), &suggestions); err != nil {
		return nil, fmt.Errorf("parse LLM response: %w", err)
	}

	return suggestions, nil
}

func (s *LLMService) buildSuggestionsPrompt(specContent string) string {
	if s.library != nil {
		rendered, err := s.library.RenderLLM("generate-suggestions", map[string]any{
			"SpecContent": specContent,
		})
		if err == nil && rendered.User != "" {
			return rendered.User
		}
	}

	return fmt.Sprintf(`Analyze this NAEOS specification and return a JSON array of suggestions.
Each suggestion should have: category, title, description, priority (high/medium/low).
Return ONLY the JSON array, no other text.

Specification:
%s`, specContent)
}

// ExplainArchitecture asks the LLM to explain an architecture pattern in the context of the specification.
func (s *LLMService) ExplainArchitecture(specContent, architecture string) (string, error) {
	prompt := s.buildExplainPrompt(specContent, architecture)
	return s.callLLM(prompt)
}

func (s *LLMService) buildExplainPrompt(specContent, arch string) string {
	if s.library != nil {
		rendered, err := s.library.RenderLLM("explain-architecture", map[string]any{
			"SpecContent":  specContent,
			"Architecture": arch,
		})
		if err == nil && rendered.User != "" {
			return rendered.User
		}
	}

	return fmt.Sprintf(`Explain the architecture pattern "%s" in the context of this specification.
Provide a clear, concise explanation suitable for a developer.

Specification:
%s

Architecture explanation:`, arch, specContent)
}

// modelContextWindows maps known model names to their context window sizes (in tokens).
var modelContextWindows = map[string]int{
	"gpt-4o":                128000,
	"gpt-4o-mini":           128000,
	"gpt-4-turbo":           128000,
	"gpt-4":                 8192,
	"gpt-3.5-turbo":         16385,
	"claude-3-opus-20240229":   200000,
	"claude-3-sonnet-20240229": 200000,
	"claude-3-haiku-20240307":  200000,
	"claude-3-5-sonnet-20241022": 200000,
	"claude-3-5-haiku-20241022": 200000,
	"llama3.2":              8192,
	"llama3.1":              8192,
	"llama3":                8192,
	"mistral":               8192,
	"codellama":             16384,
	"mixtral":               32768,
}

// estimateTokens returns a rough estimate of the number of tokens in a string.
// Uses the common heuristic of ~4 characters per token for English text.
func estimateTokens(s string) int {
	return len(s) / 4
}

// truncatePrompt truncates the prompt to fit within the model's context window,
// reserving config.MaxTokens output tokens.
func (s *LLMService) truncatePrompt(prompt string) string {
	window, ok := modelContextWindows[s.config.Model]
	if !ok {
		window = 8192
	}

	available := window - s.config.MaxTokens
	if available < 256 {
		available = 256
	}

	estimated := estimateTokens(prompt)
	if estimated <= available {
		return prompt
	}

	maxChars := available * 4
	truncated := prompt[:maxChars]

	slog.Warn("prompt truncated",
		"model", s.config.Model,
		"estimated_tokens", estimated,
		"context_window", window,
		"available_tokens", available,
	)

	return truncated
}

func (s *LLMService) callLLM(prompt string) (string, error) {
	prompt = s.truncatePrompt(prompt)
	switch s.config.Provider {
	case ProviderOpenAI, ProviderOllama:
		response, err := s.callOpenAI(prompt)
		if err != nil {
			slog.Error("openai call failed", "error", err)
		}
		return response, err
	case ProviderAnthropic:
		response, err := s.callAnthropic(prompt)
		if err != nil {
			slog.Error("anthropic call failed", "error", err)
		}
		return response, err
	default:
		slog.Error("unsupported LLM provider", "provider", s.config.Provider)
		return "", fmt.Errorf("unsupported LLM provider: %s", s.config.Provider)
	}
}

func (s *LLMService) callOpenAI(prompt string) (string, error) {
	reqBody := openAIRequest{
		Model: s.config.Model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   s.config.MaxTokens,
		Temperature: 0.3,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	baseURL := s.config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	url := baseURL + "/v1/chat/completions"

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		slog.Error("openai request failed", "error", err)
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("openai read response failed", "error", err)
		return "", err
	}

	var result openAIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		slog.Error("openai parse response failed", "error", err)
		return "", err
	}

	if result.Error != nil {
		slog.Error("openai api error", "message", result.Error.Message)
		return "", fmt.Errorf("openai error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		slog.Error("openai no choices returned")
		return "", fmt.Errorf("openai: no choices returned")
	}

	slog.Info("openai call succeeded", "model", s.config.Model)
	return result.Choices[0].Message.Content, nil
}

func (s *LLMService) callAnthropic(prompt string) (string, error) {
	reqBody := anthropicRequest{
		Model: s.config.Model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens: s.config.MaxTokens,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	baseURL := s.config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	url := baseURL + "/v1/messages"

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result anthropicResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("anthropic error: %s", result.Error.Message)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic: no content returned")
	}

	return result.Content[0].Text, nil
}

func (s *LLMService) StreamEnrichSpec(specContent string, w io.Writer) error {
	prompt := s.buildEnrichPrompt(specContent)
	return s.streamLLM(prompt, w)
}

func (s *LLMService) StreamExplainArchitecture(specContent, architecture string, w io.Writer) error {
	prompt := s.buildExplainPrompt(specContent, architecture)
	return s.streamLLM(prompt, w)
}

func (s *LLMService) streamLLM(prompt string, w io.Writer) error {
	flusher, flushable := w.(http.Flusher)

	writeEvent := func(event, data string) error {
		_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
		if flushable {
			flusher.Flush()
		}
		return err
	}

	if err := writeEvent("start", "{}"); err != nil {
		return err
	}

	result, err := s.callLLM(prompt)
	if err != nil {
		_ = writeEvent("error", fmt.Sprintf(`{"message":"%s"}`, err.Error()))
		return err
	}

	words := strings.Fields(result)
	var buf strings.Builder
	for _, word := range words {
		buf.WriteString(word)
		buf.WriteByte(' ')
		chunk := strings.TrimSpace(buf.String())
		if len(chunk) >= 80 || word == words[len(words)-1] {
			escaped := strings.ReplaceAll(chunk, "\n", "\\n")
			escaped = strings.ReplaceAll(escaped, `"`, `\"`)
			if err := writeEvent("chunk", fmt.Sprintf(`{"text":"%s"}`, escaped)); err != nil {
				return err
			}
			buf.Reset()
		}
	}

	if buf.Len() > 0 {
		remaining := strings.TrimSpace(buf.String())
		escaped := strings.ReplaceAll(remaining, "\n", "\\n")
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		_ = writeEvent("chunk", fmt.Sprintf(`{"text":"%s"}`, escaped))
	}

	return writeEvent("done", `{"status":"completed"}`)
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
