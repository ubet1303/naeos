package cloud

import (
	"fmt"
	"strings"
	"time"
)

type GCPAdapter struct{}

func (a *GCPAdapter) Name() string {
	return "GCP"
}

func (a *GCPAdapter) Provider() CloudProvider {
	return GCP
}

func (a *GCPAdapter) Validate(config *DeployConfig) error {
	if config.Project == "" {
		return fmt.Errorf("GCP project is required")
	}
	if config.Region == "" {
		return fmt.Errorf("GCP region is required")
	}
	return nil
}

func (a *GCPAdapter) Plan(config *DeployConfig) ([]Resource, error) {
	resources := []Resource{}
	for _, res := range config.Resources {
		switch res.Type {
		case "storage":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "google_storage_bucket",
				Spec: map[string]interface{}{
					"name":     fmt.Sprintf("%s-%s-%s", config.Project, config.Environment, res.Name),
					"location": config.Region,
				},
			})
		case "compute":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "google_cloud_run_service",
				Spec: map[string]interface{}{
					"name":     res.Name,
					"location": config.Region,
				},
			})
		case "database":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "google_sql_database_instance",
				Spec: map[string]interface{}{
					"name":    fmt.Sprintf("%s-%s-%s", config.Project, config.Environment, res.Name),
					"region":  config.Region,
					"db_version": "POSTGRES_15",
				},
			})
		}
	}
	return resources, nil
}

func (a *GCPAdapter) Deploy(config *DeployConfig) (*DeployResult, error) {
	plan, err := a.Plan(config)
	if err != nil {
		return nil, err
	}

	deployed := []DeployedResource{}
	for _, res := range plan {
		deployed = append(deployed, DeployedResource{
			Name: res.Name,
			Type: res.Type,
			ID:   fmt.Sprintf("projects/%s/%s/%s", config.Project, res.Type, res.Name),
		})
	}

	tf, _ := a.ExportTerraform(config)

	return &DeployResult{
		Provider:  GCP,
		Resources: deployed,
		Terraform: tf,
		Status:    "deployed",
		Timestamp: time.Now(),
	}, nil
}

func (a *GCPAdapter) Destroy(config *DeployConfig) error {
	return nil
}

func (a *GCPAdapter) ExportTerraform(config *DeployConfig) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = "%s"
  region  = "%s"
}

`, config.Project, config.Region))

	for _, res := range config.Resources {
		switch res.Type {
		case "storage":
			sb.WriteString(fmt.Sprintf(`resource "google_storage_bucket" "%s" {
  name     = "%s-%s-%s"
  location = "%s"

  labels = {
    environment = "%s"
    project     = "%s"
  }
}

`, res.Name, config.Project, config.Environment, res.Name, config.Region, config.Environment, config.Project))
		}
	}

	return sb.String(), nil
}
