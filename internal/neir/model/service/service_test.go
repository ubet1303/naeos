package service

import (
	"testing"
)

func TestServiceKindConstants(t *testing.T) {
	tests := []struct {
		kind ServiceKind
		want string
	}{
		{KindHTTP, "http"},
		{KindGRPC, "grpc"},
		{KindWorker, "worker"},
		{KindCLI, "cli"},
		{KindJob, "job"},
	}
	for _, tt := range tests {
		if string(tt.kind) != tt.want {
			t.Errorf("ServiceKind = %q, want %q", tt.kind, tt.want)
		}
	}
}

func TestServiceZeroValue(t *testing.T) {
	s := &Service{}
	if s.Name != "" {
		t.Errorf("zero-value Service.Name = %q, want empty", s.Name)
	}
	if s.Port != 0 {
		t.Errorf("zero-value Service.Port = %d, want 0", s.Port)
	}
}

func TestServiceWithEndpoints(t *testing.T) {
	s := &Service{
		Name: "gateway",
		Kind: KindHTTP,
		Port: 8080,
		Endpoints: []Endpoint{
			{Method: "GET", Path: "/health", Action: "healthCheck"},
			{Method: "POST", Path: "/api/v1/users", Action: "createUser"},
		},
	}
	if len(s.Endpoints) != 2 {
		t.Errorf("Endpoints has %d entries, want 2", len(s.Endpoints))
	}
	if s.Endpoints[0].Method != "GET" {
		t.Errorf("Endpoints[0].Method = %q, want %q", s.Endpoints[0].Method, "GET")
	}
	if s.Endpoints[1].Path != "/api/v1/users" {
		t.Errorf("Endpoints[1].Path = %q, want %q", s.Endpoints[1].Path, "/api/v1/users")
	}
}

func TestServiceWithMiddleware(t *testing.T) {
	s := &Service{
		Name:       "api",
		Middleware: []string{"logging", "auth", "cors"},
	}
	if len(s.Middleware) != 3 {
		t.Errorf("Middleware has %d entries, want 3", len(s.Middleware))
	}
}
