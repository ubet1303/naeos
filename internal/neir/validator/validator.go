package validator

import (
	"fmt"
	"strings"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
)

type Validator interface {
	Validate(neir any) error
}

type ValidationResult struct {
	Valid  bool
	Errors []string
	Warns  []string
}

type DefaultValidator struct{}

func NewValidator() Validator {
	return DefaultValidator{}
}

func (DefaultValidator) Validate(neir any) error {
	result := ValidateDetailed(neir)
	if !result.Valid {
		return fmt.Errorf("validation failed:\n  - %s", strings.Join(result.Errors, "\n  - "))
	}
	return nil
}

func ValidateDetailed(neir any) ValidationResult {
	result := ValidationResult{Valid: true}

	if neir == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "NEIR model is nil — ensure the specification was parsed correctly")
		return result
	}

	neirStruct, ok := neir.(*model.NEIR)
	if !ok {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("expected *model.NEIR, got %T — this is an internal error", neir))
		return result
	}

	if neirStruct.Project == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "project is required — add a 'project:' field to your specification")
	} else if neirStruct.Project.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "project name must not be empty — set 'project: <name>' in your specification")
	}

	if len(neirStruct.Modules) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "at least one module is required — add a 'modules:' section to your specification")
	}

	for i, mod := range neirStruct.Modules {
		if mod.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module[%d] name is required — each module needs a 'name:' field", i))
		}
		if mod.Path == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module %q (index %d) path is required — add a 'path:' field (e.g., ./internal/%s)", mod.Name, i, mod.Name))
		}
	}

	seenModules := make(map[string]int)
	for i, mod := range neirStruct.Modules {
		if prev, exists := seenModules[mod.Name]; exists {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("duplicate module name %q at index %d and %d — module names must be unique", mod.Name, prev, i))
		}
		seenModules[mod.Name] = i
	}

	for i, svc := range neirStruct.Services {
		if svc.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("service[%d] name is required — each service needs a 'name:' field", i))
		}
		if svc.Port < 0 || svc.Port > 65535 {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("service %q port %d is out of range — must be between 0 and 65535", svc.Name, svc.Port))
		}
	}

	seenPorts := make(map[int]string)
	for _, svc := range neirStruct.Services {
		if svc.Port > 0 {
			if prev, exists := seenPorts[svc.Port]; exists {
				result.Warns = append(result.Warns, fmt.Sprintf("service %q and %q share port %d — this may cause conflicts", prev, svc.Name, svc.Port))
			}
			seenPorts[svc.Port] = svc.Name
		}
	}

	if neirStruct.Metadata != nil {
		if neirStruct.Metadata.NEIRVersion == "" {
			result.Warns = append(result.Warns, "metadata.neir_version is recommended for traceability")
		}
	}

	if neirStruct.Generation != nil {
		if len(neirStruct.Generation.Languages) == 0 {
			result.Warns = append(result.Warns, "generation.languages is empty — defaulting to Go")
		}
		for _, lang := range neirStruct.Generation.Languages {
			if !language.IsValid(lang) {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported language %q — supported: go, typescript, python, java, rust", lang))
			}
		}
	}

	return result
}
