package cloud

import (
	"testing"
)

func TestGetAdapterAWS(t *testing.T) {
	adapter, err := GetAdapter(AWS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if adapter.Name() != "AWS" {
		t.Errorf("expected name 'AWS', got %s", adapter.Name())
	}
}

func TestGetAdapterGCP(t *testing.T) {
	adapter, err := GetAdapter(GCP)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if adapter.Name() != "GCP" {
		t.Errorf("expected name 'GCP', got %s", adapter.Name())
	}
}

func TestGetAdapterAzure(t *testing.T) {
	adapter, err := GetAdapter(Azure)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if adapter.Name() != "Azure" {
		t.Errorf("expected name 'Azure', got %s", adapter.Name())
	}
}

func TestGetAdapterInvalid(t *testing.T) {
	_, err := GetAdapter("invalid")
	if err == nil {
		t.Error("expected error for invalid provider")
	}
}

func TestAWSValidate(t *testing.T) {
	adapter := &AWSAdapter{}

	validConfig := &DeployConfig{
		Provider: AWS,
		Region:   "us-east-1",
		Project:  "test",
	}
	if err := adapter.Validate(validConfig); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	invalidConfig := &DeployConfig{
		Provider: AWS,
		Region:   "invalid-region",
		Project:  "test",
	}
	if err := adapter.Validate(invalidConfig); err == nil {
		t.Error("expected error for invalid region")
	}
}

func TestGCPValidate(t *testing.T) {
	adapter := &GCPAdapter{}

	validConfig := &DeployConfig{
		Provider: GCP,
		Project:  "my-project",
		Region:   "us-central1",
	}
	if err := adapter.Validate(validConfig); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	invalidConfig := &DeployConfig{
		Provider: GCP,
		Region:   "us-central1",
	}
	if err := adapter.Validate(invalidConfig); err == nil {
		t.Error("expected error for missing project")
	}
}

func TestAzureValidate(t *testing.T) {
	adapter := &AzureAdapter{}

	validConfig := &DeployConfig{
		Provider: Azure,
		Project:  "my-rg",
		Region:   "eastus",
	}
	if err := adapter.Validate(validConfig); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	invalidConfig := &DeployConfig{
		Provider: Azure,
		Region:   "eastus",
	}
	if err := adapter.Validate(invalidConfig); err == nil {
		t.Error("expected error for missing project")
	}
}

func TestAWSPlan(t *testing.T) {
	adapter := &AWSAdapter{}
	config := &DeployConfig{
		Provider:    AWS,
		Region:      "us-east-1",
		Project:     "myapp",
		Environment: "prod",
		Resources: []Resource{
			{Name: "uploads", Type: "storage"},
			{Name: "api", Type: "compute"},
			{Name: "db", Type: "database"},
		},
	}

	plan, err := adapter.Plan(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plan) != 3 {
		t.Errorf("expected 3 resources, got %d", len(plan))
	}

	if plan[0].Type != "aws_s3_bucket" {
		t.Errorf("expected aws_s3_bucket, got %s", plan[0].Type)
	}
}

func TestAWSTerraformExport(t *testing.T) {
	adapter := &AWSAdapter{}
	config := &DeployConfig{
		Provider:    AWS,
		Region:      "us-east-1",
		Project:     "myapp",
		Environment: "prod",
		Resources: []Resource{
			{Name: "uploads", Type: "storage"},
		},
	}

	tf, err := adapter.ExportTerraform(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tf == "" {
		t.Error("expected non-empty terraform output")
	}
}
