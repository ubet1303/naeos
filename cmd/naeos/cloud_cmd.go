package main

import (
	"fmt"
	"github.com/NAEOS-foundation/naeos/internal/cloud"
	"github.com/spf13/cobra"
)

var (
	cloudProvider string
	cloudRegion   string
	cloudProject  string
	cloudEnv      string
)

func newCloudCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Cloud deployment commands",
		Long:  `Deploy NAEOS projects to AWS, GCP, or Azure.`,
	}

	cmd.AddCommand(newCloudDeployCommand())
	cmd.AddCommand(newCloudPlanCommand())
	cmd.AddCommand(newCloudExportCommand())

	cmd.PersistentFlags().StringVarP(&cloudProvider, "provider", "p", "aws", "Cloud provider (aws, gcp, azure)")
	cmd.PersistentFlags().StringVarP(&cloudRegion, "region", "r", "", "Cloud region")
	cmd.PersistentFlags().StringVarP(&cloudProject, "project", "j", "", "Cloud project name")
	cmd.PersistentFlags().StringVarP(&cloudEnv, "env", "e", "dev", "Environment (dev, staging, prod)")

	return cmd
}

func newCloudDeployCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy to cloud provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := cloud.CloudProvider(cloudProvider)
			adapter, err := cloud.GetAdapter(provider)
			if err != nil {
				return err
			}

			config := &cloud.DeployConfig{
				Provider:    provider,
				Region:      cloudRegion,
				Project:     cloudProject,
				Environment: cloudEnv,
			}

			if err := adapter.Validate(config); err != nil {
				return err
			}

			result, err := adapter.Deploy(config)
			if err != nil {
				return err
			}

			fmt.Printf("Deployed to %s: %d resources\n", result.Provider, len(result.Resources))
			return nil
		},
	}
}

func newCloudPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Plan cloud deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := cloud.CloudProvider(cloudProvider)
			adapter, err := cloud.GetAdapter(provider)
			if err != nil {
				return err
			}

			config := &cloud.DeployConfig{
				Provider:    provider,
				Region:      cloudRegion,
				Project:     cloudProject,
				Environment: cloudEnv,
			}

			plan, err := adapter.Plan(config)
			if err != nil {
				return err
			}

			fmt.Printf("Plan: %d resources to deploy\n", len(plan))
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
			provider := cloud.CloudProvider(cloudProvider)
			adapter, err := cloud.GetAdapter(provider)
			if err != nil {
				return err
			}

			config := &cloud.DeployConfig{
				Provider:    provider,
				Region:      cloudRegion,
				Project:     cloudProject,
				Environment: cloudEnv,
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
