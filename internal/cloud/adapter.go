package cloud

import (
	"fmt"
	"time"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	GCP   CloudProvider = "gcp"
	Azure CloudProvider = "azure"
)

// ResourceTypes maps abstract resource types to supported kinds.
const (
	ResourceStorage  = "storage"
	ResourceCompute  = "compute"
	ResourceDatabase = "database"
	ResourceCache    = "cache"
	ResourceQueue    = "queue"
	ResourceCDN      = "cdn"
)

var SupportedResourceTypes = []string{
	ResourceStorage,
	ResourceCompute,
	ResourceDatabase,
	ResourceCache,
	ResourceQueue,
	ResourceCDN,
}

type DeployConfig struct {
	Provider    CloudProvider
	Region      string
	Project     string
	Environment string
	Resources   []Resource
}

type Resource struct {
	Name string
	Type string
	Spec map[string]interface{}
}

type DeployResult struct {
	Provider   CloudProvider
	Resources  []DeployedResource
	Terraform  string
	Status     string
	Timestamp  time.Time
}

type DeployedResource struct {
	Name string
	Type string
	ID   string
	ARN  string
}

type CloudAdapter interface {
	Name() string
	Provider() CloudProvider
	Validate(config *DeployConfig) error
	Plan(config *DeployConfig) ([]Resource, error)
	Deploy(config *DeployConfig) (*DeployResult, error)
	Destroy(config *DeployConfig) error
	ExportTerraform(config *DeployConfig) (string, error)
}

func GetAdapter(provider CloudProvider) (CloudAdapter, error) {
	switch provider {
	case AWS:
		return &AWSAdapter{}, nil
	case GCP:
		return &GCPAdapter{}, nil
	case Azure:
		return &AzureAdapter{}, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
