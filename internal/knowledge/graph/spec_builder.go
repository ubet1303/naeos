package graph

import (
	"fmt"
	"strings"

	"github.com/NAEOS-foundation/naeos/internal/specification/parser"
)

// BuildFromSpec creates a knowledge graph from a NAEOS specification.
func BuildFromSpec(specContent string) (*KnowledgeGraph, error) {
	p := parser.NewParser("")
	doc, err := p.Parse(specContent)
	if err != nil {
		return nil, fmt.Errorf("parse spec: %w", err)
	}
	if doc == nil {
		return New(), nil
	}

	kg := New()

	if doc.Project != "" {
		_ = kg.AddNode(Node{
			ID:        "project",
			Type:      NodeTypeComponent,
			Topic:     doc.Project,
			Component: doc.Project,
			Content:   specContent,
		})
	}

	for i, m := range doc.Modules {
		id := fmt.Sprintf("module-%s", m.Name)
		if m.Name == "" {
			id = fmt.Sprintf("module-%d", i)
		}
		_ = kg.AddNode(Node{
			ID:        id,
			Type:      NodeTypeModule,
			Topic:     m.Name,
			Component: m.Name,
		})
		if doc.Project != "" {
			_ = kg.AddEdge(Edge{From: "project", To: id, Type: EdgeTypeContains})
		}
		for _, dep := range m.Dependencies {
			depID := fmt.Sprintf("module-%s", dep)
			_ = kg.AddEdge(Edge{From: id, To: depID, Type: EdgeTypeDependsOn})
		}
	}

	for i, svc := range doc.Services {
		id := fmt.Sprintf("service-%s", svc.Name)
		if svc.Name == "" {
			id = fmt.Sprintf("service-%d", i)
		}
		_ = kg.AddNode(Node{
			ID:        id,
			Type:      NodeTypeService,
			Topic:     svc.Name,
			Component: svc.Name,
		})
		if doc.Project != "" {
			_ = kg.AddEdge(Edge{From: "project", To: id, Type: EdgeTypeContains})
		}
		for _, ep := range svc.Endpoints {
			epID := fmt.Sprintf("api-%s-%s", svc.Name, strings.ToLower(ep.Method))
			_ = kg.AddNode(Node{
				ID:        epID,
				Type:      NodeTypeAPI,
				Topic:     ep.Action,
				Component: svc.Name,
			})
			_ = kg.AddEdge(Edge{From: id, To: epID, Type: EdgeTypeExposes})
		}
	}

	if doc.Architecture != nil && doc.Architecture.Pattern != "" {
		_ = kg.AddNode(Node{
			ID:      "architecture",
			Type:    NodeTypeDecision,
			Topic:   doc.Architecture.Pattern,
			Content: doc.Architecture.Description,
		})
		if doc.Project != "" {
			_ = kg.AddEdge(Edge{From: "project", To: "architecture", Type: EdgeTypeImplements})
		}
	}

	if doc.Deployment != nil && doc.Deployment.Strategy != "" {
		_ = kg.AddNode(Node{
			ID:        "deployment",
			Type:      NodeTypeDeployment,
			Topic:     doc.Deployment.Strategy,
			Component: doc.Project,
		})
		if doc.Project != "" {
			_ = kg.AddEdge(Edge{From: "project", To: "deployment", Type: EdgeTypeDeploysTo})
		}
	}

	if doc.Testing != nil {
		_ = kg.AddNode(Node{
			ID:        "testing",
			Type:      NodeTypeTesting,
			Topic:     doc.Testing.Strategy,
			Component: doc.Project,
		})
		if doc.Project != "" {
			_ = kg.AddEdge(Edge{From: "project", To: "testing", Type: EdgeTypeTests})
		}
	}

	return kg, nil
}

// AnalyzeSpecGraph analyzes a knowledge graph and returns suggestions.
func AnalyzeSpecGraph(kg *KnowledgeGraph) []Suggestion {
	var suggestions []Suggestion

	services := kg.FindByType(NodeTypeService)
	modules := kg.FindByType(NodeTypeModule)

	if len(services) > 3 && len(modules) == 0 {
		suggestions = append(suggestions, Suggestion{
			Category: "structure",
			Title:    "Consider adding modules",
			Description: fmt.Sprintf(
				"Found %d services but no modules. Consider organizing services into logical modules.",
				len(services)),
			Priority: "medium",
		})
	}

	if len(modules) > 0 {
		cycles := kg.DetectCycles()
		if len(cycles) > 0 {
			suggestions = append(suggestions, Suggestion{
				Category:    "architecture",
				Title:       "Circular dependency detected",
				Description: "The module dependency graph has cycles. This may cause build or initialization issues.",
				Priority:    "high",
			})
		}
	}

	apis := kg.FindByType(NodeTypeAPI)
	if len(apis) == 0 && len(services) > 0 {
		suggestions = append(suggestions, Suggestion{
			Category:    "api",
			Title:       "No API endpoints defined",
			Description: "Services are defined but no API endpoints. Add endpoints to expose functionality.",
			Priority:    "medium",
		})
	}

	if kg.FindByType(NodeTypeSecurity) == nil {
		suggestions = append(suggestions, Suggestion{
			Category:    "security",
			Title:       "No security configuration",
			Description: "Consider adding security policies to your specification.",
			Priority:    "high",
		})
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, Suggestion{
			Category:    "general",
			Title:       "Graph analysis passed",
			Description: "No issues found in the knowledge graph analysis.",
			Priority:    "low",
		})
	}

	return suggestions
}

// Suggestion is a local copy to avoid circular imports.
type Suggestion struct {
	Category    string
	Title       string
	Description string
	Priority    string
}
