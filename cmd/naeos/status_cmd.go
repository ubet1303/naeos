package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cfgpkg "github.com/NAEOS-foundation/naeos/pkg/config"
)

func newStatusCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current pipeline and project status",
		Long: `Display the current status of the NAEOS project and pipeline configuration.

Example:
  naeos status
  naeos status --config config.yaml
  naeos status -o json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := resolveConfigPath(configPath)
			if err != nil {
				return err
			}

			fileCfg, err := cfgpkg.LoadFile(resolved)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			type statusOutput struct {
				Config    string   `json:"config" yaml:"config"`
				Pipeline  string   `json:"pipeline" yaml:"pipeline"`
				Mode      string   `json:"mode" yaml:"mode"`
				OutputDir string   `json:"output_dir" yaml:"output_dir"`
				Languages []string `json:"languages" yaml:"languages"`
				Verbose   bool     `json:"verbose" yaml:"verbose"`
				CheckedAt string   `json:"checked_at" yaml:"checked_at"`
			}

			status := statusOutput{
				Config:    resolved,
				Pipeline:  fileCfg.Pipeline.Name,
				Mode:      fileCfg.Pipeline.Mode,
				OutputDir: fileCfg.Pipeline.OutputDir,
				Languages: fileCfg.Pipeline.Language,
				Verbose:   fileCfg.Pipeline.Verbose,
				CheckedAt: time.Now().Format(time.RFC3339),
			}

			switch cliOutputFormat {
			case "json", "yaml":
				return FormatOutput(cmd.OutOrStdout(), status, cliOutputFormat)
			default:
				out := cmd.OutOrStdout()
				fmt.Fprintf(out, "NAEOS Status\n")
				fmt.Fprintf(out, "%s\n", "================================")
				fmt.Fprintf(out, "Config:        %s\n", status.Config)
				fmt.Fprintf(out, "Pipeline:      %s\n", status.Pipeline)
				fmt.Fprintf(out, "Mode:          %s\n", status.Mode)
				fmt.Fprintf(out, "Output Dir:    %s\n", status.OutputDir)
				fmt.Fprintf(out, "Languages:     %s\n", joinStrings(status.Languages))
				fmt.Fprintf(out, "Verbose:       %t\n", status.Verbose)
				fmt.Fprintf(out, "Checked At:    %s\n", status.CheckedAt)
				return nil
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to JSON or YAML config file (auto-detected if omitted)")
	return cmd
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return "(none)"
	}
	return strings.Join(ss, ", ")
}
