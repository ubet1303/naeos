package create

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Wizard struct {
	reader *bufio.Reader
}

type ProjectConfig struct {
	Name             string
	ModulePath       string
	Language         string
	Architecture     string
	Deployment       string
	Port             int
	OutputDir        string
	Description      string
	EnableAuth       bool
	EnableTesting    bool
	EnableDocker     bool
	EnableCI         bool
}

func NewWizard() *Wizard {
	return &Wizard{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (w *Wizard) Run() (*ProjectConfig, error) {
	cfg := &ProjectConfig{}

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║     NAEOS Project Creation Wizard    ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	cfg.Name = w.askRequired("Project name")
	cfg.ModulePath = w.askDefault("Module path", "./"+strings.ToLower(strings.ReplaceAll(cfg.Name, " ", "-")))
	cfg.Description = w.askDefault("Description", "A NAEOS project")
	cfg.Language = w.askChoice("Language", []string{"go", "typescript", "python", "java", "rust"}, "go")
	cfg.Architecture = w.askChoice("Architecture pattern", []string{"hexagonal", "layered", "clean", "event-driven", "cqrs", "monolith"}, "hexagonal")
	cfg.Deployment = w.askChoice("Deployment strategy", []string{"rolling", "blue-green", "canary", "recreate"}, "rolling")
	cfg.Port = w.askInt("Default port", 8080)
	cfg.OutputDir = w.askDefault("Output directory", cfg.Name)
	cfg.EnableAuth = w.askYesNo("Enable authentication", false)
	cfg.EnableTesting = w.askYesNo("Enable test generation", true)
	cfg.EnableDocker = w.askYesNo("Generate Dockerfile", true)
	cfg.EnableCI = w.askYesNo("Generate CI workflow", true)

	fmt.Println()
	fmt.Println("Configuration complete!")
	return cfg, nil
}

func (w *Wizard) askRequired(prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		text, _ := w.reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text != "" {
			return text
		}
		fmt.Println("  This field is required.")
	}
}

func (w *Wizard) askDefault(prompt, defaultVal string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultVal)
	text, _ := w.reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultVal
	}
	return text
}

func (w *Wizard) askChoice(prompt string, options []string, defaultVal string) string {
	fmt.Printf("%s:\n", prompt)
	for i, opt := range options {
		marker := "  "
		if opt == defaultVal {
			marker = "→ "
		}
		fmt.Printf("  %s%d) %s\n", marker, i+1, opt)
	}
	fmt.Printf("  Choose [1-%d] (default: %s): ", len(options), defaultVal)
	text, _ := w.reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultVal
	}
	var idx int
	if _, err := fmt.Sscanf(text, "%d", &idx); err == nil && idx >= 1 && idx <= len(options) {
		return options[idx-1]
	}
	return defaultVal
}

func (w *Wizard) askInt(prompt string, defaultVal int) int {
	fmt.Printf("%s [%d]: ", prompt, defaultVal)
	text, _ := w.reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultVal
	}
	var val int
	if _, err := fmt.Sscanf(text, "%d", &val); err == nil {
		return val
	}
	return defaultVal
}

func (w *Wizard) askYesNo(prompt string, defaultVal bool) bool {
	defaultStr := "y/N"
	if defaultVal {
		defaultStr = "Y/n"
	}
	fmt.Printf("%s [%s]: ", prompt, defaultStr)
	text, _ := w.reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))
	if text == "" {
		return defaultVal
	}
	return text == "y" || text == "yes"
}

func (cfg *ProjectConfig) ToSpec() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("project: %s\n", strings.ToLower(strings.ReplaceAll(cfg.Name, " ", "-"))))
	if cfg.Description != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", cfg.Description))
	}
	sb.WriteString("\nmodules:\n")
	sb.WriteString(fmt.Sprintf("  - name: core\n    path: %s\n", cfg.ModulePath))
	sb.WriteString("\nservices:\n")
	sb.WriteString(fmt.Sprintf("  - name: api\n    kind: http\n    port: %d\n", cfg.Port))
	sb.WriteString("\narchitecture:\n")
	sb.WriteString(fmt.Sprintf("  pattern: %s\n", cfg.Architecture))
	sb.WriteString("\ndeployment:\n")
	sb.WriteString(fmt.Sprintf("  strategy: %s\n", cfg.Deployment))
	if cfg.EnableTesting {
		sb.WriteString("\ntesting:\n")
		sb.WriteString("  strategy: unit\n")
		sb.WriteString("  coverage: standard\n")
	}
	return sb.String()
}
