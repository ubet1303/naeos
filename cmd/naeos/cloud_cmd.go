package main

import (
	"fmt"
	"os"

	"github.com/NAEOS-foundation/naeos/internal/cloud"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cloudProvider string
	cloudRegion   string
	cloudProject  string
	cloudEnv      string
	cloudInput    string
)

type cloudSpec struct {
	Cloud struct {
		Provider    string `yaml:"provider"`
		Region      string `yaml:"region"`
		Project     string `yaml:"project"`
		Environment string `yaml:"environment"`
		Resources   []struct {
			Name string            `yaml:"name"`
			Kind string            `yaml:"kind"`
			Type string            `yaml:"type"`
			Spec map[string]string `yaml:"spec"`
		} `yaml:"resources"`
	} `yaml:"cloud"`
}

func loadCloudConfigFromSpec(path string) (*cloud.DeployConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spec file: %w", err)
	}

	var spec cloudSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parse spec file: %w", err)
	}

	if spec.Cloud.Provider == "" {
		return nil, fmt.Errorf("cloud.provider is required in spec")
	}

	config := &cloud.DeployConfig{
		Provider:    cloud.CloudProvider(spec.Cloud.Provider),
		Region:      spec.Cloud.Region,
		Project:     spec.Cloud.Project,
		Environment: spec.Cloud.Environment,
	}

	for _, r := range spec.Cloud.Resources {
		resType := r.Type
		if resType == "" {
			resType = r.Kind
		}
		specMap := make(map[string]interface{})
		for k, v := range r.Spec {
			specMap[k] = v
		}
		config.Resources = append(config.Resources, cloud.Resource{
			Name: r.Name,
			Type: resType,
			Spec: specMap,
		})
	}

	return config, nil
}

func resolveCloudConfig() (*cloud.DeployConfig, error) {
	if cloudInput != "" {
		return loadCloudConfigFromSpec(cloudInput)
	}

	config := &cloud.DeployConfig{
		Provider:    cloud.CloudProvider(cloudProvider),
		Region:      cloudRegion,
		Project:     cloudProject,
		Environment: cloudEnv,
	}

	return config, nil
}

func newCloudCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Cloud deployment commands",
		Long:  `Deploy NAEOS projects to AWS, GCP, or Azure.`,
	}

	cmd.AddCommand(newCloudDeployCommand())
	cmd.AddCommand(newCloudPlanCommand())
	cmd.AddCommand(newCloudExportCommand())
	cmd.AddCommand(newCloudTypesCommand())

	cmd.PersistentFlags().StringVarP(&cloudProvider, "provider", "p", "aws", "Cloud provider (aws, gcp, azure)")
	cmd.PersistentFlags().StringVarP(&cloudRegion, "region", "r", "", "Cloud region")
	cmd.PersistentFlags().StringVarP(&cloudProject, "project", "j", "", "Cloud project name")
	cmd.PersistentFlags().StringVarP(&cloudEnv, "env", "e", "dev", "Environment (dev, staging, prod)")
	cmd.PersistentFlags().StringVarP(&cloudInput, "input-file", "i", "", "Spec file with cloud configuration (overrides flags)")

	return cmd
}

func newCloudDeployCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy to cloud provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := resolveCloudConfig()
			if err != nil {
				return err
			}

			adapter, err := cloud.GetAdapter(config.Provider)
			if err != nil {
				return err
			}

			if err := adapter.Validate(config); err != nil {
				return err
			}

			result, err := adapter.Deploy(config)
			if err != nil {
				return err
			}

			fmt.Printf("Deployed to %s: %d resources\n", result.Provider, len(result.Resources))
			for _, r := range result.Resources {
				fmt.Printf("  - %s (%s) -> %s\n", r.Name, r.Type, r.ID)
			}
			return nil
		},
	}
}

func newCloudPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Plan cloud deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := resolveCloudConfig()
			if err != nil {
				return err
			}

			adapter, err := cloud.GetAdapter(config.Provider)
			if err != nil {
				return err
			}

			if err := adapter.Validate(config); err != nil {
				return err
			}

			plan, err := adapter.Plan(config)
			if err != nil {
				return err
			}

			fmt.Printf("Plan: %d resources to deploy (%s/%s)\n", len(plan), config.Provider, config.Region)
			for _, res := range plan {
				fmt.Printf("  - %s (%s)\n", res.Name, res.Type)
			}
			return nil
		},
	}
}

func newCloudExportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export Terraform configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := resolveCloudConfig()
			if err != nil {
				return err
			}

			adapter, err := cloud.GetAdapter(config.Provider)
			if err != nil {
				return err
			}

			if err := adapter.Validate(config); err != nil {
				return err
			}

			tf, err := adapter.ExportTerraform(config)
			if err != nil {
				return err
			}

			fmt.Println(tf)
			return nil
		},
	}
}

func newCloudTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List supported resource types",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Supported resource types:")
			for _, t := range cloud.SupportedResourceTypes {
				fmt.Printf("  - %s\n", t)
			}
			return nil
		},
	}
}
