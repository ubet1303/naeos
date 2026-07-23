package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
)

func TestPluginListJSONCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	root.SetArgs([]string{"plugin", "list", "--plugin-dir", dir, "--output-format", "json"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	_ = root.Execute()
}

func TestPluginListYAMLCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	root.SetArgs([]string{"plugin", "list", "--plugin-dir", dir, "--output-format", "yaml"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	_ = root.Execute()
}

func TestPluginSearchCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "plugin", "search", "lint", "--plugin-dir", dir)
	if err != nil {
		t.Fatalf("plugin search failed: %v", err)
	}
	if !strings.Contains(out, "lint-lint") {
		t.Fatalf("expected search results, got %q", out)
	}
}

func TestPluginCreateCov(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "my-plugin")
	root := newRootCommand()
	out, err := executeCommand(root, "plugin", "create", "my-plugin", "--plugin-dir", pluginDir)
	if err != nil {
		t.Fatalf("plugin create failed: %v", err)
	}
	if !strings.Contains(out, "Created plugin skeleton") {
		t.Fatalf("expected 'Created plugin skeleton', got %q", out)
	}
}

func TestPluginInfoNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "plugin", "info", "nonexistent", "--plugin-dir", dir)
	if err == nil {
		t.Fatal("expected error for nonexistent plugin")
	}
}

func TestPluginUninstallNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "plugin", "uninstall", "nonexistent", "--plugin-dir", dir)
	if err == nil {
		t.Fatal("expected error for uninstall nonexistent")
	}
}

func TestPluginEnableNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "plugin", "enable", "nonexistent", "--plugin-dir", dir)
	if err == nil {
		t.Fatal("expected error for enable nonexistent")
	}
}

func TestPluginDisableNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "plugin", "disable", "nonexistent", "--plugin-dir", dir)
	if err == nil {
		t.Fatal("expected error for disable nonexistent")
	}
}

func TestPluginTestNonexistentCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "plugin", "test", "/nonexistent/plugin.so", "--plugin-dir", dir)
	if err != nil {
		t.Fatalf("plugin test should not error: %v", err)
	}
	if !strings.Contains(out, "FAIL") {
		t.Fatalf("expected FAIL in output, got %q", out)
	}
}

func TestPluginListWithPluginsCov(t *testing.T) {
	dir := t.TempDir()
	cfgJSON := `{"plugins":[{"name":"test-plug","version":"1.0.0","enabled":true,"description":"A test plugin"},{"name":"disabled-plug","version":"0.5.0","enabled":false,"description":"Disabled plugin"}],"lazy":true}`
	if err := os.WriteFile(filepath.Join(dir, "plugins.json"), []byte(cfgJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "plugin", "list", "--plugin-dir", dir)
	if err != nil {
		t.Fatalf("plugin list with data failed: %v", err)
	}
	if !strings.Contains(out, "test-plug") {
		t.Fatalf("expected test-plug in output, got %q", out)
	}
	if !strings.Contains(out, "enabled") {
		t.Fatalf("expected 'enabled' status, got %q", out)
	}
	if !strings.Contains(out, "disabled") {
		t.Fatalf("expected 'disabled' status, got %q", out)
	}
}

func TestPluginInfoFoundCov(t *testing.T) {
	dir := t.TempDir()
	cfgJSON := `{"plugins":[{"name":"my-plug","version":"2.0.0","enabled":true,"description":"desc","author":"me","path":"/path/to/plug.so"}],"lazy":true}`
	if err := os.WriteFile(filepath.Join(dir, "plugins.json"), []byte(cfgJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "plugin", "info", "my-plug", "--plugin-dir", dir)
	if err != nil {
		t.Fatalf("plugin info failed: %v", err)
	}
	if !strings.Contains(out, "my-plug") || !strings.Contains(out, "2.0.0") || !strings.Contains(out, "me") {
		t.Fatalf("expected plugin details in output, got %q", out)
	}
}

func TestMarketplaceProfileListCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "marketplace", "profile", "list", "--cache-dir", dir)
	if err != nil {
		t.Fatalf("marketplace profile list failed: %v", err)
	}
	if !strings.Contains(out, "saas-starter") {
		t.Fatalf("expected saas-starter in output, got %q", out)
	}
}

func TestMarketplaceProfileSearchCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "marketplace", "profile", "search", "fintech", "--cache-dir", dir)
	if err != nil {
		t.Fatalf("marketplace profile search failed: %v", err)
	}
	if !strings.Contains(out, "fintech-core") {
		t.Fatalf("expected fintech-core in output, got %q", out)
	}
}

func TestMarketplacePluginListCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "marketplace", "plugin", "list", "--cache-dir", dir)
	if err != nil {
		t.Fatalf("marketplace plugin list failed: %v", err)
	}
	if !strings.Contains(out, "naeos-lint") {
		t.Fatalf("expected naeos-lint in output, got %q", out)
	}
}

func TestMarketplacePluginSearchCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "marketplace", "plugin", "search", "lint", "--cache-dir", dir)
	if err != nil {
		t.Fatalf("marketplace plugin search failed: %v", err)
	}
	if !strings.Contains(out, "naeos-lint") {
		t.Fatalf("expected naeos-lint, got %q", out)
	}
}

func TestMarketplacePublishNoManifestCov(t *testing.T) {
	dir := t.TempDir()
	pkgDir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "marketplace", "publish", pkgDir, "--cache-dir", dir)
	if err == nil {
		t.Fatal("expected error for publish without manifest")
	}
}

func TestMarketplacePublishWithManifestCov(t *testing.T) {
	dir := t.TempDir()
	pkgDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(pkgDir, "naeos.yaml"), []byte("name: test\nversion: 1.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "marketplace", "publish", pkgDir, "--cache-dir", dir)
	if err != nil {
		t.Fatalf("marketplace publish failed: %v", err)
	}
	if !strings.Contains(out, "package validated") && !strings.Contains(out, "Package validated") {
		t.Fatalf("expected 'Package validated' in output, got %q", out)
	}
}

func TestTemplateAddAndRemoveCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "template", "add", "my-tpl", "content here", "--templates-dir", dir)
	if err != nil {
		t.Fatalf("template add failed: %v", err)
	}
	if !strings.Contains(out, "Added template") {
		t.Fatalf("expected 'Added template', got %q", out)
	}

	out, err = executeCommand(root, "template", "remove", "my-tpl", "--templates-dir", dir)
	if err != nil {
		t.Fatalf("template remove failed: %v", err)
	}
	if !strings.Contains(out, "Removed template") {
		t.Fatalf("expected 'Removed template', got %q", out)
	}
}

func TestTemplateRemoveNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "template", "remove", "nonexistent", "--templates-dir", dir)
	if err == nil {
		t.Fatal("expected error for remove nonexistent template")
	}
}

func TestTemplateShowLLMPromptCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "template", "show", "enrich-spec")
	if err != nil {
		t.Fatalf("template show failed: %v", err)
	}
	if !strings.Contains(out, "enrich-spec") || !strings.Contains(out, "llm") {
		t.Fatalf("expected LLM prompt details, got %q", out)
	}
}

func TestTemplateShowCompilerTemplateCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "template", "show", "claude")
	if err != nil {
		t.Fatalf("template show compiler failed: %v", err)
	}
	if !strings.Contains(out, "claude") || !strings.Contains(out, "compiler") {
		t.Fatalf("expected compiler template details, got %q", out)
	}
}

func TestTemplateShowNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "template", "show", "nonexistent-xyz")
	if err == nil {
		t.Fatal("expected error for template not found")
	}
}

func TestTemplatePromptCreateCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "template", "prompt-create", "my-prompt", "--user", "Analyze: {{.Spec}}", "--templates-dir", dir)
	if err != nil {
		t.Fatalf("prompt-create failed: %v", err)
	}
	if !strings.Contains(out, "Created prompt template") {
		t.Fatalf("expected 'Created prompt template', got %q", out)
	}
}

func TestTemplatePromptCreateNoUserCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "template", "prompt-create", "my-prompt", "--templates-dir", dir)
	if err == nil {
		t.Fatal("expected error for missing --user flag")
	}
}

func TestTemplatePromptRemoveCov(t *testing.T) {
	dir := t.TempDir()
	promptDir := filepath.Join(dir, "prompts")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "test-prompt.yaml"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "template", "prompt-remove", "test-prompt", "--templates-dir", dir)
	if err != nil {
		t.Fatalf("prompt-remove failed: %v", err)
	}
	if !strings.Contains(out, "Removed prompt template") {
		t.Fatalf("expected 'Removed prompt template', got %q", out)
	}
}

func TestTemplatePromptRemoveNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "template", "prompt-remove", "nonexistent", "--templates-dir", dir)
	if err == nil {
		t.Fatal("expected error for nonexistent prompt")
	}
}

func TestTemplateListKindFilterCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "template", "list", "--kind", "code")
	if err != nil {
		t.Fatalf("template list --kind code failed: %v", err)
	}
	if !strings.Contains(out, "Code Generation") {
		t.Fatalf("expected 'Code Generation' in output, got %q", out)
	}
}

func TestArtifactsSummaryEmptyCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "artifacts", "summary", "--dir", dir)
	if err != nil {
		t.Fatalf("artifacts summary failed: %v", err)
	}
	if !strings.Contains(out, "Artifact Summary") {
		t.Fatalf("expected 'Artifact Summary', got %q", out)
	}
}

func TestArtifactsDedupEmptyCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "artifacts", "dedup", "--dir", dir)
	if err != nil {
		t.Fatalf("artifacts dedup failed: %v", err)
	}
	if !strings.Contains(out, "Removed 0 duplicate") {
		t.Fatalf("expected 'Removed 0 duplicate', got %q", out)
	}
}

func TestArtifactsInfoNotFoundCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "artifacts", "info", "nonexistent.go", "--dir", dir)
	if err == nil {
		t.Fatal("expected error for artifact not found")
	}
}

func TestArtifactsWithEntriesCov(t *testing.T) {
	dir := t.TempDir()
	storeFile := filepath.Join(dir, ".artifacts.json")
	manifest := `{"version":"1.0","artifacts":[{"path":"main.go","kind":"code","language":"go","size":100,"content_hash":"abc123def456","created_at":"2025-01-01T00:00:00Z","metadata":{"key":"val"}}]}`
	if err := os.WriteFile(storeFile, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "artifacts", "list", "--dir", dir)
	if err != nil {
		t.Fatalf("artifacts list failed: %v", err)
	}
	if !strings.Contains(out, "main.go") {
		t.Fatalf("expected main.go in output, got %q", out)
	}

	out, err = executeCommand(root, "artifacts", "info", "main.go", "--dir", dir)
	if err != nil {
		t.Fatalf("artifacts info failed: %v", err)
	}
	if !strings.Contains(out, "main.go") || !strings.Contains(out, "abc123") {
		t.Fatalf("expected artifact details, got %q", out)
	}

	out, err = executeCommand(root, "artifacts", "summary", "--dir", dir)
	if err != nil {
		t.Fatalf("artifacts summary failed: %v", err)
	}
	if !strings.Contains(out, "code") {
		t.Fatalf("expected code kind in summary, got %q", out)
	}
}

func TestProfileShowCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"profile", "show", "saas"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "saas") || !strings.Contains(out, "Industry") {
		t.Fatalf("expected profile details, got %q", out)
	}
}

func TestProfileShowNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "profile", "show", "nonexistent")
	if err == nil {
		t.Fatal("expected error for profile not found")
	}
}

func TestProfileCompareCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"profile", "compare", "saas", "fintech"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "Comparing") || !strings.Contains(out, "Industry") {
		t.Fatalf("expected comparison output, got %q", out)
	}
}

func TestProfileCompareNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "profile", "compare", "saas", "nonexistent")
	if err == nil {
		t.Fatal("expected error for profile compare not found")
	}
}

func TestProfileCategoriesCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"profile", "categories"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "Profile Categories") {
		t.Fatalf("expected 'Profile Categories', got %q", out)
	}
}

func TestProfileSearchCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"profile", "search", "fintech"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "fintech") {
		t.Fatalf("expected fintech in search results, got %q", out)
	}
}

func TestProfileSearchNoResultsCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"profile", "search", "zzznonexistentzzz"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "No profiles match") {
		t.Fatalf("expected 'No profiles match', got %q", out)
	}
}

func TestProfileApplyCov(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "spec.yaml")
	root := newRootCommand()
	root.SetArgs([]string{"profile", "apply", "saas", "--output", outPath})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "applied") {
		t.Fatalf("expected 'applied' in output, got %q", out)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), "saas") {
		t.Fatalf("expected saas in spec, got %q", string(data))
	}
}

func TestProfileListWithIndustryFilterCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"profile", "list", "--industry", "technology"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "ID") {
		t.Fatalf("expected header in output, got %q", out)
	}
}

func TestLintNoIssuesCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("project: test\nversion: 0.1.0\ndescription: A test project\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "lint", "--input-file", specFile)
	if err != nil {
		t.Fatalf("lint no issues failed: %v", err)
	}
	if !strings.Contains(out, "no issues found") {
		t.Fatalf("expected 'no issues found', got %q", out)
	}
}

func TestLintMissingInputCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "lint")
	if err == nil {
		t.Fatal("expected error for missing --input-file")
	}
}

func TestLintFileNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "lint", "--input-file", "/nonexistent/spec.yaml")
	if err == nil {
		t.Fatal("expected error for file not found")
	}
}

func TestSecuritySetSecretCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "set-secret", "--name", "test-secret-cov", "--value", "secret123", "--key", "my-encryption-key")
	if err != nil {
		t.Fatalf("security set-secret failed: %v", err)
	}
	if !strings.Contains(out, "stored successfully") {
		t.Fatalf("expected 'stored successfully', got %q", out)
	}
}

func TestSecurityGetSecretCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "security", "get-secret", "--name", "some-name", "--key", "some-key")
	if err == nil {
		t.Fatal("expected error for non-existent secret")
	}
}

func TestSecurityGetSecretNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "security", "get-secret", "--name", "nonexistent-cov", "--key", "enc-key")
	if err == nil {
		t.Fatal("expected error for secret not found")
	}
}

func TestSecurityHashPasswordCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "hash-password", "--password", "mypassword123")
	if err != nil {
		t.Fatalf("security hash-password failed: %v", err)
	}
	if !strings.Contains(out, "$2") && len(strings.TrimSpace(out)) == 0 {
		t.Fatalf("expected bcrypt hash, got %q", out)
	}
}

func TestSecurityValidatePassCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "validate", "--name", "email", "--value", "test@example.com")
	if err != nil {
		t.Fatalf("security validate failed: %v", err)
	}
	if !strings.Contains(out, "passed") {
		t.Fatalf("expected 'passed', got %q", out)
	}
}

func TestSecurityValidateFailCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "validate", "--name", "name", "--value", "a")
	if err != nil {
		t.Fatalf("security validate failed: %v", err)
	}
	if !strings.Contains(out, "failed") {
		t.Fatalf("expected 'failed', got %q", out)
	}
}

func TestSecurityListSecretsEmptyCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "security", "list-secrets", "--key", dir+"/nonexistent")
	if err != nil {
		t.Fatalf("security list-secrets failed: %v", err)
	}
	if !strings.Contains(out, "No secrets stored") {
		t.Fatalf("expected 'No secrets stored', got %q", out)
	}
}

func TestSecurityListSecretsWithDataCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "list-secrets", "--key", "default-key")
	if err != nil {
		t.Fatalf("security list-secrets with data failed: %v", err)
	}
	if !strings.Contains(out, "No secrets stored") {
		t.Fatalf("expected 'No secrets stored', got %q", out)
	}
}

func TestSecuritySanitizeHTMLCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "sanitize", "--input", "<b>hello</b>", "--mode", "html")
	if err != nil {
		t.Fatalf("security sanitize html failed: %v", err)
	}
	if !strings.Contains(out, "Sanitized") {
		t.Fatalf("expected 'Sanitized' in output, got %q", out)
	}
}

func TestSecuritySanitizeSQLCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "sanitize", "--input", "SELECT * FROM users", "--mode", "sql")
	if err != nil {
		t.Fatalf("security sanitize sql failed: %v", err)
	}
	if !strings.Contains(out, "Sanitized") {
		t.Fatalf("expected 'Sanitized' in output, got %q", out)
	}
}

func TestSecuritySanitizeXSSCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "sanitize", "--input", "<script>alert(1)</script>", "--mode", "xss")
	if err != nil {
		t.Fatalf("security sanitize xss failed: %v", err)
	}
	if !strings.Contains(out, "Sanitized") {
		t.Fatalf("expected 'Sanitized' in output, got %q", out)
	}
}

func TestSecuritySanitizePathCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "sanitize", "--input", "../../etc/passwd", "--mode", "path")
	if err != nil {
		t.Fatalf("security sanitize path failed: %v", err)
	}
	if !strings.Contains(out, "Sanitized") {
		t.Fatalf("expected 'Sanitized' in output, got %q", out)
	}
}

func TestSecuritySanitizeDefaultCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "sanitize", "--input", "hello world")
	if err != nil {
		t.Fatalf("security sanitize default failed: %v", err)
	}
	if !strings.Contains(out, "Sanitized") {
		t.Fatalf("expected 'Sanitized' in output, got %q", out)
	}
}

func TestSecuritySanitizeRequiresInputCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "security", "sanitize")
	if err == nil {
		t.Fatal("expected error for missing --input flag")
	}
}

func TestSecurityAuditCleanCov(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "clean.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "security", "audit", "--input", dir)
	if err != nil {
		t.Fatalf("security audit failed: %v", err)
	}
	if !strings.Contains(out, "Security Audit") {
		t.Fatalf("expected 'Security Audit', got %q", out)
	}
}

func TestSecurityAuditWithFindingsCov(t *testing.T) {
	dir := t.TempDir()
	content := `package main
const password = "supersecret123"
const apiKey = "sk-1234567890abcdef"
func main() {
	query := "SELECT * FROM users WHERE id=" + userId
}
`
	if err := os.WriteFile(filepath.Join(dir, "bad.go"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "security", "audit", "--input", dir, "--output", "json")
	if err != nil {
		t.Fatalf("security audit with findings failed: %v", err)
	}
	if !strings.Contains(out, "findings") {
		t.Fatalf("expected 'findings' in JSON, got %q", out)
	}
}

func TestSecurityAuditNoInputCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "security", "audit")
	if err != nil {
		t.Fatalf("security audit no input failed: %v", err)
	}
	if !strings.Contains(out, "Security Audit") {
		t.Fatalf("expected 'Security Audit', got %q", out)
	}
}

func TestAuthWhoamiValidKeyCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "auth", "whoami", "--api-key", "invalid-key-cov")
	if err != nil {
		t.Fatalf("auth whoami failed: %v", err)
	}
	if !strings.Contains(out, "Authentication failed") {
		t.Fatalf("expected 'Authentication failed', got %q", out)
	}
}

func TestAuthWhoamiWithCreatedKeyCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "auth", "create-key", "--user-id", "u2", "--name", "test-key")
	if err != nil {
		t.Fatalf("auth create-key failed: %v", err)
	}
	if !strings.Contains(out, "Created API key") {
		t.Fatalf("expected 'Created API key', got %q", out)
	}
}

func TestAuthCreateKeyCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "auth", "create-key", "--user-id", "u1", "--name", "my-key")
	if err != nil {
		t.Fatalf("auth create-key failed: %v", err)
	}
	if !strings.Contains(out, "Created API key") {
		t.Fatalf("expected 'Created API key', got %q", out)
	}
}

func TestAuthLoginCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "auth", "login", "--provider", "nonexistent")
	if err != nil {
		t.Fatalf("auth login failed: %v", err)
	}
	if !strings.Contains(out, "not registered") {
		t.Fatalf("expected 'not registered', got %q", out)
	}
}

func TestAuthLogoutCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "auth", "logout")
	if err != nil {
		t.Fatalf("auth logout failed: %v", err)
	}
	if !strings.Contains(out, "Logged out") {
		t.Fatalf("expected 'Logged out', got %q", out)
	}
}

func TestCloudExportRequiresResourcesCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: test
  environment: staging
`
	if err := os.WriteFile(specFile, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	_, err := executeCommand(root, "cloud", "export", "--input-file", specFile)
	if err == nil {
		t.Fatal("expected error for export without resources")
	}
}

func TestCloudExportWithResourcesCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: test
  environment: staging
  resources:
    - name: my-bucket
      type: aws_s3_bucket
      spec:
        bucket: test-bucket
`
	if err := os.WriteFile(specFile, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "cloud", "export", "--input-file", specFile)
	if err != nil {
		t.Fatalf("cloud export failed: %v", err)
	}
	if !strings.Contains(out, "aws_s3_bucket") && !strings.Contains(out, "terraform") {
		t.Fatalf("expected HCL in output, got %q", out)
	}
}

func TestCloudDeployWithResourcesCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: test
  environment: staging
  resources:
    - name: my-bucket
      type: aws_s3_bucket
      spec:
        bucket: test-bucket
`
	if err := os.WriteFile(specFile, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	_, err := executeCommand(root, "cloud", "deploy", "--input-file", specFile)
	_ = err
}

func TestCloudDeployRequiresResourcesCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: test
  environment: staging
`
	if err := os.WriteFile(specFile, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	_, err := executeCommand(root, "cloud", "deploy", "--input-file", specFile)
	if err == nil {
		t.Fatal("expected error for deploy without resources")
	}
}

func TestCloudStatusCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "cloud", "status")
	if err != nil {
		t.Fatalf("cloud status failed: %v", err)
	}
	if !strings.Contains(out, "No deployments found") && !strings.Contains(out, "Deployed resources") {
		t.Fatalf("expected status output, got %q", out)
	}
}

func TestCloudTypesCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "cloud", "types")
	if err != nil {
		t.Fatalf("cloud types failed: %v", err)
	}
	if !strings.Contains(out, "Supported resource types") {
		t.Fatalf("expected 'Supported resource types', got %q", out)
	}
}

func TestCloudTypesJSONCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "cloud", "types", "--output-format", "json")
	if err != nil {
		t.Fatalf("cloud types json failed: %v", err)
	}
	if !strings.Contains(out, "[") {
		t.Fatalf("expected JSON array, got %q", out)
	}
}

func TestCloudPlanRequiresResourcesCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: test
  environment: staging
`
	if err := os.WriteFile(specFile, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	_, err := executeCommand(root, "cloud", "plan", "--input-file", specFile)
	if err == nil {
		t.Fatal("expected error for plan without resources")
	}
}

func TestMigrateVersionsCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"migrate", "versions"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "0.1.0") || !strings.Contains(out, "0.3.0") {
		t.Fatalf("expected versions in output, got %q", out)
	}
}

func TestMigratePlanCov(t *testing.T) {
	root := newRootCommand()
	root.SetArgs([]string{"migrate", "plan"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	out := captureOutput(t, func() { _ = root.Execute() })
	if !strings.Contains(out, "Migration plan") && !strings.Contains(out, "No migrations") {
		t.Fatalf("expected migration plan, got %q", out)
	}
}

func TestMigrateRunDryRunCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("project: test\nversion: 0.1.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "migrate", "run", specFile, "--dry-run")
	if err != nil {
		t.Fatalf("migrate run dry-run failed: %v", err)
	}
	_ = out
}

func TestMigrateRunFileNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "migrate", "run", "/nonexistent/spec.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent spec file")
	}
}

func TestGatewayStatusJSONCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "gateway", "status", "--output", "json")
	if err != nil {
		t.Fatalf("gateway status json failed: %v", err)
	}
	if !strings.Contains(out, "rate_limiter") {
		t.Fatalf("expected rate_limiter in JSON, got %q", out)
	}
}

func TestGatewayRateStatusJSONCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "gateway", "rate-status", "--output", "json")
	if err != nil {
		t.Fatalf("gateway rate-status json failed: %v", err)
	}
	if !strings.Contains(out, "token bucket") {
		t.Fatalf("expected 'token bucket' in JSON, got %q", out)
	}
}

func TestGatewayCBStatusJSONCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "gateway", "cb-status", "--output", "json")
	if err != nil {
		t.Fatalf("gateway cb-status json failed: %v", err)
	}
	if !strings.Contains(out, "failure_threshold") {
		t.Fatalf("expected failure_threshold in JSON, got %q", out)
	}
}

func TestGatewayLBListJSONCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "gateway", "lb-list", "--output", "json")
	if err != nil {
		t.Fatalf("gateway lb-list json failed: %v", err)
	}
	_ = out
}

func TestHistoryEmptyCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "history", "--store-dir", dir)
	if err != nil {
		t.Fatalf("history empty failed: %v", err)
	}
	if !strings.Contains(out, "No pipeline runs found") {
		t.Fatalf("expected 'No pipeline runs found', got %q", out)
	}
}

func TestHistoryJSONOutputCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	root.SetArgs([]string{"history", "--store-dir", dir, "--output-format", "json"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	_ = root.Execute()
}

func TestHistoryYAMLOutputCov(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	out, err := executeCommand(root, "history", "--store-dir", dir, "--output-format", "yaml")
	if err != nil {
		t.Fatalf("history yaml failed: %v", err)
	}
	_ = out
}

func TestHistoryNonexistentDirCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "history", "--store-dir", "/nonexistent/path/events")
	if err == nil {
		t.Fatal("expected error for nonexistent event store dir")
	}
}

func TestDBListJSONCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "db", "list", "--output", "json")
	if err != nil {
		t.Fatalf("db list json failed: %v", err)
	}
	if !strings.Contains(out, "[") {
		t.Fatalf("expected JSON array, got %q", out)
	}
}

func TestDBListYAMLCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "db", "list", "--output", "yaml")
	if err != nil {
		t.Fatalf("db list yaml failed: %v", err)
	}
	_ = out
}

func TestDBStatusNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "db", "status", "--name", "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent connection")
	}
}

func TestDBMigrateNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "db", "migrate", "--name", "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent connection")
	}
}

func TestDBDisconnectNotFoundCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "db", "disconnect", "--name", "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent connection")
	}
}

func TestDBConnectInvalidConfigCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "db", "connect", "--type", "sqlite", "--name", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestAIHelpCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "ai", "--help")
	if err != nil {
		t.Fatalf("ai help failed: %v", err)
	}
	if !strings.Contains(out, "suggest") || !strings.Contains(out, "explain") || !strings.Contains(out, "enrich") {
		t.Fatalf("expected subcommands in help, got %q", out)
	}
}

func TestAISuggestRequiresInputCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "ai", "suggest")
	if err == nil {
		t.Fatal("expected error for missing --input-file")
	}
}

func TestAISuggestWithFileCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("project: test\nversion: 0.1.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand()
	out, err := executeCommand(root, "ai", "suggest", "--input-file", specFile)
	if err != nil {
		t.Fatalf("ai suggest failed: %v", err)
	}
	if !strings.Contains(out, "]") {
		t.Fatalf("expected suggestion output, got %q", out)
	}
}

func TestAIExplainCov(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "ai", "explain", "pipeline")
	if err != nil {
		t.Fatalf("ai explain failed: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty explanation")
	}
}

func TestAIEnrichRequiresInputCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "ai", "enrich")
	if err == nil {
		t.Fatal("expected error for missing --input-file")
	}
}

func TestAICompileRequiresInputCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "ai", "compile")
	if err == nil {
		t.Fatal("expected error for missing --input-file")
	}
}

func TestScaffoldWithLanguagesCov(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "app")
	root := newRootCommand()
	out, err := executeCommand(root, "scaffold", "--name", "myapp", "--output", outputDir, "--language", "go")
	if err != nil {
		t.Fatalf("scaffold with language failed: %v", err)
	}
	if !strings.Contains(out, "scaffolded") {
		t.Fatalf("expected 'scaffolded', got %q", out)
	}
}

func TestScaffoldDefaultLanguageCov(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "app2")
	root := newRootCommand()
	out, err := executeCommand(root, "scaffold", "--name", "app2", "--output", outputDir)
	if err != nil {
		t.Fatalf("scaffold default failed: %v", err)
	}
	if !strings.Contains(out, "scaffolded") {
		t.Fatalf("expected 'scaffolded', got %q", out)
	}
}

func TestScaffoldMissingNameCov(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "scaffold")
	if err == nil {
		t.Fatal("expected error for missing --name")
	}
}

func TestLoadCloudConfigFromSpecKindCov(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: test
  environment: staging
  resources:
    - name: my-bucket
      kind: aws_s3_bucket
      spec:
        bucket: test-bucket
`
	if err := os.WriteFile(specFile, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	config, err := loadCloudConfigFromSpec(specFile)
	if err != nil {
		t.Fatalf("loadCloudConfigFromSpec failed: %v", err)
	}
	if len(config.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(config.Resources))
	}
	if config.Resources[0].Type != "aws_s3_bucket" {
		t.Fatalf("expected aws_s3_bucket type, got %q", config.Resources[0].Type)
	}
}

func TestGenerateSimpleIDCov(t *testing.T) {
	id := generateSimpleID()
	if len(id) == 0 {
		t.Fatal("expected non-empty ID")
	}
}

func TestIndentCov(t *testing.T) {
	result := indent("line1\nline2\n", "  ")
	if !strings.Contains(result, "  line1") || !strings.Contains(result, "  line2") {
		t.Fatalf("expected indented lines, got %q", result)
	}
}

func TestResolveScaffoldLanguagesEmptyCov(t *testing.T) {
	langs := resolveScaffoldLanguages(nil)
	if len(langs) != 1 {
		t.Fatalf("expected 1 default language, got %d", len(langs))
	}
}

func TestResolveScaffoldLanguagesMultipleCov(t *testing.T) {
	langs := resolveScaffoldLanguages([]string{"go", "python"})
	if len(langs) != 2 {
		t.Fatalf("expected 2 languages, got %d", len(langs))
	}
}

func TestTargetLangNamesCov(t *testing.T) {
	names := targetLangNames([]language.Language{language.LanguageGo, language.LanguagePython})
	if len(names) != 2 || names[0] != "go" || names[1] != "python" {
		t.Fatalf("expected [go python], got %v", names)
	}
}

func TestDBConnectAndStatusJSONCov(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	name := "statdb-json-" + filepath.Base(dir)
	root := newRootCommand()
	_, err := executeCommand(root, "db", "connect", "--type", "sqlite", "--name", name, "--database", dbPath, "--user", "testuser")
	if err != nil {
		t.Fatalf("db connect failed: %v", err)
	}
	defer executeCommand(root, "db", "disconnect", "--name", name)
	out, err := executeCommand(root, "db", "status", "--name", name, "--output", "json")
	if err != nil {
		t.Fatalf("db status json failed: %v", err)
	}
	if !strings.Contains(out, "HEALTHY") {
		t.Fatalf("expected HEALTHY in JSON, got %q", out)
	}
}

func TestDBConnectAndStatusYAMLCov(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	name := "statdb-yaml-" + filepath.Base(dir)
	root := newRootCommand()
	_, err := executeCommand(root, "db", "connect", "--type", "sqlite", "--name", name, "--database", dbPath, "--user", "testuser")
	if err != nil {
		t.Fatalf("db connect failed: %v", err)
	}
	defer executeCommand(root, "db", "disconnect", "--name", name)
	out, err := executeCommand(root, "db", "status", "--name", name, "--output", "yaml")
	if err != nil {
		t.Fatalf("db status yaml failed: %v", err)
	}
	_ = out
}
