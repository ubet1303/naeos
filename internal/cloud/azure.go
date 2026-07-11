package cloud

import (
	"fmt"
	"strings"
	"time"
)

type AzureAdapter struct{}

func (a *AzureAdapter) Name() string {
	return "Azure"
}

func (a *AzureAdapter) Provider() CloudProvider {
	return Azure
}

func (a *AzureAdapter) Validate(config *DeployConfig) error {
	if config.Project == "" {
		return fmt.Errorf("Azure resource group is required")
	}
	if config.Region == "" {
		return fmt.Errorf("Azure region is required")
	}
	return nil
}

func (a *AzureAdapter) Plan(config *DeployConfig) ([]Resource, error) {
	resources := []Resource{}
	for _, res := range config.Resources {
		switch res.Type {
		case "storage":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "azurerm_storage_account",
				Spec: map[string]interface{}{
					"name":     fmt.Sprintf("%s%s%s", config.Project, config.Environment, res.Name),
					"location": config.Region,
				},
			})
		case "compute":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "azurerm_container_group",
				Spec: map[string]interface{}{
					"name":     res.Name,
					"location": config.Region,
				},
			})
		case "database":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "azurerm_postgresql_flexible_server",
				Spec: map[string]interface{}{
					"name":     fmt.Sprintf("%s-%s-%s", config.Project, config.Environment, res.Name),
					"location": config.Region,
				},
			})
		}
	}
	return resources, nil
}

func (a *AzureAdapter) Deploy(config *DeployConfig) (*DeployResult, error) {
	plan, err := a.Plan(config)
	if err != nil {
		return nil, err
	}

	deployed := []DeployedResource{}
	for _, res := range plan {
		deployed = append(deployed, DeployedResource{
			Name: res.Name,
			Type: res.Type,
			ID:   fmt.Sprintf("/subscriptions/.../resourceGroups/%s/providers/%s/%s", config.Project, res.Type, res.Name),
		})
	}

	tf, _ := a.ExportTerraform(config)

	return &DeployResult{
		Provider:  Azure,
		Resources: deployed,
		Terraform: tf,
		Status:    "deployed",
		Timestamp: time.Now(),
	}, nil
}

func (a *AzureAdapter) Destroy(config *DeployConfig) error {
	return nil
}

func (a *AzureAdapter) ExportTerraform(config *DeployConfig) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "main" {
  name     = "%s"
  location = "%s"
}

`, config.Project, config.Region))

	for _, res := range config.Resources {
		switch res.Type {
		case "storage":
			sb.WriteString(fmt.Sprintf(`resource "azurerm_storage_account" "%s" {
  name                     = "%s%s%s"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "%s"
    project     = "%s"
  }
}

`, res.Name, config.Project, config.Environment, res.Name, config.Environment, config.Project))
		}
	}

	return sb.String(), nil
}
