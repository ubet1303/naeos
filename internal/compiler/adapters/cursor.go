package adapters

import (
	"fmt"
	"strings"
	"time"

	"github.com/NAEOS-foundation/naeos/internal/compiler"
	"github.com/NAEOS-foundation/naeos/internal/neir/model"
)

type cursorAdapter struct{}

func NewCursorAdapter() compiler.Adapter {
	return &cursorAdapter{}
}

func (a *cursorAdapter) Target() compiler.Target {
	return compiler.TargetCursor
}

func (a *cursorAdapter) Compile(neir *model.NEIR) (*compiler.CompiledOutput, error) {
	if neir == nil {
		return nil, fmt.Errorf("nil NEIR")
	}

	var files []compiler.OutputFile

	rulesFile := a.buildRulesFile(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".cursorrules",
		Content: rulesFile,
		Kind:    "rules",
	})

	contextFile := a.buildContextFile(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".cursor/context.md",
		Content: contextFile,
		Kind:    "context",
	})

	projectName := "unknown"
	if neir.Project != nil {
		projectName = neir.Project.Name
	}

	return &compiler.CompiledOutput{
		Target:     compiler.TargetCursor,
		Files:      files,
		Summary:    fmt.Sprintf("Generated %d files for Cursor (%s)", len(files), projectName),
		CompiledAt: time.Now(),
	}, nil
}

func (a *cursorAdapter) buildRulesFile(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Cursor Rules\n\n")

	if neir.Project != nil {
		sb.WriteString(fmt.Sprintf("project_name: %s\n", neir.Project.Name))
	}
	if neir.Architecture != nil {
		sb.WriteString(fmt.Sprintf("architecture: %s\n", neir.Architecture.Pattern))
	}

	sb.WriteString("\n## Instructions\n\n")
	sb.WriteString("You are working on a project with the following structure:\n\n")

	if len(neir.Modules) > 0 {
		sb.WriteString("### Modules\n\n")
		for _, m := range neir.Modules {
			sb.WriteString(fmt.Sprintf("- `%s` at `%s`\n", m.Name, m.Path))
			if m.Description != "" {
				sb.WriteString(fmt.Sprintf("  > %s\n", m.Description))
			}
		}
		sb.WriteString("\n")
	}

	if len(neir.Services) > 0 {
		sb.WriteString("### Services\n\n")
		for _, s := range neir.Services {
			sb.WriteString(fmt.Sprintf("- **%s** (%s, port %d)\n", s.Name, s.Kind, s.Port))
			for _, ep := range s.Endpoints {
				sb.WriteString(fmt.Sprintf("  - %s %s -> %s\n", ep.Method, ep.Path, ep.Action))
			}
		}
		sb.WriteString("\n")
	}

	if len(neir.Components) > 0 {
		sb.WriteString("### Components\n\n")
		for _, c := range neir.Components {
			sb.WriteString(fmt.Sprintf("- `%s` [%s] module: `%s`\n", c.Name, c.Kind, c.Module))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Style Rules\n\n")
	sb.WriteString("- Use early returns\n")
	sb.WriteString("- Prefer const/readonly where possible\n")
	sb.WriteString("- Write self-documenting code with clear names\n")
	sb.WriteString("- Handle all error paths\n")
	sb.WriteString("- Keep functions under 50 lines\n")

	return sb.String()
}

func (a *cursorAdapter) buildContextFile(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Cursor Context\n\n")
	sb.WriteString("Additional project context for AI-assisted coding.\n\n")

	if neir.Project != nil {
		sb.WriteString(fmt.Sprintf("Project: %s v%s\n", neir.Project.Name, neir.Project.Version))
	}

	if len(neir.Modules) > 0 {
		sb.WriteString("\n## Module Dependency Map\n\n")
		for _, m := range neir.Modules {
			if len(m.Dependencies) > 0 {
				sb.WriteString(fmt.Sprintf("%s depends on: %s\n", m.Name, strings.Join(m.Dependencies, ", ")))
			}
		}
	}

	if len(neir.APIs) > 0 {
		sb.WriteString("\n## API Endpoints\n\n")
		for _, api := range neir.APIs {
			sb.WriteString(fmt.Sprintf("### %s v%s (%s)\n", api.Name, api.Version, api.Protocol))
			for _, ep := range api.Endpoints {
				sb.WriteString(fmt.Sprintf("- %s %s: %s\n", ep.Method, ep.Path, ep.Summary))
			}
		}
	}

	return sb.String()
}
