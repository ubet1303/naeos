package contracts

import (
	"testing"
)

type testContract struct {
	valid bool
}

func (c testContract) Validate() error {
	if !c.valid {
		return errInvalid
	}
	return nil
}

type testSchemaAware struct {
	version string
}

func (s testSchemaAware) SchemaVersion() string {
	return s.version
}

type testVersioned struct {
	version string
}

func (v testVersioned) Version() string {
	return v.version
}

type testIdentifiable struct {
	id string
}

func (i testIdentifiable) ID() string {
	return i.id
}

type testNamed struct {
	name string
}

func (n testNamed) Name() string {
	return n.name
}

var errInvalid = &testError{"invalid contract"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestContractInterface(t *testing.T) {
	var c Contract = testContract{valid: true}
	if err := c.Validate(); err != nil {
		t.Errorf("valid contract should return nil, got %v", err)
	}

	c = testContract{valid: false}
	if err := c.Validate(); err == nil {
		t.Error("invalid contract should return error")
	}
}

func TestSchemaAwareInterface(t *testing.T) {
	var s SchemaAware = testSchemaAware{version: "1.0.0"}
	if got := s.SchemaVersion(); got != "1.0.0" {
		t.Errorf("SchemaVersion() = %q, want %q", got, "1.0.0")
	}
}

func TestVersionedInterface(t *testing.T) {
	var v Versioned = testVersioned{version: "2.1.0"}
	if got := v.Version(); got != "2.1.0" {
		t.Errorf("Version() = %q, want %q", got, "2.1.0")
	}
}

func TestIdentifiableInterface(t *testing.T) {
	var i Identifiable = testIdentifiable{id: "ID-001"}
	if got := i.ID(); got != "ID-001" {
		t.Errorf("ID() = %q, want %q", got, "ID-001")
	}
}

func TestNamedInterface(t *testing.T) {
	var n Named = testNamed{name: "my-service"}
	if got := n.Name(); got != "my-service" {
		t.Errorf("Name() = %q, want %q", got, "my-service")
	}
}
