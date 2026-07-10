package model

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model/generation"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/service"
)

func TestNEIRZeroValue(t *testing.T) {
	neir := &NEIR{}
	if neir.Project != nil {
		t.Error("zero-value NEIR should have nil Project")
	}
	if neir.Modules != nil {
		t.Error("zero-value NEIR should have nil Modules")
	}
	if neir.Services != nil {
		t.Error("zero-value NEIR should have nil Services")
	}
	if neir.Generation != nil {
		t.Error("zero-value NEIR should have nil Generation")
	}
}

func TestNEIRWithProject(t *testing.T) {
	neir := &NEIR{
		Project: &project.Project{
			Name: "test-api",
		},
	}
	if neir.Project.Name != "test-api" {
		t.Errorf("Project.Name = %q, want %q", neir.Project.Name, "test-api")
	}
}

func TestNEIRWithModules(t *testing.T) {
	neir := &NEIR{
		Modules: []module.Module{
			{Name: "auth", Path: "./internal/auth"},
			{Name: "user", Path: "./internal/user"},
		},
	}
	if len(neir.Modules) != 2 {
		t.Errorf("Modules has %d entries, want 2", len(neir.Modules))
	}
	if neir.Modules[0].Name != "auth" {
		t.Errorf("Modules[0].Name = %q, want %q", neir.Modules[0].Name, "auth")
	}
}

func TestNEIRWithServices(t *testing.T) {
	neir := &NEIR{
		Services: []service.Service{
			{Name: "gateway", Kind: "http", Port: 8080},
		},
	}
	if len(neir.Services) != 1 {
		t.Errorf("Services has %d entries, want 1", len(neir.Services))
	}
	if neir.Services[0].Port != 8080 {
		t.Errorf("Services[0].Port = %d, want 8080", neir.Services[0].Port)
	}
}

func TestNEIRWithGeneration(t *testing.T) {
	neir := &NEIR{
		Generation: &generation.GenerationConfig{
			Languages: []language.Language{language.LanguageGo, language.LanguageTypeScript},
			OutputDir: "./generated",
		},
	}
	if !neir.Generation.HasLanguage(language.LanguageGo) {
		t.Error("Generation should contain Go")
	}
	if !neir.Generation.HasLanguage(language.LanguageTypeScript) {
		t.Error("Generation should contain TypeScript")
	}
	if neir.Generation.HasLanguage(language.LanguagePython) {
		t.Error("Generation should not contain Python")
	}
}
