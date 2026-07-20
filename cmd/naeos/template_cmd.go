package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/promptlib"
	"github.com/NAEOS-foundation/naeos/internal/templates"
)

func newTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage generation templates and prompt library",
		Long: `Manage NAEOS generation templates and prompt library.

Example:
  naeos template list
  naeos template list --kind code
  naeos template list --kind prompt-llm
  naeos template list --kind prompt-compiler
  naeos template show enrich-spec
  naeos template show claude
  naeos template add my-template "template content"
  naeos template remove my-template`,
	}

	var templatesDir string

	templateList := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			kindFilter, _ := cmd.Flags().GetString("kind")

			if kindFilter == "" || kindFilter == "code" {
				mgr := templates.NewManager(templatesDir)
				tmpls, err := mgr.List()
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Code Generation Templates:")
				for _, t := range tmpls {
					custom := ""
					if t.IsCustom {
						custom = " (custom)"
					}
					fmt.Fprintf(cmd.OutOrStdout(), "  %-20s %s%s\n", t.Name, t.Path, custom)
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}

			if kindFilter == "" || kindFilter == "prompt-llm" || kindFilter == "prompt-compiler" {
				lib := promptlib.NewWithDefaults()

				if kindFilter == "" || kindFilter == "prompt-llm" {
					llmPrompts := lib.ListLLMPrompts()
					fmt.Fprintln(cmd.OutOrStdout(), "LLM Prompt Templates:")
					for _, name := range llmPrompts {
						p, _ := lib.GetLLMPrompt(name)
						desc := ""
						if p != nil && p.Description != "" {
							desc = " - " + p.Description
						}
						fmt.Fprintf(cmd.OutOrStdout(), "  %-30s%s\n", name, desc)
					}
					fmt.Fprintln(cmd.OutOrStdout())
				}

				if kindFilter == "" || kindFilter == "prompt-compiler" {
					compilerTpls := lib.ListCompilerTemplates()
					fmt.Fprintln(cmd.OutOrStdout(), "Compiler Templates:")
					for _, name := range compilerTpls {
						t, _ := lib.GetCompilerTemplate(name)
						target := ""
						if t != nil {
							target = " (target: " + t.Target + ")"
						}
						fmt.Fprintf(cmd.OutOrStdout(), "  %-20s%s\n", name, target)
					}
				}
			}

			return nil
		},
	}
	templateList.Flags().String("kind", "", "filter by kind: code, prompt-llm, prompt-compiler")

	templateShow := &cobra.Command{
		Use:   "show [name]",
		Short: "Show details of a prompt template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			lib := promptlib.NewWithDefaults()

			if p, ok := lib.GetLLMPrompt(name); ok {
				fmt.Fprintf(cmd.OutOrStdout(), "Name:        %s\n", p.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "Kind:        llm\n")
				fmt.Fprintf(cmd.OutOrStdout(), "Version:     %s\n", p.Version)
				fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", p.Description)
				fmt.Fprintf(cmd.OutOrStdout(), "Provider:    %s\n", p.Provider)
				if p.Constraints != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "MaxTokens:   %d\n", p.Constraints.MaxTokens)
					fmt.Fprintf(cmd.OutOrStdout(), "Temperature: %.1f\n", p.Constraints.Temperature)
				}
				if p.System != "" {
					fmt.Fprintln(cmd.OutOrStdout(), "\nSystem Prompt:")
					fmt.Fprintln(cmd.OutOrStdout(), indent(p.System, "  "))
				}
				fmt.Fprintln(cmd.OutOrStdout(), "\nUser Prompt:")
				fmt.Fprintln(cmd.OutOrStdout(), indent(p.User, "  "))
				if len(p.Variables) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "\nVariables:")
					for _, v := range p.Variables {
						req := ""
						if v.Required {
							req = " (required)"
						}
						fmt.Fprintf(cmd.OutOrStdout(), "  - %s [%s]%s: %s\n", v.Name, v.Type, req, v.Description)
					}
				}
				return nil
			}

			if t, ok := lib.GetCompilerTemplate(name); ok {
				fmt.Fprintf(cmd.OutOrStdout(), "Name:    %s\n", t.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "Kind:    compiler\n")
				fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", t.Version)
				fmt.Fprintf(cmd.OutOrStdout(), "Target:  %s\n", t.Target)
				fmt.Fprintf(cmd.OutOrStdout(), "\nOutput Files:\n")
				for _, f := range t.Files {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s)\n", f.Path, f.Kind)
				}
				if len(t.Variables) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "\nVariables:")
					for _, v := range t.Variables {
						fmt.Fprintf(cmd.OutOrStdout(), "  - %s [%s]\n", v.Name, v.Type)
					}
				}
				return nil
			}

			return fmt.Errorf("template %q not found (searched LLM prompts and compiler templates)", name)
		},
	}

	templateAdd := &cobra.Command{
		Use:   "add [name] [content]",
		Short: "Add a custom template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := templates.NewManager(templatesDir)
			if err := mgr.AddCustom(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added template %s\n", args[0])
			return nil
		},
	}

	templateRemove := &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a custom template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := templates.NewManager(templatesDir)
			if err := mgr.RemoveCustom(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed template %s\n", args[0])
			return nil
		},
	}

	promptCreate := &cobra.Command{
		Use:   "prompt-create [name]",
		Short: "Create a custom LLM prompt template",
		Long: `Create a custom LLM prompt template that can be used with 'naeos ai enrich'.

The command opens an interactive editor or accepts --system and --user flags.
Example:
  naeos template prompt-create my-custom-prompt --system "You are an expert" --user "Analyze this: {{.SpecContent}}"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			system, _ := cmd.Flags().GetString("system")
			user, _ := cmd.Flags().GetString("user")
			provider, _ := cmd.Flags().GetString("provider")
			desc, _ := cmd.Flags().GetString("description")

			if user == "" {
				return fmt.Errorf("--user prompt content is required")
			}
			if system == "" {
				system = "You are a helpful assistant."
			}
			if provider == "" {
				provider = "openai"
			}

			promptDir := filepath.Join(templatesDir, "prompts")
			if err := os.MkdirAll(promptDir, 0o755); err != nil {
				return err
			}

			content := fmt.Sprintf(`kind: llm
name: %s
version: "1.0.0"
description: %s
provider: %s
system: |
  %s
user: |
  %s
variables:
  - name: SpecContent
    type: string
    description: Specification content
    required: true
`, name, desc, provider, strings.ReplaceAll(system, "\n", "\n  "), strings.ReplaceAll(user, "\n", "\n  "))

			path := filepath.Join(promptDir, name+".yaml")
			if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
				return fmt.Errorf("write prompt: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created prompt template: %s\n", path)
			return nil
		},
	}
	promptCreate.Flags().String("system", "", "system prompt content")
	promptCreate.Flags().String("user", "", "user prompt content (required)")
	promptCreate.Flags().String("provider", "openai", "LLM provider (openai, anthropic, ollama)")
	promptCreate.Flags().String("description", "", "description of the prompt")

	promptRemove := &cobra.Command{
		Use:   "prompt-remove [name]",
		Short: "Remove a custom prompt template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			promptDir := filepath.Join(templatesDir, "prompts")
			path := filepath.Join(promptDir, name+".yaml")
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("prompt %q not found at %s", name, path)
			}
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("remove prompt: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed prompt template: %s\n", name)
			return nil
		},
	}

	cmd.AddCommand(templateList)
	cmd.AddCommand(templateShow)
	cmd.AddCommand(templateAdd)
	cmd.AddCommand(templateRemove)
	cmd.AddCommand(promptCreate)
	cmd.AddCommand(promptRemove)
	cmd.PersistentFlags().StringVar(&templatesDir, "templates-dir", filepath.Join(".", ".naeos", "templates"), "templates directory")
	return cmd
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}
