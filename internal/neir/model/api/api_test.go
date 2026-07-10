package api

import "testing"

func TestProtocolConstants(t *testing.T) {
	tests := []struct {
		constant Protocol
		expected string
	}{
		{ProtocolHTTP, "http"},
		{ProtocolGRPC, "grpc"},
		{ProtocolGraphQL, "graphql"},
		{ProtocolWS, "websocket"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("Protocol %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var a API
	if a.Name != "" {
		t.Errorf("expected empty Name, got %q", a.Name)
	}
	if a.Protocol != "" {
		t.Errorf("expected empty Protocol, got %q", a.Protocol)
	}
	if a.Endpoints != nil {
		t.Errorf("expected nil Endpoints, got %v", a.Endpoints)
	}
	if a.Schemas != nil {
		t.Errorf("expected nil Schemas, got %v", a.Schemas)
	}

	var ep APIEndpoint
	if ep.Method != "" {
		t.Errorf("expected empty Method, got %q", ep.Method)
	}
	if ep.Path != "" {
		t.Errorf("expected empty Path, got %q", ep.Path)
	}

	var s Schema
	if s.Name != "" {
		t.Errorf("expected empty Name, got %q", s.Name)
	}
	if s.Fields != nil {
		t.Errorf("expected nil Fields, got %v", s.Fields)
	}
}

func TestInitialization(t *testing.T) {
	a := API{
		Name:     "users-api",
		Version:  "v1",
		Protocol: ProtocolGRPC,
		Endpoints: []APIEndpoint{
			{Method: "GET", Path: "/users", Summary: "List users"},
		},
		Schemas: []Schema{
			{Name: "User", Fields: map[string]string{"id": "string", "name": "string"}},
		},
		Attributes: map[string]string{"team": "backend"},
	}

	if a.Protocol != ProtocolGRPC {
		t.Errorf("expected Protocol %q, got %q", ProtocolGRPC, a.Protocol)
	}
	if a.Endpoints[0].Method != "GET" {
		t.Errorf("expected Method GET, got %q", a.Endpoints[0].Method)
	}
	if a.Schemas[0].Fields["id"] != "string" {
		t.Errorf("expected field id=string, got %q", a.Schemas[0].Fields["id"])
	}
}
