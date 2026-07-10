package domain

import "testing"

func TestZeroValue(t *testing.T) {
	var d Domain
	if d.Name != "" {
		t.Errorf("expected empty Name, got %q", d.Name)
	}
	if d.BoundedContexts != nil {
		t.Errorf("expected nil BoundedContexts, got %v", d.BoundedContexts)
	}
	if d.Aggregates != nil {
		t.Errorf("expected nil Aggregates, got %v", d.Aggregates)
	}
	if d.Entities != nil {
		t.Errorf("expected nil Entities, got %v", d.Entities)
	}
	if d.ValueObjects != nil {
		t.Errorf("expected nil ValueObjects, got %v", d.ValueObjects)
	}

	var bc BoundedContext
	if bc.Name != "" {
		t.Errorf("expected empty Name, got %q", bc.Name)
	}
	if bc.Modules != nil {
		t.Errorf("expected nil Modules, got %v", bc.Modules)
	}

	var agg Aggregate
	if agg.Name != "" {
		t.Errorf("expected empty Name, got %q", agg.Name)
	}
	if agg.Entities != nil {
		t.Errorf("expected nil Entities, got %v", agg.Entities)
	}

	var e Entity
	if e.Name != "" {
		t.Errorf("expected empty Name, got %q", e.Name)
	}
	if e.Attributes != nil {
		t.Errorf("expected nil Attributes, got %v", e.Attributes)
	}

	var vo ValueObject
	if vo.Name != "" {
		t.Errorf("expected empty Name, got %q", vo.Name)
	}
}

func TestInitialization(t *testing.T) {
	d := Domain{
		Name:        "ecommerce",
		Description: "E-commerce bounded contexts",
		BoundedContexts: []BoundedContext{
			{Name: "catalog", Modules: []string{"product", "category"}},
		},
		Aggregates: []Aggregate{
			{Name: "Order", RootEntity: "order", Entities: []string{"OrderItem", "ShippingAddress"}},
		},
		Entities: []Entity{
			{Name: "Order", Attributes: map[string]string{"id": "uuid", "total": "decimal"}},
		},
		ValueObjects: []ValueObject{
			{Name: "Money", Attributes: map[string]string{"amount": "decimal", "currency": "string"}},
		},
	}

	if d.Name != "ecommerce" {
		t.Errorf("expected Name 'ecommerce', got %q", d.Name)
	}
	if len(d.BoundedContexts) != 1 || d.BoundedContexts[0].Modules[0] != "product" {
		t.Errorf("unexpected BoundedContexts: %v", d.BoundedContexts)
	}
	if d.Aggregates[0].RootEntity != "order" {
		t.Errorf("expected RootEntity 'order', got %q", d.Aggregates[0].RootEntity)
	}
	if d.Entities[0].Attributes["id"] != "uuid" {
		t.Errorf("expected Entity attribute id=uuid, got %q", d.Entities[0].Attributes["id"])
	}
	if d.ValueObjects[0].Name != "Money" {
		t.Errorf("expected ValueObject name 'Money', got %q", d.ValueObjects[0].Name)
	}
}
