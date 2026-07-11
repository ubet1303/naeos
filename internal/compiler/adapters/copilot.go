package adapters

import (
	"fmt"
	"strings"
	"time"

	"github.com/NAEOS-foundation/naeos/internal/compiler"
	"github.com/NAEOS-foundation/naeos/internal/neir/model"
)

type copilotAdapter struct{}

func NewCopilotAdapter() compiler.Adapter {
	return &copilotAdapter{}
}

func (a *copilotAdapter) Target() compiler.Target {
	return compiler.TargetCopilot
}

func (a *copilotAdapter) Compile(neir *model.NEIR) (*compiler.CompiledOutput, error) {
	if neir == nil {
		return nil, fmt.Errorf("nil NEIR")
	}

	var files []compiler.OutputFile

	instructions := a.buildInstructions(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".github/copilot-instructions.md",
		Content: instructions,
		Kind:    "instructions",
	})

	context := a.buildContextFile(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".github/copilot-context.md",
		Content: context,
		Kind:    "context",
	})

	rules := a.buildRulesFile(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".github/copilot-rules.md",
		Content: rules,
		Kind:    "rules",
	})

	projectName := "unknown"
	if neir.Project != nil {
		projectName = neir.Project.Name
	}

	return &compiler.CompiledOutput{
		Target:  compiler.TargetCopilot,
		Files:   files,
		Summary: fmt.Sprintf("Generated %d files for GitHub Copilot (%s)", len(files), projectName),
		CompiledAt: time.Now(),
	}, nil
}

func (a *copilotAdapter) buildInstructions(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# GitHub Copilot Instructions\n\n")
	sb.WriteString("This file contains project-specific instructions for GitHub Copilot.\n\n")

	if neir.Project != nil {
		sb.WriteString(fmt.Sprintf("## Project: %s\n\n", neir.Project.Name))
		if neir.Project.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", neir.Project.Description))
		}
	}

	if neir.Architecture != nil {
		sb.WriteString(fmt.Sprintf("## Architecture Pattern: %s\n\n", neir.Architecture.Pattern))
		sb.WriteString("Follow these architectural principles:\n")
		for _, p := range neir.Architecture.Principles {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\n")
	}

	if len(neir.Modules) > 0 {
		sb.WriteString("## Module Structure\n\n")
		for _, m := range neir.Modules {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", m.Name, m.Path))
			if m.Description != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", m.Description))
			}
		}
		sb.WriteString("\n")
	}

	if len(neir.Components) > 0 {
		sb.WriteString("## Components\n\n")
		for _, c := range neir.Components {
			sb.WriteString(fmt.Sprintf("- `%s` (%s) in module `%s`\n", c.Name, c.Kind, c.Module))
		}
		sb.WriteString("\n")
	}

	if len(neir.Services) > 0 {
		sb.WriteString("## Services\n\n")
		for _, s := range neir.Services {
			sb.WriteString(fmt.Sprintf("### %s (%s, port %d)\n", s.Name, s.Kind, s.Port))
			for _, ep := range s.Endpoints {
				sb.WriteString(fmt.Sprintf("- `%s %s` -> `%s`\n", ep.Method, ep.Path, ep.Action))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("## Coding Guidelines\n\n")
	sb.WriteString("- Write clean, idiomatic code following the project's architecture pattern\n")
	sb.WriteString("- Include proper error handling\n")
	sb.WriteString("- Add comments for public APIs\n")
	sb.WriteString("- Follow the module dependency structure\n")

	return sb.String()
}

func (a *copilotAdapter) buildContextFile(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Project Context for Copilot\n\n")
	sb.WriteString("Use this file as additional context when generating code.\n\n")
	sb.WriteString("```yaml\n")
	sb.WriteString(fmt.Sprintf("project: %s\n", neir.Project.Name))
	if neir.Architecture != nil {
		sb.WriteString(fmt.Sprintf("architecture: %s\n", neir.Architecture.Pattern))
	}
	if len(neir.Modules) > 0 {
		sb.WriteString("modules:\n")
		for _, m := range neir.Modules {
			sb.WriteString(fmt.Sprintf("  - name: %s\n    path: %s\n", m.Name, m.Path))
		}
	}
	sb.WriteString("```\n")
	return sb.String()
}

func (a *copilotAdapter) buildRulesFile(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Copilot Rules\n\n")
	sb.WriteString("## File Organization\n\n")

	if len(neir.Modules) > 0 {
		for _, m := range neir.Modules {
			sb.WriteString(fmt.Sprintf("- Files in `%s` belong to the `%s` module\n", m.Path, m.Name))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Code Style\n\n")
	sb.WriteString("- Use early returns for error handling\n")
	sb.WriteString("- Prefer composition over inheritance\n")
	sb.WriteString("- Keep functions small and focused\n")
	sb.WriteString("- Use meaningful variable and function names\n")

	return sb.String()
}
