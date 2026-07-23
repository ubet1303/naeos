package parser

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func FuzzParse(f *testing.F) {
	f.Add(`project: test`)
	f.Add(`project: test
modules:
  - name: core
    path: ./core`)
	f.Add(`project: test
modules:
  - name: auth
    dependencies: [core]
  - name: core`)
	f.Add(`{invalid yaml`)
	f.Add(`project: ""`)
	f.Add(`modules: []`)
	f.Add(`services:
  - name: api
    port: 8080
    endpoints:
      - method: GET
        path: /test`)
	f.Add(`version: 0.3.0
project: versioned
architecture:
  pattern: hexagonal`)

	f.Fuzz(func(t *testing.T, input string) {
		p := NewParser(".")
		doc, err := p.Parse(input)

		if err != nil {
			if doc != nil {
				t.Error("doc should be nil on error")
			}
			return
		}

		if doc == nil {
			t.Error("doc should not be nil on success")
			return
		}

		if doc.Project == "" && len(doc.Modules) == 0 && len(doc.Services) == 0 {
			if doc.Raw != input {
				t.Error("raw should match input")
			}
		}
	})
}

func FuzzParseYAMLNode(f *testing.F) {
	f.Add(`project: test`)
	f.Add(`key: value`)
	f.Add(`list:
  - item1
  - item2`)
	f.Add(`nested:
  inner: value`)
	f.Add(`null_val: null`)
	f.Add(`bool_val: true`)
	f.Add(`int_val: 42`)
	f.Add(`float_val: 3.14`)

	f.Fuzz(func(t *testing.T, input string) {
		var root yaml.Node
		if err := yaml.Unmarshal([]byte(input), &root); err != nil {
			return
		}

		if len(root.Content) == 0 {
			return
		}

		result, err := parseYAMLNode(root.Content[0])
		if err != nil {
			return
		}

		_ = result
	})
}

func FuzzVariableResolver(f *testing.F) {
	f.Add("${var}", "value")
	f.Add("$env{HOME}", "/home/user")
	f.Add("no variables here", "")
	f.Add("${a}-${b}", "x")
	f.Add("$env{NONEXISTENT}", "")

	f.Fuzz(func(t *testing.T, input, value string) {
		resolver := NewVariableResolver()
		resolver.SetVar("var", value)
		resolver.SetVar("a", value)
		resolver.SetVar("b", value)

		result, err := resolver.Resolve(input)
		if err != nil {
			t.Fatal(err)
		}

		if result == "" && input != "" {
			if !strings.Contains(input, "${") && !strings.Contains(input, "$env{") {
				if result != input {
					t.Errorf("expected %q, got %q", input, result)
				}
			}
		}
	})
}

func FuzzSchemaVersionParse(f *testing.F) {
	f.Add("0.1.0")
	f.Add("0.3.0")
	f.Add("1.0.0")
	f.Add("v2.0.0")
	f.Add("invalid")
	f.Add("")
	f.Add("1.2.3.4")
	f.Add("0.0")
	f.Add("-1.0.0")
	f.Add("abc.def.ghi")

	f.Fuzz(func(t *testing.T, input string) {
		result, err := ParseSchemaVersion(input)
		if err != nil {
			return
		}

		if result.Major < 0 {
			t.Error("major version should not be negative")
		}
		if result.Minor < 0 {
			t.Error("minor version should not be negative")
		}
		if result.Patch < 0 {
			t.Error("patch version should not be negative")
		}
	})
}

func FuzzValidateModules(f *testing.F) {
	f.Add("auth", "core")
	f.Add("api", "auth")
	f.Add("core", "")
	f.Add("", "something")
	f.Add("a", "b")
	f.Add("b", "a")

	f.Fuzz(func(t *testing.T, name, dep string) {
		modules := []Module{
			{Name: name, Dependencies: []string{dep}},
		}

		v := NewSpecValidator()
		issues := v.ValidateModules(modules)

		for _, issue := range issues {
			if issue.Severity != "error" && issue.Severity != "warning" {
				t.Errorf("invalid severity: %s", issue.Severity)
			}
			if issue.Rule == "" {
				t.Error("rule should not be empty")
			}
		}
	})
}

func FuzzSlugify(f *testing.F) {
	f.Add("Hello World")
	f.Add("my-project_name.test")
	f.Add("")
	f.Add("   spaces   ")
	f.Add("UPPERCASE")
	f.Add("special!@#$%chars")
	f.Add("already-slugified")
	f.Add("äöü ñ café")
	f.Add("123-numbers")
	f.Add("--dashes--")

	f.Fuzz(func(t *testing.T, input string) {
		result := Slugify(input)
		if result == "" {
			t.Error("slugify should never return empty")
		}
		for _, c := range result {
			if c < 'a' || c > 'z' {
				if c < '0' || c > '9' {
					if c != '-' {
						t.Errorf("slug contains unexpected character: %c", c)
					}
				}
			}
		}
	})
}

func FuzzCheckSpecVersion(f *testing.F) {
	f.Add("0.3.0")
	f.Add("0.1.0")
	f.Add("1.0.0")
	f.Add("")
	f.Add("invalid")
	f.Add("v2.0.0")
	f.Add("0.0.1")

	f.Fuzz(func(t *testing.T, input string) {
		result := CheckSpecVersion(input)
		if result == nil {
			t.Fatal("result should not be nil")
		}
		if result.Message == "" {
			t.Error("message should not be empty")
		}
	})
}

func FuzzDefaultProjectNameForInput(f *testing.F) {
	f.Add("my cool project")
	f.Add("hello-world")
	f.Add("")
	f.Add("   spaces everywhere   ")
	f.Add("UPPERCASE PROJECT")
	f.Add("project.with.dots")
	f.Add("123-numeric-start")

	f.Fuzz(func(t *testing.T, input string) {
		result := DefaultProjectNameForInput(input)
		if result == "" {
			t.Error("should not return empty")
		}
	})
}

func FuzzDefaultModuleNameForProject(f *testing.F) {
	f.Add("my-project")
	f.Add("auth service")
	f.Add("")
	f.Add("UPPERCASE")
	f.Add("project-with-many-modules")

	f.Fuzz(func(t *testing.T, input string) {
		result := DefaultModuleNameForProject(input)
		if result == "" {
			t.Error("should not return empty")
		}
	})
}
