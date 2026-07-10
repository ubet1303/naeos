package infrastructure

import "testing"

func TestProviderConstants(t *testing.T) {
	tests := []struct {
		constant Provider
		expected string
	}{
		{ProviderAWS, "aws"},
		{ProviderGCP, "gcp"},
		{ProviderAzure, "azure"},
		{ProviderLocal, "local"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("Provider %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var infra Infrastructure
	if infra.Provider != "" {
		t.Errorf("expected empty Provider, got %q", infra.Provider)
	}
	if infra.Region != "" {
		t.Errorf("expected empty Region, got %q", infra.Region)
	}
	if infra.Resources != nil {
		t.Errorf("expected nil Resources, got %v", infra.Resources)
	}
	if infra.Networking != nil {
		t.Errorf("expected nil Networking, got %v", infra.Networking)
	}

	var r Resource
	if r.Name != "" {
		t.Errorf("expected empty Name, got %q", r.Name)
	}
	if r.Spec != nil {
		t.Errorf("expected nil Spec, got %v", r.Spec)
	}

	var n Network
	if n.Name != "" {
		t.Errorf("expected empty Name, got %q", n.Name)
	}
	if n.Ports != nil {
		t.Errorf("expected nil Ports, got %v", n.Ports)
	}
}

func TestInitialization(t *testing.T) {
	infra := Infrastructure{
		Provider: ProviderAWS,
		Region:   "us-east-1",
		Resources: []Resource{
			{Name: "main-db", Kind: "rds", Spec: map[string]string{"engine": "postgres"}},
		},
		Networking: []Network{
			{Name: "vpc", Kind: "private", Ports: []int{5432, 8080}},
		},
		Attributes: map[string]string{"cost-center": "engineering"},
	}

	if infra.Provider != ProviderAWS {
		t.Errorf("expected Provider %q, got %q", ProviderAWS, infra.Provider)
	}
	if infra.Resources[0].Spec["engine"] != "postgres" {
		t.Errorf("expected engine=postgres, got %q", infra.Resources[0].Spec["engine"])
	}
	if infra.Networking[0].Ports[0] != 5432 {
		t.Errorf("expected first port 5432, got %d", infra.Networking[0].Ports[0])
	}
}
