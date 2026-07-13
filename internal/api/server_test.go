package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})
	if s == nil {
		t.Fatal("expected server to be created")
	}
	if s.Addr != ":8080" {
		t.Errorf("expected addr ':8080', got %s", s.Addr)
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	s.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}
}

func TestSpecsEndpointGET(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("GET", "/api/v1/specs", nil)
	w := httptest.NewRecorder()

	s.handleSpecs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSpecsEndpointPOST(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	body, _ := json.Marshal(map[string]string{
		"spec": "project: test\nmodules:\n  - name: core\n    path: ./core\n",
	})
	req := httptest.NewRequest("POST", "/api/v1/specs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSpecs(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestSpecValidateEndpoint(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	body, _ := json.Marshal(map[string]string{
		"spec": "project: test\nmodules:\n  - name: core\n    path: ./core\n",
	})
	req := httptest.NewRequest("POST", "/api/v1/specs/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSpecValidate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSpecValidateEndpointInvalid(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	body, _ := json.Marshal(map[string]string{
		"spec": "",
	})
	req := httptest.NewRequest("POST", "/api/v1/specs/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSpecValidate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp APIResponse
	json.NewDecoder(w.Body).Decode(&resp)
	data, _ := json.Marshal(resp.Data)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	if result["valid"].(bool) {
		t.Error("expected valid to be false for empty spec")
	}
}

func TestPipelineRunEndpoint(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	body, _ := json.Marshal(map[string]string{
		"spec": "project: test\nmodules:\n  - name: core\n    path: ./core\n",
	})
	req := httptest.NewRequest("POST", "/api/v1/pipeline/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handlePipelineRun(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("job_id")) {
		t.Errorf("expected job_id in response, got %s", w.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("DELETE", "/api/v1/specs", nil)
	w := httptest.NewRecorder()

	s.handleSpecs(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestOIDCDiscovery(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: true, JWTSecret: "test-secret"})

	req := httptest.NewRequest("GET", "/.well-known/openid-configuration", nil)
	req.Host = "localhost:8080"
	w := httptest.NewRecorder()

	s.handleOIDCDiscovery(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var doc OIDCDiscovery
	if err := json.NewDecoder(w.Body).Decode(&doc); err != nil {
		t.Fatalf("failed to decode OIDC discovery: %v", err)
	}

	if doc.Issuer != "http://localhost:8080" {
		t.Errorf("expected issuer 'http://localhost:8080', got %s", doc.Issuer)
	}

	if doc.JWKSURI != "http://localhost:8080/.well-known/jwks.json" {
		t.Errorf("expected jwks_uri 'http://localhost:8080/.well-known/jwks.json', got %s", doc.JWKSURI)
	}

	if len(doc.IDTokenSigningAlgValuesSupported) != 1 || doc.IDTokenSigningAlgValuesSupported[0] != "HS256" {
		t.Errorf("expected HS256 signing alg, got %v", doc.IDTokenSigningAlgValuesSupported)
	}
}

func TestOIDCDiscoveryNotConfigured(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: false})

	req := httptest.NewRequest("GET", "/.well-known/openid-configuration", nil)
	w := httptest.NewRecorder()

	s.handleOIDCDiscovery(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestJWKS(t *testing.T) {
	s := NewServer(":8080", &AuthConfig{Enabled: true, JWTSecret: "test-secret"})

	req := httptest.NewRequest("GET", "/.well-known/jwks.json", nil)
	w := httptest.NewRecorder()

	s.handleJWKS(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var jwks JWKS
	if err := json.NewDecoder(w.Body).Decode(&jwks); err != nil {
		t.Fatalf("failed to decode JWKS: %v", err)
	}

	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}

	key := jwks.Keys[0]
	if key.Kty != "oct" {
		t.Errorf("expected kty 'oct', got %s", key.Kty)
	}
	if key.Alg != "HS256" {
		t.Errorf("expected alg 'HS256', got %s", key.Alg)
	}
	if key.Use != "sig" {
		t.Errorf("expected use 'sig', got %s", key.Use)
	}
}
