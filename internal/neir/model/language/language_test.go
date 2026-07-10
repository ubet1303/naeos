package language

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		lang Language
		want bool
	}{
		{LanguageGo, true},
		{LanguageTypeScript, true},
		{LanguagePython, true},
		{LanguageJava, true},
		{LanguageRust, true},
		{"csharp", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsValid(tt.lang); got != tt.want {
			t.Errorf("IsValid(%q) = %v, want %v", tt.lang, got, tt.want)
		}
	}
}

func TestAll(t *testing.T) {
	all := All()
	if len(all) != 5 {
		t.Errorf("All() returned %d languages, want 5", len(all))
	}
	seen := make(map[Language]bool)
	for _, l := range all {
		if seen[l] {
			t.Errorf("All() contains duplicate %q", l)
		}
		seen[l] = true
		if !IsValid(l) {
			t.Errorf("All() contains invalid language %q", l)
		}
	}
}

func TestExtensions(t *testing.T) {
	tests := []struct {
		lang Language
		want int
	}{
		{LanguageGo, 3},
		{LanguageTypeScript, 4},
		{LanguagePython, 3},
		{LanguageJava, 3},
		{LanguageRust, 2},
		{"unknown", 0},
	}
	for _, tt := range tests {
		exts := Extensions(tt.lang)
		if len(exts) != tt.want {
			t.Errorf("Extensions(%q) returned %d extensions, want %d", tt.lang, len(exts), tt.want)
		}
	}
}

func TestBuildFile(t *testing.T) {
	tests := []struct {
		lang Language
		want string
	}{
		{LanguageGo, "go.mod"},
		{LanguageTypeScript, "package.json"},
		{LanguagePython, "pyproject.toml"},
		{LanguageJava, "pom.xml"},
		{LanguageRust, "Cargo.toml"},
		{"unknown", ""},
	}
	for _, tt := range tests {
		if got := BuildFile(tt.lang); got != tt.want {
			t.Errorf("BuildFile(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}

func TestDockerBaseImage(t *testing.T) {
	tests := []struct {
		lang Language
		want string
	}{
		{LanguageGo, "golang:1.22-alpine"},
		{LanguageTypeScript, "node:22-alpine"},
		{LanguagePython, "python:3.12-slim"},
		{LanguageJava, "eclipse-temurin:21-jdk-alpine"},
		{LanguageRust, "rust:1.78-alpine"},
		{"unknown", "alpine:latest"},
	}
	for _, tt := range tests {
		if got := DockerBaseImage(tt.lang); got != tt.want {
			t.Errorf("DockerBaseImage(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}

func TestDockerRuntimeImage(t *testing.T) {
	tests := []struct {
		lang Language
		want string
	}{
		{LanguageGo, "alpine:3.19"},
		{LanguageTypeScript, "node:22-alpine"},
		{LanguagePython, "python:3.12-slim"},
		{LanguageJava, "eclipse-temurin:21-jre-alpine"},
		{LanguageRust, "alpine:3.19"},
		{"unknown", "alpine:latest"},
	}
	for _, tt := range tests {
		if got := DockerRuntimeImage(tt.lang); got != tt.want {
			t.Errorf("DockerRuntimeImage(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}
