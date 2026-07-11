package cloud

import (
	"fmt"
	"strings"
	"time"
)

type AWSAdapter struct{}

func (a *AWSAdapter) Name() string {
	return "AWS"
}

func (a *AWSAdapter) Provider() CloudProvider {
	return AWS
}

func (a *AWSAdapter) Validate(config *DeployConfig) error {
	if config.Region == "" {
		return fmt.Errorf("AWS region is required")
	}
	validRegions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}
	valid := false
	for _, r := range validRegions {
		if config.Region == r {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid AWS region: %s", config.Region)
	}
	return nil
}

func (a *AWSAdapter) Plan(config *DeployConfig) ([]Resource, error) {
	resources := []Resource{}
	for _, res := range config.Resources {
		switch res.Type {
		case "storage":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "aws_s3_bucket",
				Spec: map[string]interface{}{
					"bucket": fmt.Sprintf("%s-%s-%s", config.Project, config.Environment, res.Name),
					"region": config.Region,
				},
			})
		case "compute":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "aws_ecs_service",
				Spec: map[string]interface{}{
					"cluster": fmt.Sprintf("%s-%s", config.Project, config.Environment),
					"service": res.Name,
				},
			})
		case "database":
			resources = append(resources, Resource{
				Name: res.Name,
				Type: "aws_rds_instance",
				Spec: map[string]interface{}{
					"identifier": fmt.Sprintf("%s-%s-%s", config.Project, config.Environment, res.Name),
					"engine":     "postgres",
				},
			})
		}
	}
	return resources, nil
}

func (a *AWSAdapter) Deploy(config *DeployConfig) (*DeployResult, error) {
	plan, err := a.Plan(config)
	if err != nil {
		return nil, err
	}

	deployed := []DeployedResource{}
	for _, res := range plan {
		deployed = append(deployed, DeployedResource{
			Name: res.Name,
			Type: res.Type,
			ID:   fmt.Sprintf("arn:aws:%s:%s:%s", res.Type, config.Region, res.Name),
		})
	}

	tf, _ := a.ExportTerraform(config)

	return &DeployResult{
		Provider:  AWS,
		Resources: deployed,
		Terraform: tf,
		Status:    "deployed",
		Timestamp: time.Now(),
	}, nil
}

func (a *AWSAdapter) Destroy(config *DeployConfig) error {
	return nil
}

func (a *AWSAdapter) ExportTerraform(config *DeployConfig) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "%s"
}

`, config.Region))

	for _, res := range config.Resources {
		switch res.Type {
		case "storage":
			sb.WriteString(fmt.Sprintf(`resource "aws_s3_bucket" "%s" {
  bucket = "%s-%s-%s"

  tags = {
    Environment = "%s"
    Project     = "%s"
  }
}

`, res.Name, config.Project, config.Environment, res.Name, config.Environment, config.Project))
		}
	}

	return sb.String(), nil
}
