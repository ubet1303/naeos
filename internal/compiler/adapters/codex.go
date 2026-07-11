package adapters

import (
	"fmt"
	"strings"
	"time"

	"github.com/NAEOS-foundation/naeos/internal/compiler"
	"github.com/NAEOS-foundation/naeos/internal/neir/model"
)

type codexAdapter struct{}

func NewCodexAdapter() compiler.Adapter {
	return &codexAdapter{}
}

func (a *codexAdapter) Target() compiler.Target {
	return compiler.TargetCodex
}

func (a *codexAdapter) Compile(neir *model.NEIR) (*compiler.CompiledOutput, error) {
	if neir == nil {
		return nil, fmt.Errorf("nil NEIR")
	}

	var files []compiler.OutputFile

	instructions := a.buildInstructions(neir)
	files = append(files, compiler.OutputFile{
		Path:    "AGENTS.md",
		Content: instructions,
		Kind:    "instructions",
	})

	contextFile := a.buildContextFile(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".codex/context.md",
		Content: contextFile,
		Kind:    "context",
	})

	projectName := "unknown"
	if neir.Project != nil {
		projectName = neir.Project.Name
	}

	return &compiler.CompiledOutput{
		Target:     compiler.TargetCodex,
		Files:      files,
		Summary:    fmt.Sprintf("Generated %d files for Codex (%s)", len(files), projectName),
		CompiledAt: time.Now(),
	}, nil
}

func (a *codexAdapter) buildInstructions(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# AGENTS.md\n\n")
	sb.WriteString("Instructions for AI agents working on this project.\n\n")

	if neir.Project != nil {
		sb.WriteString(fmt.Sprintf("## Project: %s\n\n", neir.Project.Name))
		if neir.Project.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", neir.Project.Description))
		}
	}

	if neir.Architecture != nil {
		sb.WriteString(fmt.Sprintf("## Architecture: %s\n\n", neir.Architecture.Pattern))
		for _, p := range neir.Architecture.Principles {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Module Structure\n\n")
	if len(neir.Modules) > 0 {
		for _, m := range neir.Modules {
			sb.WriteString(fmt.Sprintf("### %s\nPath: `%s`\n", m.Name, m.Path))
			if m.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n", m.Description))
			}
			if len(m.Dependencies) > 0 {
				sb.WriteString(fmt.Sprintf("Dependencies: %s\n", strings.Join(m.Dependencies, ", ")))
			}
			sb.WriteString("\n")
		}
	}

	if len(neir.Services) > 0 {
		sb.WriteString("## Services\n\n")
		for _, s := range neir.Services {
			sb.WriteString(fmt.Sprintf("### %s (%s, port %d)\n", s.Name, s.Kind, s.Port))
			for _, ep := range s.Endpoints {
				sb.WriteString(fmt.Sprintf("- %s %s → %s\n", ep.Method, ep.Path, ep.Action))
			}
			sb.WriteString("\n")
		}
	}

	if len(neir.Components) > 0 {
		sb.WriteString("## Components\n\n")
		for _, c := range neir.Components {
			sb.WriteString(fmt.Sprintf("- `%s` [%s] in `%s`\n", c.Name, c.Kind, c.Module))
		}
		sb.WriteString("\n")
	}

	if neir.Deployment != nil {
		sb.WriteString(fmt.Sprintf("## Deployment: %s\n\n", neir.Deployment.Strategy))
		if len(neir.Deployment.Environments) > 0 {
			for _, env := range neir.Deployment.Environments {
				sb.WriteString(fmt.Sprintf("- %s (%s)\n", env.Name, env.Kind))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Agent Guidelines\n\n")
	sb.WriteString("1. Follow the architecture pattern\n")
	sb.WriteString("2. Write clean, idiomatic code for the target language\n")
	sb.WriteString("3. Handle errors explicitly\n")
	sb.WriteString("4. Write tests for new code\n")
	sb.WriteString("5. Keep functions focused and small\n")
	sb.WriteString("6. Document public APIs\n")

	return sb.String()
}

func (a *codexAdapter) buildContextFile(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Codex Context\n\n")

	if len(neir.Storage) > 0 {
		sb.WriteString("## Storage\n\n")
		for _, st := range neir.Storage {
			sb.WriteString(fmt.Sprintf("- %s (%s) via %s\n", st.Name, st.Type, st.Provider))
			for _, col := range st.Collections {
				sb.WriteString(fmt.Sprintf("  - %s\n", col.Name))
			}
		}
		sb.WriteString("\n")
	}

	if neir.Infrastructure != nil {
		sb.WriteString(fmt.Sprintf("## Infrastructure: %s\n\n", neir.Infrastructure.Provider))
		for _, r := range neir.Infrastructure.Resources {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", r.Name, r.Kind))
		}
	}

	if neir.AI != nil && len(neir.AI.Models) > 0 {
		sb.WriteString("\n## AI Models\n\n")
		for _, m := range neir.AI.Models {
			sb.WriteString(fmt.Sprintf("- %s (%s) v%s\n", m.Name, m.Kind, m.Version))
		}
	}

	return sb.String()
}
