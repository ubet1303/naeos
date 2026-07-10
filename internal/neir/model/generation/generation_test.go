package generation

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
)

func TestHasLanguage(t *testing.T) {
	gc := &GenerationConfig{
		Languages: []language.Language{language.LanguageGo, language.LanguageTypeScript},
	}
	tests := []struct {
		lang language.Language
		want bool
	}{
		{language.LanguageGo, true},
		{language.LanguageTypeScript, true},
		{language.LanguagePython, false},
		{language.LanguageJava, false},
		{language.LanguageRust, false},
	}
	for _, tt := range tests {
		if got := gc.HasLanguage(tt.lang); got != tt.want {
			t.Errorf("HasLanguage(%q) = %v, want %v", tt.lang, got, tt.want)
		}
	}
}

func TestHasLanguageEmpty(t *testing.T) {
	gc := &GenerationConfig{}
	if gc.HasLanguage(language.LanguageGo) {
		t.Error("HasLanguage() should return false for empty config")
	}
}

func TestDefaultLanguage(t *testing.T) {
	gc := &GenerationConfig{
		Languages: []language.Language{language.LanguageTypeScript, language.LanguageGo},
	}
	if got := gc.DefaultLanguage(); got != language.LanguageTypeScript {
		t.Errorf("DefaultLanguage() = %q, want %q", got, language.LanguageTypeScript)
	}
}

func TestDefaultLanguageEmpty(t *testing.T) {
	gc := &GenerationConfig{}
	if got := gc.DefaultLanguage(); got != language.LanguageGo {
		t.Errorf("DefaultLanguage() = %q, want %q (default)", got, language.LanguageGo)
	}
}

func TestGenerationConfigJSON(t *testing.T) {
	gc := &GenerationConfig{
		Languages: []language.Language{language.LanguageGo, language.LanguagePython},
		OutputDir: "./output",
		ModuleDir: "./internal",
	}
	if gc.OutputDir != "./output" {
		t.Errorf("OutputDir = %q, want %q", gc.OutputDir, "./output")
	}
	if gc.ModuleDir != "./internal" {
		t.Errorf("ModuleDir = %q, want %q", gc.ModuleDir, "./internal")
	}
	if len(gc.Languages) != 2 {
		t.Errorf("Languages has %d entries, want 2", len(gc.Languages))
	}
}
