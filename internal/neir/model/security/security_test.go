package security

import "testing"

func TestZeroValue(t *testing.T) {
	var s Security
	if s.Authentication != nil {
		t.Errorf("expected nil Authentication, got %v", s.Authentication)
	}
	if s.Authorization != nil {
		t.Errorf("expected nil Authorization, got %v", s.Authorization)
	}
	if s.Encryption != nil {
		t.Errorf("expected nil Encryption, got %v", s.Encryption)
	}
	if s.Secrets != nil {
		t.Errorf("expected nil Secrets, got %v", s.Secrets)
	}

	var auth Authentication
	if auth.Method != "" {
		t.Errorf("expected empty Method, got %q", auth.Method)
	}

	var az Authorization
	if az.Model != "" {
		t.Errorf("expected empty Model, got %q", az.Model)
	}
	if az.Roles != nil {
		t.Errorf("expected nil Roles, got %v", az.Roles)
	}

	var enc Encryption
	if enc.InTransit {
		t.Error("expected false InTransit")
	}
	if enc.AtRest {
		t.Error("expected false AtRest")
	}
	if enc.Algorithm != "" {
		t.Errorf("expected empty Algorithm, got %q", enc.Algorithm)
	}

	var sec Secret
	if sec.Name != "" {
		t.Errorf("expected empty Name, got %q", sec.Name)
	}
}

func TestInitialization(t *testing.T) {
	s := Security{
		Authentication: &Authentication{Method: "jwt", Provider: "auth0"},
		Authorization:  &Authorization{Model: "rbac", Roles: []string{"admin", "viewer"}},
		Encryption:     &Encryption{InTransit: true, AtRest: true, Algorithm: "AES-256"},
		Secrets: []Secret{
			{Name: "db-password", Kind: "env"},
			{Name: "api-key", Kind: "vault"},
		},
	}

	if s.Authentication == nil || s.Authentication.Method != "jwt" {
		t.Errorf("expected Authentication.Method 'jwt', got %v", s.Authentication)
	}
	if s.Authorization == nil || s.Authorization.Model != "rbac" {
		t.Errorf("expected Authorization.Model 'rbac', got %v", s.Authorization)
	}
	if !s.Encryption.InTransit || !s.Encryption.AtRest {
		t.Error("expected both InTransit and AtRest to be true")
	}
	if len(s.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(s.Secrets))
	}
}

func TestSecurityNilPointers(t *testing.T) {
	var s Security
	if s.Authentication != nil || s.Authorization != nil || s.Encryption != nil {
		t.Error("expected all pointer fields to be nil")
	}
	s.Authentication = &Authentication{Method: "oauth2"}
	s.Encryption = &Encryption{Algorithm: "RSA"}
	if s.Authentication == nil || s.Encryption == nil {
		t.Error("expected non-nil after assignment")
	}
	if s.Authorization != nil {
		t.Error("expected Authorization to still be nil")
	}
}
