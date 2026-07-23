package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "completion", "bash")
		if err != nil {
			t.Fatalf("completion bash failed: %v", err)
		}
	})
	if len(out) == 0 {
		t.Fatal("expected bash completion output")
	}
}

func TestCompletionZsh(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "completion", "zsh")
		if err != nil {
			t.Fatalf("completion zsh failed: %v", err)
		}
	})
	if len(out) == 0 {
		t.Fatal("expected zsh completion output")
	}
}

func TestCompletionFish(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "completion", "fish")
		if err != nil {
			t.Fatalf("completion fish failed: %v", err)
		}
	})
	if len(out) == 0 {
		t.Fatal("expected fish completion output")
	}
}

func TestCompletionPowerShell(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "completion", "powershell")
		if err != nil {
			t.Fatalf("completion powershell failed: %v", err)
		}
	})
	if len(out) == 0 {
		t.Fatal("expected powershell completion output")
	}
}

func TestCompletionInvalidShell(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "completion", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid shell")
	}
}

func TestBuildHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "build", "--help")
	if err != nil {
		t.Fatalf("build --help failed: %v", err)
	}
	if !strings.Contains(output, "Build artifacts") {
		t.Fatalf("expected build help, got %q", output)
	}
}

func TestBuildLocalDryRun(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "build", "--config", configPath, "--input", "project: demo", "--dry-run")
	if err != nil {
		t.Fatalf("build --dry-run failed: %v", err)
	}
	if !strings.Contains(output, "build=local") && !strings.Contains(output, "dry_run") {
		t.Fatalf("expected build output, got %q", output)
	}
}

func TestBuildDistributed(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "build", "--config", configPath, "--distributed", "--workers", "2")
	if err != nil {
		t.Fatalf("build distributed failed: %v", err)
	}
	if !strings.Contains(output, "distributed") {
		t.Fatalf("expected distributed output, got %q", output)
	}
}

func TestBuildJSONOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "build", "--config", configPath, "--input", "project: demo", "--output", "json")
	if err != nil {
		t.Fatalf("build json output failed: %v", err)
	}
	if !strings.Contains(output, `"pipeline"`) && !strings.Contains(output, `"build"`) {
		t.Fatalf("expected JSON output, got %q", output)
	}
}

func TestBuildLanguages(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "build", "--config", configPath, "--input", "project: demo", "--language", "go", "--language", "python")
	if err != nil {
		t.Fatalf("build with languages failed: %v", err)
	}
	if !strings.Contains(output, "build=local") {
		t.Fatalf("expected build output, got %q", output)
	}
}

func TestConfigShow(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "config", "show")
	if err != nil {
		t.Fatalf("config show failed: %v", err)
	}
	if !strings.Contains(output, "NAEOS Configuration Schema") {
		t.Fatalf("expected schema output, got %q", output)
	}
}

func TestConfigValidate(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: test\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "config", "validate", "--input", configPath)
	if err != nil {
		t.Fatalf("config validate failed: %v", err)
	}
	if !strings.Contains(output, "Config is valid") && !strings.Contains(output, "validation error") {
		t.Fatalf("expected validation output, got %q", output)
	}
}

func TestConfigEncryptDecrypt(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "config.yaml")
	encryptedPath := filepath.Join(dir, "config.enc")
	decryptedPath := filepath.Join(dir, "config.dec.yaml")

	if err := os.WriteFile(inputPath, []byte("pipeline:\n  name: secret\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	_, err := executeCommand(root, "config", "encrypt", "--input", inputPath, "--output", encryptedPath, "--passphrase", "test123")
	if err != nil {
		t.Fatalf("config encrypt failed: %v", err)
	}
	if _, err := os.Stat(encryptedPath); err != nil {
		t.Fatalf("expected encrypted file to exist: %v", err)
	}

	_, err = executeCommand(root, "config", "decrypt", "--input", encryptedPath, "--output", decryptedPath, "--passphrase", "test123")
	if err != nil {
		t.Fatalf("config decrypt failed: %v", err)
	}
	data, err := os.ReadFile(decryptedPath)
	if err != nil {
		t.Fatalf("read decrypted file: %v", err)
	}
	if !strings.Contains(string(data), "secret") {
		t.Fatalf("expected decrypted content to contain 'secret', got %q", string(data))
	}
}

func TestConfigEncryptRequiresInput(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "config", "encrypt", "--passphrase", "test")
	if err == nil {
		t.Fatal("expected error when --input is missing")
	}
}

func TestConfigEncryptRequiresPassphrase(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: test\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	_, err := executeCommand(root, "config", "encrypt", "--input", configPath)
	if err == nil {
		t.Fatal("expected error when --passphrase is missing")
	}
}

func TestConfigDecryptRequiresInput(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "config", "decrypt", "--passphrase", "test")
	if err == nil {
		t.Fatal("expected error when --input is missing")
	}
}

func TestConfigDecryptRequiresPassphrase(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, "config.enc")
	if err := os.WriteFile(encPath, []byte("encrypted-data"), 0o644); err != nil {
		t.Fatalf("write enc: %v", err)
	}

	root := newRootCommand()
	_, err := executeCommand(root, "config", "decrypt", "--input", encPath)
	if err == nil {
		t.Fatal("expected error when --passphrase is missing")
	}
}

func TestConfigEncryptStdout(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(inputPath, []byte("pipeline:\n  name: test\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "config", "encrypt", "--input", inputPath, "--passphrase", "test123")
	if err != nil {
		t.Fatalf("config encrypt stdout failed: %v", err)
	}
	if len(strings.TrimSpace(output)) == 0 {
		t.Fatal("expected encrypted output to stdout")
	}
}

func TestConfigDecryptStdout(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "config.yaml")
	encPath := filepath.Join(dir, "config.enc")

	if err := os.WriteFile(inputPath, []byte("pipeline:\n  name: test\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	_, err := executeCommand(root, "config", "encrypt", "--input", inputPath, "--output", encPath, "--passphrase", "test123")
	if err != nil {
		t.Fatalf("encrypt first: %v", err)
	}

	output, err := executeCommand(root, "config", "decrypt", "--input", encPath, "--passphrase", "test123")
	if err != nil {
		t.Fatalf("config decrypt stdout failed: %v", err)
	}
	if !strings.Contains(output, "test") {
		t.Fatalf("expected decrypted content, got %q", output)
	}
}

func TestContextCommand(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "context", "--input", "project: myapp\nversion: 1.0")
	if err != nil {
		t.Fatalf("context failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected context output")
	}
}

func TestContextJSONOutput(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "context", "--input", "project: myapp", "--output", "json")
	if err != nil {
		t.Fatalf("context json failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected JSON context output")
	}
}

func TestHealthCommand(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "health")
	if err != nil {
		t.Fatalf("health failed: %v", err)
	}
	if !strings.Contains(output, "NAEOS Health Report") {
		t.Fatalf("expected health report, got %q", output)
	}
}

func TestHealthJSONOutput(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "health", "--output-format", "json")
	if err != nil {
		t.Fatalf("health json failed: %v", err)
	}
	if !strings.Contains(output, `"status"`) {
		t.Fatalf("expected JSON health output, got %q", output)
	}
}

func TestHealthYAMLOutput(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "health", "--output-format", "yaml")
	if err != nil {
		t.Fatalf("health yaml failed: %v", err)
	}
	if !strings.Contains(output, "status:") {
		t.Fatalf("expected YAML health output, got %q", output)
	}
}

func TestHistoryCommand(t *testing.T) {
	dir := t.TempDir()
	storeDir := filepath.Join(dir, "events")
	if err := os.MkdirAll(storeDir, 0o755); err != nil {
		t.Fatalf("create events dir: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "history", "--store-dir", storeDir)
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if !strings.Contains(output, "No pipeline runs found") {
		t.Fatalf("expected no runs message, got %q", output)
	}
}

func TestLockGenerate(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	root := newRootCommand()
	output, err := executeCommand(root, "lock", "generate", testFile)
	if err != nil {
		t.Fatalf("lock generate failed: %v", err)
	}
	if !strings.Contains(output, "Generated") {
		t.Fatalf("expected generated message, got %q", output)
	}
	if _, err := os.Stat(filepath.Join(dir, "naeos.lock")); err != nil {
		t.Fatalf("expected lock file to exist: %v", err)
	}
}

func TestLockVerifyNoChanges(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	root := newRootCommand()
	_, err := executeCommand(root, "lock", "generate", testFile)
	if err != nil {
		t.Fatalf("lock generate failed: %v", err)
	}

	output, err := executeCommand(root, "lock", "verify", testFile)
	if err != nil {
		t.Fatalf("lock verify failed: %v", err)
	}
	if !strings.Contains(output, "no changes detected") {
		t.Fatalf("expected no changes, got %q", output)
	}
}

func TestLockHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "lock", "--help")
	if err != nil {
		t.Fatalf("lock --help failed: %v", err)
	}
	if !strings.Contains(output, "reproducible builds") {
		t.Fatalf("expected lock help, got %q", output)
	}
}

func TestRollbackListEmpty(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	root := newRootCommand()
	output, err := executeCommand(root, "rollback", "list")
	if err != nil {
		t.Fatalf("rollback list failed: %v", err)
	}
	if !strings.Contains(output, "No snapshots found") {
		t.Fatalf("expected no snapshots, got %q", output)
	}
}

func TestDocgenFull(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "docgen", "--input", "project: myapp\nversion: 1.0")
	if err != nil {
		t.Fatalf("docgen failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected docgen output")
	}
}

func TestDocgenAPI(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "docgen", "--input", "project: myapp", "--output", "api")
	if err != nil {
		t.Fatalf("docgen api failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected docgen api output")
	}
}

func TestDocgenModules(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "docgen", "--input", "project: myapp", "--output", "modules")
	if err != nil {
		t.Fatalf("docgen modules failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected docgen modules output")
	}
}

func TestDocsAPI(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "docs", "api", "--project", "my-app")
	if err != nil {
		t.Fatalf("docs api failed: %v", err)
	}
	if !strings.Contains(output, "Health check") && !strings.Contains(output, "GET") {
		t.Fatalf("expected API docs, got %q", output)
	}
}

func TestDocsAPIWithOutputDir(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "docs")

	root := newRootCommand()
	output, err := executeCommand(root, "docs", "api", "--project", "my-app", "--output", outDir)
	if err != nil {
		t.Fatalf("docs api with output failed: %v", err)
	}
	if !strings.Contains(output, "Generated api.md") {
		t.Fatalf("expected generated message, got %q", output)
	}
	if _, err := os.Stat(filepath.Join(outDir, "api.md")); err != nil {
		t.Fatalf("expected api.md to exist: %v", err)
	}
}

func TestDocsArchitecture(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "docs", "architecture", "--project", "my-app")
	if err != nil {
		t.Fatalf("docs architecture failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected architecture output")
	}
}

func TestDocsArchitectureWithOutputDir(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "docs")

	root := newRootCommand()
	output, err := executeCommand(root, "docs", "architecture", "--project", "my-app", "--output", outDir)
	if err != nil {
		t.Fatalf("docs architecture with output failed: %v", err)
	}
	if !strings.Contains(output, "Generated architecture.md") {
		t.Fatalf("expected generated message, got %q", output)
	}
}

func TestBenchmarkCommand(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "benchmark", "--config", configPath, "--iterations", "2")
	if err != nil {
		t.Fatalf("benchmark failed: %v", err)
	}
	if !strings.Contains(output, "NAEOS Benchmark Results") {
		t.Fatalf("expected benchmark results, got %q", output)
	}
}

func TestBenchmarkJSONOutput(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "benchmark", "--config", configPath, "--iterations", "2", "--output", "json")
	if err != nil {
		t.Fatalf("benchmark json failed: %v", err)
	}
	if !strings.Contains(output, `"iterations"`) {
		t.Fatalf("expected JSON benchmark output, got %q", output)
	}
}

func TestComplianceExportJSON(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "audit.json")

	root := newRootCommand()
	output, err := executeCommand(root, "compliance", "export", "--format", "json", "--output", outputPath)
	if err != nil {
		t.Fatalf("compliance export json failed: %v", err)
	}
	if !strings.Contains(output, "Compliance report exported") {
		t.Fatalf("expected export message, got %q", output)
	}
}

func TestComplianceExportCSV(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "audit.csv")

	root := newRootCommand()
	output, err := executeCommand(root, "compliance", "export", "--format", "csv", "--output", outputPath)
	if err != nil {
		t.Fatalf("compliance export csv failed: %v", err)
	}
	if !strings.Contains(output, "Compliance report exported") {
		t.Fatalf("expected export message, got %q", output)
	}
}

func TestComplianceExportInvalidFormat(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "audit.txt")

	root := newRootCommand()
	_, err := executeCommand(root, "compliance", "export", "--format", "xml", "--output", outputPath)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestComplianceExportRequiresOutput(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "compliance", "export")
	if err == nil {
		t.Fatal("expected error when --output is missing")
	}
}

func TestDistributedCommand(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "distributed", "--config", configPath, "--workers", "2")
	if err != nil {
		t.Fatalf("distributed failed: %v", err)
	}
	if !strings.Contains(output, "Distributed pipeline") {
		t.Fatalf("expected distributed output, got %q", output)
	}
}

func TestEventsListNoFile(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "events", "list", "--input", "/nonexistent/events.json")
	if err == nil {
		t.Fatal("expected error for nonexistent events file")
	}
}

func TestEventsReplayNoFile(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "events", "replay", "--input", "/nonexistent/events.json")
	if err == nil {
		t.Fatal("expected error for nonexistent events file")
	}
}

func TestEventsReplayEmptyFile(t *testing.T) {
	dir := t.TempDir()
	eventsFile := filepath.Join(dir, "events.json")
	if err := os.WriteFile(eventsFile, []byte("[]"), 0o644); err != nil {
		t.Fatalf("write events file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "events", "replay", "--input", eventsFile)
	if err != nil {
		t.Fatalf("events replay failed: %v", err)
	}
	if !strings.Contains(output, "No events to replay") {
		t.Fatalf("expected no events message, got %q", output)
	}
}

func TestEventsReplayWithEvents(t *testing.T) {
	dir := t.TempDir()
	eventsFile := filepath.Join(dir, "events.json")
	events := `[{"stream_id":"run-1","type":"pipeline.started","timestamp":"2024-01-01T00:00:00Z"},{"stream_id":"run-1","type":"pipeline.completed","timestamp":"2024-01-01T00:01:00Z"}]`
	if err := os.WriteFile(eventsFile, []byte(events), 0o644); err != nil {
		t.Fatalf("write events file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "events", "replay", "--input", eventsFile)
	if err != nil {
		t.Fatalf("events replay failed: %v", err)
	}
	if !strings.Contains(output, "pipeline") {
		t.Fatalf("expected replay output, got %q", output)
	}
}

func TestEventsReplayToFile(t *testing.T) {
	dir := t.TempDir()
	eventsFile := filepath.Join(dir, "events.json")
	outputFile := filepath.Join(dir, "snapshot.json")
	events := `[{"stream_id":"run-1","type":"pipeline.started","timestamp":"2024-01-01T00:00:00Z"}]`
	if err := os.WriteFile(eventsFile, []byte(events), 0o644); err != nil {
		t.Fatalf("write events file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "events", "replay", "--input", eventsFile, "--output", outputFile)
	if err != nil {
		t.Fatalf("events replay to file failed: %v", err)
	}
	if !strings.Contains(output, "Replayed") {
		t.Fatalf("expected replayed message, got %q", output)
	}
}

func TestEventsList(t *testing.T) {
	dir := t.TempDir()
	eventsFile := filepath.Join(dir, "events.json")
	events := `[{"stream_id":"run-1","type":"pipeline.started","timestamp":"2024-01-01T00:00:00Z"}]`
	if err := os.WriteFile(eventsFile, []byte(events), 0o644); err != nil {
		t.Fatalf("write events file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "events", "list", "--input", eventsFile)
	if err != nil {
		t.Fatalf("events list failed: %v", err)
	}
	if !strings.Contains(output, "pipeline.started") {
		t.Fatalf("expected event in list, got %q", output)
	}
}

func TestExportCompose(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("pipeline:\n  name: demo\n  mode: development\n  verbose: true\n  output_dir: ./out\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	outDir := filepath.Join(dir, "docker")

	root := newRootCommand()
	output, err := executeCommand(root, "export", "compose", "--config", configPath, "--input", "project: demo\nservices:\n  - name: api\n    port: 8080", "--output-dir", outDir)
	if err != nil {
		t.Fatalf("export compose failed: %v", err)
	}
	if !strings.Contains(output, "generated") {
		t.Fatalf("expected generated output, got %q", output)
	}
}

func TestExportComposeRequiresInput(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "export", "compose")
	if err == nil {
		t.Fatal("expected error when --input is missing")
	}
}

func TestImportHCL(t *testing.T) {
	dir := t.TempDir()
	hclFile := filepath.Join(dir, "spec.hcl")
	if err := os.WriteFile(hclFile, []byte(`project "myapp" {
  version = "1.0.0"
}

service "api" {
  port = 8080
  type = "backend"
}`), 0o644); err != nil {
		t.Fatalf("write HCL file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "import", "--input", hclFile)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if !strings.Contains(output, "myapp") {
		t.Fatalf("expected project name in output, got %q", output)
	}
}

func TestImportHCLToJSON(t *testing.T) {
	dir := t.TempDir()
	hclFile := filepath.Join(dir, "spec.hcl")
	if err := os.WriteFile(hclFile, []byte(`project "myapp" {
  version = "1.0.0"
}`), 0o644); err != nil {
		t.Fatalf("write HCL file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "import", "--input", hclFile, "--format", "json")
	if err != nil {
		t.Fatalf("import json failed: %v", err)
	}
	if !strings.Contains(output, `"name"`) && !strings.Contains(output, "myapp") {
		t.Fatalf("expected JSON output, got %q", output)
	}
}

func TestImportHCLToFile(t *testing.T) {
	dir := t.TempDir()
	hclFile := filepath.Join(dir, "spec.hcl")
	outputFile := filepath.Join(dir, "output.yaml")
	if err := os.WriteFile(hclFile, []byte(`project "myapp" { version = "1.0.0" }`), 0o644); err != nil {
		t.Fatalf("write HCL file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "import", "--input", hclFile, "--output", outputFile)
	if err != nil {
		t.Fatalf("import to file failed: %v", err)
	}
	if !strings.Contains(output, "Imported") {
		t.Fatalf("expected imported message, got %q", output)
	}
}

func TestImportRequiresInput(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "import")
	if err == nil {
		t.Fatal("expected error when --input is missing")
	}
}

func TestMigrationStatus(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "migration", "status")
	if err != nil {
		t.Fatalf("migration status failed: %v", err)
	}
	if !strings.Contains(output, "Migration Status") {
		t.Fatalf("expected migration status output, got %q", output)
	}
}

func TestMigrationHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "migration", "--help")
	if err != nil {
		t.Fatalf("migration --help failed: %v", err)
	}
	if !strings.Contains(output, "Database migration") {
		t.Fatalf("expected migration help, got %q", output)
	}
}

func TestSearchHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "search")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if !strings.Contains(output, "search") {
		t.Fatalf("expected search help, got %q", output)
	}
}

func TestSearchIndexAndQuery(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "search", "index", "--id", "doc1", "--title", "Hello World", "--content", "This is a test document")
	if err != nil {
		t.Fatalf("search index failed: %v", err)
	}

	output, err := executeCommand(root, "search", "query", "--term", "hello")
	if err != nil {
		t.Fatalf("search query failed: %v", err)
	}
	if !strings.Contains(output, "Found") {
		t.Fatalf("expected search results, got %q", output)
	}
}

func TestSearchCount(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "search", "count")
	if err != nil {
		t.Fatalf("search count failed: %v", err)
	}
	if !strings.Contains(output, "Documents in") {
		t.Fatalf("expected count output, got %q", output)
	}
}

func TestSearchList(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "search", "list")
	if err != nil {
		t.Fatalf("search list failed: %v", err)
	}
	if !strings.Contains(output, "Search indexes") {
		t.Fatalf("expected list output, got %q", output)
	}
}

func TestSearchQueryJSONOutput(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "search", "index", "--id", "doc1", "--title", "Test", "--content", "content")
	if err != nil {
		t.Fatalf("search index failed: %v", err)
	}

	output, err := executeCommand(root, "search", "query", "--term", "test", "--output", "json")
	if err != nil {
		t.Fatalf("search query json failed: %v", err)
	}
	if !strings.Contains(output, `"index"`) {
		t.Fatalf("expected JSON output, got %q", output)
	}
}

func TestSearchDelete(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "search", "index", "--id", "doc-to-delete", "--title", "Delete me", "--content", "content")
	if err != nil {
		t.Fatalf("search index failed: %v", err)
	}

	output, err := executeCommand(root, "search", "delete", "--id", "doc-to-delete")
	if err != nil {
		t.Fatalf("search delete failed: %v", err)
	}
	if !strings.Contains(output, "Deleted") {
		t.Fatalf("expected deleted message, got %q", output)
	}
}

func TestTestCommand(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	output, err := executeCommand(root, "test", "--dir", dir)
	if err != nil {
		t.Fatalf("test failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected test output")
	}
}

func TestTestCommandJSONOutput(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	output, err := executeCommand(root, "test", "--dir", dir, "--output", "json")
	if err != nil {
		t.Fatalf("test json failed: %v", err)
	}
	if !strings.Contains(output, `"status"`) {
		t.Fatalf("expected JSON test output, got %q", output)
	}
}

func TestWorkspaceInit(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	output, err := executeCommand(root, "workspace", "--root", dir, "init", "my-workspace")
	if err != nil {
		t.Fatalf("workspace init failed: %v", err)
	}
	if !strings.Contains(output, "Initialized workspace") {
		t.Fatalf("expected init message, got %q", output)
	}
}

func TestWorkspaceAddAndList(t *testing.T) {
	dir := t.TempDir()
	wsDir := filepath.Join(dir, "ws")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatalf("create ws dir: %v", err)
	}
	root := newRootCommand()
	_, err := executeCommand(root, "workspace", "--root", wsDir, "init", "my-workspace")
	if err != nil {
		t.Fatalf("workspace init failed: %v", err)
	}

	_, err = executeCommand(root, "workspace", "--root", wsDir, "add", "api-module", "./modules/api")
	if err != nil {
		t.Fatalf("workspace add failed: %v", err)
	}

	output, err := executeCommand(root, "workspace", "--root", wsDir, "list")
	if err != nil {
		t.Fatalf("workspace list failed: %v", err)
	}
	if !strings.Contains(output, "my-workspace") && !strings.Contains(output, "modules") {
		t.Fatalf("expected modules in list, got %q", output)
	}
}

func TestWorkspaceRemove(t *testing.T) {
	dir := t.TempDir()
	wsDir := filepath.Join(dir, "ws")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatalf("create ws dir: %v", err)
	}
	root := newRootCommand()
	_, err := executeCommand(root, "workspace", "--root", wsDir, "init", "my-workspace")
	if err != nil {
		t.Fatalf("workspace init failed: %v", err)
	}

	// AddModule creates the directory, so we remove by the directory name
	output, err := executeCommand(root, "workspace", "--root", wsDir, "remove", "my-workspace")
	if err != nil {
		t.Fatalf("workspace remove failed: %v", err)
	}
	if !strings.Contains(output, "Removed module") {
		t.Fatalf("expected remove message, got %q", output)
	}
}

func TestWorkspaceInfo(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "workspace", "--root", dir, "init", "my-workspace")
	if err != nil {
		t.Fatalf("workspace init failed: %v", err)
	}

	output, err := executeCommand(root, "workspace", "--root", dir, "info")
	if err != nil {
		t.Fatalf("workspace info failed: %v", err)
	}
	if !strings.Contains(output, "Workspace Information") {
		t.Fatalf("expected workspace info, got %q", output)
	}
}

func TestWorkspaceLock(t *testing.T) {
	dir := t.TempDir()
	root := newRootCommand()
	_, err := executeCommand(root, "workspace", "--root", dir, "init", "my-workspace")
	if err != nil {
		t.Fatalf("workspace init failed: %v", err)
	}

	output, err := executeCommand(root, "workspace", "--root", dir, "lock")
	if err != nil {
		t.Fatalf("workspace lock failed: %v", err)
	}
	if !strings.Contains(output, "Created") {
		t.Fatalf("expected lock message, got %q", output)
	}
}

func TestAuthHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "auth")
	if err != nil {
		t.Fatalf("auth failed: %v", err)
	}
	if !strings.Contains(output, "auth") {
		t.Fatalf("expected auth help, got %q", output)
	}
}

func TestAuthWhoamiNoKey(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "auth", "whoami")
	if err != nil {
		t.Fatalf("auth whoami failed: %v", err)
	}
	if !strings.Contains(output, "No API key provided") {
		t.Fatalf("expected no key message, got %q", output)
	}
}

func TestAuthWhoamiInvalidKey(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "auth", "whoami", "--api-key", "invalid-key")
	if err != nil {
		t.Fatalf("auth whoami invalid key failed: %v", err)
	}
	if !strings.Contains(output, "Authentication failed") {
		t.Fatalf("expected auth failed message, got %q", output)
	}
}

func TestAuthCreateUser(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "auth", "create-user", "--name", "testuser", "--email", "test@example.com", "--role", "admin")
	if err != nil {
		t.Fatalf("auth create-user failed: %v", err)
	}
	if !strings.Contains(output, "Created user") {
		t.Fatalf("expected created user message, got %q", output)
	}
}

func TestAuthListUsers(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "auth", "list-users")
	if err != nil {
		t.Fatalf("auth list-users failed: %v", err)
	}
	if !strings.Contains(output, "Users") && !strings.Contains(output, "user") && !strings.Contains(output, "No") {
		t.Fatalf("expected users list, got %q", output)
	}
}

func TestAuthListRoles(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "auth", "list-roles")
	if err != nil {
		t.Fatalf("auth list-roles failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatalf("expected roles output")
	}
}

func TestBrokerHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "broker")
	if err != nil {
		t.Fatalf("broker failed: %v", err)
	}
	if !strings.Contains(output, "broker") {
		t.Fatalf("expected broker help, got %q", output)
	}
}

func TestBrokerConnectMemory(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "broker", "connect", "--type", "memory", "--name", "test-mem-broker")
	if err != nil {
		t.Fatalf("broker connect memory failed: %v", err)
	}
	if !strings.Contains(output, "Connected to") {
		t.Fatalf("expected connected message, got %q", output)
	}
}

func TestBrokerConnectInvalidType(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "broker", "connect", "--type", "invalid", "--name", "fail-broker")
	if err == nil {
		t.Fatal("expected error for invalid broker type")
	}
}

func TestBrokerList(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "broker", "list")
	if err != nil {
		t.Fatalf("broker list failed: %v", err)
	}
	if !strings.Contains(output, "NAME") && !strings.Contains(output, "No broker") {
		t.Fatalf("expected broker list, got %q", output)
	}
}

func TestBrokerDisconnect(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "broker", "connect", "--type", "memory", "--name", "disc-broker")
	if err != nil {
		t.Fatalf("broker connect failed: %v", err)
	}

	output, err := executeCommand(root, "broker", "disconnect", "--name", "disc-broker")
	if err != nil {
		t.Fatalf("broker disconnect failed: %v", err)
	}
	if !strings.Contains(output, "Disconnected") {
		t.Fatalf("expected disconnect message, got %q", output)
	}
}

func TestDXHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx")
	if err != nil {
		t.Fatalf("dx failed: %v", err)
	}
	if !strings.Contains(output, "dx") {
		t.Fatalf("expected dx help, got %q", output)
	}
}

func TestDXVSCodeGen(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx", "vscode-gen")
	if err != nil {
		t.Fatalf("dx vscode-gen failed: %v", err)
	}
	if !strings.Contains(output, "package.json") {
		t.Fatalf("expected VS Code extension output, got %q", output)
	}
}

func TestDXCompletionBash(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx", "completion-bash")
	if err != nil {
		t.Fatalf("dx completion-bash failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected bash completion output")
	}
}

func TestDXCompletionZsh(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx", "completion-zsh")
	if err != nil {
		t.Fatalf("dx completion-zsh failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected zsh completion output")
	}
}

func TestDXCompletionPowerShell(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx", "completion-powershell")
	if err != nil {
		t.Fatalf("dx completion-powershell failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected powershell completion output")
	}
}

func TestDXSnippetList(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx", "snippet-list")
	if err != nil {
		t.Fatalf("dx snippet-list failed: %v", err)
	}
	if !strings.Contains(output, "SNIPPET") {
		t.Fatalf("expected snippet list, got %q", output)
	}
}

func TestDXSnippetGet(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dx", "snippet-get", "--name", "project")
	if err != nil {
		t.Fatalf("dx snippet-get failed: %v", err)
	}
	if len(strings.TrimSpace(output)) == 0 {
		t.Fatal("expected snippet content")
	}
}

func TestDXSnippetGetNotFound(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "dx", "snippet-get", "--name", "nonexistent-snippet")
	if err == nil {
		t.Fatal("expected error for nonexistent snippet")
	}
}

func TestCICDHelp(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "cicd", "--project", "test-app", "--languages", "go")
		if err != nil {
			t.Fatalf("cicd failed: %v", err)
		}
	})
	if len(out) == 0 {
		t.Fatal("expected CICD output")
	}
}

func TestCICDGitlab(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "cicd", "--platform", "gitlab", "--project", "test-app")
		if err != nil {
			t.Fatalf("cicd gitlab failed: %v", err)
		}
	})
	if len(out) == 0 {
		t.Fatal("expected CICD output")
	}
}

func TestCICDJenkins(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "cicd", "--platform", "jenkins", "--project", "test-app")
		if err != nil {
			t.Fatalf("cicd jenkins failed: %v", err)
		}
	})
	if !strings.Contains(out, "Pipeline") && !strings.Contains(out, "Jenkinsfile") {
		t.Fatalf("expected Jenkins output, got %q", out)
	}
}

func TestCICDWithInputFile(t *testing.T) {
	dir := t.TempDir()
	inputFile := filepath.Join(dir, "cicd.yaml")
	if err := os.WriteFile(inputFile, []byte("project: my-app\ngithub:\n  runner: ubuntu-latest\n"), 0o644); err != nil {
		t.Fatalf("write input file: %v", err)
	}

	root := newRootCommand()
	out := captureOutput(t, func() {
		_, err := executeCommand(root, "cicd", "--input-file", inputFile)
		if err != nil {
			t.Fatalf("cicd with input file failed: %v", err)
		}
	})
	if !strings.Contains(out, "my-app") && !strings.Contains(out, "CI/CD") {
		t.Fatalf("expected CICD output, got %q", out)
	}
}

func TestTemplateList(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "template", "list")
	if err != nil {
		t.Fatalf("template list failed: %v", err)
	}
	if !strings.Contains(output, "Templates") && !strings.Contains(output, "LLM") && !strings.Contains(output, "Compiler") {
		t.Fatalf("expected template list, got %q", output)
	}
}

func TestTemplateListCodeOnly(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "template", "list", "--kind", "code")
	if err != nil {
		t.Fatalf("template list code failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected template list output")
	}
}

func TestTemplateListPromptLLM(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "template", "list", "--kind", "prompt-llm")
	if err != nil {
		t.Fatalf("template list prompt-llm failed: %v", err)
	}
	if !strings.Contains(output, "LLM") {
		t.Fatalf("expected LLM prompts, got %q", output)
	}
}

func TestTemplateListPromptCompiler(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "template", "list", "--kind", "prompt-compiler")
	if err != nil {
		t.Fatalf("template list prompt-compiler failed: %v", err)
	}
	if !strings.Contains(output, "Compiler") {
		t.Fatalf("expected compiler templates, got %q", output)
	}
}

func TestTemplateShow(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "template", "show", "enrich-spec")
	if err != nil {
		t.Fatalf("template show failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected template content")
	}
}

func TestMarketplaceSearch(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "marketplace", "search", "web-api")
	if err != nil {
		t.Fatalf("marketplace search failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected marketplace search output")
	}
}

func TestMarketplaceSearchJSON(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "marketplace", "search", "web", "--output", "json")
	if err != nil {
		t.Fatalf("marketplace search json failed: %v", err)
	}
	if !strings.Contains(output, `"query"`) {
		t.Fatalf("expected JSON output, got %q", output)
	}
}

func TestPluginListEmpty(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "plugin", "list")
	if err != nil {
		t.Fatalf("plugin list failed: %v", err)
	}
	if !strings.Contains(output, "No plugins installed") {
		t.Fatalf("expected no plugins message, got %q", output)
	}
}

func TestAuditCommand(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	root := newRootCommand()
	output, err := executeCommand(root, "audit", testFile)
	if err != nil {
		t.Fatalf("audit failed: %v", err)
	}
	if !strings.Contains(output, "Audit complete") {
		t.Fatalf("expected audit output, got %q", output)
	}
}

func TestAuditNoFiles(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "audit")
	if err == nil {
		t.Fatal("expected error for no files")
	}
}

func TestMonitorHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "monitor", "--help")
	if err != nil {
		t.Fatalf("monitor --help failed: %v", err)
	}
	if !strings.Contains(output, "Prometheus") {
		t.Fatalf("expected monitor help, got %q", output)
	}
}

func TestGraphQLHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "graphql", "--help")
	if err != nil {
		t.Fatalf("graphql --help failed: %v", err)
	}
	if !strings.Contains(output, "GraphQL") {
		t.Fatalf("expected graphql help, got %q", output)
	}
}

func TestWSHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "ws", "--help")
	if err != nil {
		t.Fatalf("ws --help failed: %v", err)
	}
	if !strings.Contains(output, "WebSocket") {
		t.Fatalf("expected ws help, got %q", output)
	}
}

func TestAPIHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "api", "--help")
	if err != nil {
		t.Fatalf("api --help failed: %v", err)
	}
	if !strings.Contains(output, "REST API") {
		t.Fatalf("expected api help, got %q", output)
	}
}

func TestDashboardHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "dashboard", "--help")
	if err != nil {
		t.Fatalf("dashboard --help failed: %v", err)
	}
	if !strings.Contains(output, "dashboard") {
		t.Fatalf("expected dashboard help, got %q", output)
	}
}

func TestMCPHelp(t *testing.T) {
	root := newRootCommand()
	output, err := executeCommand(root, "mcp", "--help")
	if err != nil {
		t.Fatalf("mcp --help failed: %v", err)
	}
	if !strings.Contains(output, "Model Context Protocol") {
		t.Fatalf("expected mcp help, got %q", output)
	}
}

func TestResolveInputFileFromInputFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(f, []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := resolveInputFile("", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != f {
		t.Fatalf("expected %s, got %s", f, got)
	}
}

func TestResolveInputFileFromInput(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(f, []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := resolveInputFile(f, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != f {
		t.Fatalf("expected %s, got %s", f, got)
	}
}

func TestResolveInputFileNotFound(t *testing.T) {
	_, err := resolveInputFile("", "/nonexistent/file.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestResolveInputFileNeither(t *testing.T) {
	_, err := resolveInputFile("", "")
	if err == nil {
		t.Fatal("expected error when neither input nor input-file specified")
	}
}

func TestResolveInputFileInputNotFound(t *testing.T) {
	got, err := resolveInputFile("/nonexistent/file.yaml", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty string, got %s", got)
	}
}

func TestWriteToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := writeToFile(path, []byte("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(data))
	}
}

func TestRenderDoctorJSON(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		results := []checkResult{
			{Name: "check1", Status: "pass"},
			{Name: "check2", Status: "warn", Detail: "something"},
			{Name: "check3", Status: "fail", Detail: "broken"},
		}
		err := renderDoctorJSON(root, results)
		if err != nil {
			t.Fatalf("renderDoctorJSON failed: %v", err)
		}
	})
	if !strings.Contains(out, `"status":"unhealthy"`) {
		t.Fatalf("expected unhealthy status in output, got %q", out)
	}
	if !strings.Contains(out, `"passed":1`) {
		t.Fatalf("expected passed=1, got %q", out)
	}
	if !strings.Contains(out, `"warned":1`) {
		t.Fatalf("expected warned=1, got %q", out)
	}
	if !strings.Contains(out, `"failed":1`) {
		t.Fatalf("expected failed=1, got %q", out)
	}
}

func TestRenderDoctorJSONAllPass(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		results := []checkResult{
			{Name: "check1", Status: "pass"},
		}
		err := renderDoctorJSON(root, results)
		if err != nil {
			t.Fatalf("renderDoctorJSON failed: %v", err)
		}
	})
	if !strings.Contains(out, `"status":"healthy"`) {
		t.Fatalf("expected healthy status, got %q", out)
	}
}

func TestRenderDoctorJSONDegraded(t *testing.T) {
	root := newRootCommand()
	out := captureOutput(t, func() {
		results := []checkResult{
			{Name: "check1", Status: "pass"},
			{Name: "check2", Status: "warn"},
		}
		err := renderDoctorJSON(root, results)
		if err != nil {
			t.Fatalf("renderDoctorJSON failed: %v", err)
		}
	})
	if !strings.Contains(out, `"status":"degraded"`) {
		t.Fatalf("expected degraded status, got %q", out)
	}
}

func TestLoadCloudConfigFromSpec(t *testing.T) {
	dir := t.TempDir()
	spec := `cloud:
  provider: aws
  region: us-east-1
  project: myproject
  environment: staging
  resources:
    - name: bucket1
      type: s3
      spec:
        bucket: my-bucket
    - name: queue1
      kind: sqs
      spec:
        name: my-queue
`
	path := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(path, []byte(spec), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := loadCloudConfigFromSpec(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Provider != "aws" {
		t.Fatalf("expected provider 'aws', got %q", config.Provider)
	}
	if config.Region != "us-east-1" {
		t.Fatalf("expected region 'us-east-1', got %q", config.Region)
	}
	if config.Project != "myproject" {
		t.Fatalf("expected project 'myproject', got %q", config.Project)
	}
	if config.Environment != "staging" {
		t.Fatalf("expected environment 'staging', got %q", config.Environment)
	}
	if len(config.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(config.Resources))
	}
	if config.Resources[0].Type != "s3" {
		t.Fatalf("expected resource type 's3', got %q", config.Resources[0].Type)
	}
	if config.Resources[1].Type != "sqs" {
		t.Fatalf("expected resource type 'sqs' from kind, got %q", config.Resources[1].Type)
	}
}

func TestLoadCloudConfigFromSpecFileNotFound(t *testing.T) {
	_, err := loadCloudConfigFromSpec("/nonexistent/spec.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadCloudConfigFromSpecBadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("{{bad yaml"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := loadCloudConfigFromSpec(path)
	if err == nil {
		t.Fatal("expected error for bad yaml")
	}
}

func TestLoadCloudConfigFromSpecMissingProvider(t *testing.T) {
	dir := t.TempDir()
	spec := `cloud:
  region: us-east-1
`
	path := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(path, []byte(spec), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := loadCloudConfigFromSpec(path)
	if err == nil {
		t.Fatal("expected error for missing provider")
	}
}

func TestMustGetDriverFromStoreNotFound(t *testing.T) {
	driver := mustGetDriverFromStore("nonexistent-connection")
	if driver != "redis" {
		t.Fatalf("expected default 'redis', got %q", driver)
	}
}

func TestDeployDryRunUnknown(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "deploy", "--target", "unknown", "--dry-run")
	if err == nil {
		t.Fatal("expected error for unknown deploy target")
	}
}

func TestDeployDryRunLocal(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "test.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "local", "--dry-run")
	if err != nil {
		t.Fatalf("deploy local dry-run failed: %v", err)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Fatalf("expected dry-run output, got %q", out)
	}
}

func TestDeployMissingOutputDir(t *testing.T) {
	dir := t.TempDir()
	cfgContent := "pipeline:\n  name: test-app\n  output_dir: /nonexistent\n"
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	_, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "docker", "--dry-run")
	if err == nil {
		t.Fatal("expected error for nonexistent output dir")
	}
}

func TestDeployDockerDryRun(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "docker", "--dry-run")
	if err != nil {
		t.Fatalf("deploy docker dry-run failed: %v", err)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Fatalf("expected dry-run output, got %q", out)
	}
}

func TestDeployK8sDryRun(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "k8s", "--dry-run")
	if err != nil {
		t.Fatalf("deploy k8s dry-run failed: %v", err)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Fatalf("expected dry-run output, got %q", out)
	}
}

func TestDeployComposeDryRun(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "compose", "--dry-run")
	if err != nil {
		t.Fatalf("deploy compose dry-run failed: %v", err)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Fatalf("expected dry-run output, got %q", out)
	}
}

func TestDeploySSHDryRun(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "ssh", "--env", "staging", "--dry-run")
	if err != nil {
		t.Fatalf("deploy ssh dry-run failed: %v", err)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Fatalf("expected dry-run output, got %q", out)
	}
}

func TestDeployRsyncDryRun(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "rsync", "--env", "production", "--dry-run")
	if err != nil {
		t.Fatalf("deploy rsync dry-run failed: %v", err)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Fatalf("expected dry-run output, got %q", out)
	}
}

func TestDeployLocalNonDryRun(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "file.txt"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgContent := fmt.Sprintf("pipeline:\n  name: test-app\n  output_dir: %s\n", outputDir)
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	out, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "local")
	if err != nil {
		t.Fatalf("deploy local failed: %v", err)
	}
	if !strings.Contains(out, "copying") {
		t.Fatalf("expected copy output, got %q", out)
	}
}

func TestDeployDefaultOutputDir(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgContent := "pipeline:\n  name: test-app\n"
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	root := newRootCommand()
	_, err := executeCommand(root, "deploy", "--config", cfgFile, "--target", "docker", "--dry-run")
	if err != nil {
		t.Fatalf("deploy default output dir failed: %v", err)
	}
}

func TestMonitorStatusHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "monitor", "--help")
	if err != nil {
		t.Fatalf("monitor --help failed: %v", err)
	}
	if !strings.Contains(out, "monitor") {
		t.Fatalf("expected monitor help, got %q", out)
	}
}

func TestMonitorMetricsHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "monitor", "metrics", "--help")
	if err != nil {
		t.Fatalf("monitor metrics --help failed: %v", err)
	}
	if !strings.Contains(out, "metrics") {
		t.Fatalf("expected metrics help, got %q", out)
	}
}

func TestLintHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "lint", "--help")
	if err != nil {
		t.Fatalf("lint --help failed: %v", err)
	}
	if !strings.Contains(out, "lint") {
		t.Fatalf("expected lint help, got %q", out)
	}
}

func TestLintRun(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "lint", "--help")
	if err != nil {
		t.Fatalf("lint --help failed: %v", err)
	}
	if !strings.Contains(out, "lint") {
		t.Fatalf("expected lint help, got %q", out)
	}
}

func TestTUIRunHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "tui", "--help")
	if err != nil {
		t.Fatalf("tui --help failed: %v", err)
	}
	if !strings.Contains(out, "tui") {
		t.Fatalf("expected tui help, got %q", out)
	}
}

func TestProfileHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "profile", "--help")
	if err != nil {
		t.Fatalf("profile --help failed: %v", err)
	}
	if !strings.Contains(out, "profile") {
		t.Fatalf("expected profile help, got %q", out)
	}
}

func TestProfileList(t *testing.T) {
	root := newRootCommand()
	_, err := executeCommand(root, "profile", "list")
	if err != nil {
		t.Fatalf("profile list failed: %v", err)
	}
}

func TestCloudHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "cloud", "--help")
	if err != nil {
		t.Fatalf("cloud --help failed: %v", err)
	}
	if !strings.Contains(out, "cloud") {
		t.Fatalf("expected cloud help, got %q", out)
	}
}

func TestCloudPlanHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "cloud", "plan", "--help")
	if err != nil {
		t.Fatalf("cloud plan --help failed: %v", err)
	}
	if !strings.Contains(out, "plan") {
		t.Fatalf("expected plan help, got %q", out)
	}
}

func TestCloudPlanDryRun(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.yaml")
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
	if err := os.WriteFile(specPath, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	_, err := executeCommand(root, "cloud", "plan", "--input-file", specPath)
	if err != nil {
		t.Fatalf("cloud plan dry-run failed: %v", err)
	}
}

func TestObservabilityDashboardHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "observability", "dashboard", "--help")
	if err != nil {
		t.Fatalf("obs dashboard --help failed: %v", err)
	}
	if !strings.Contains(out, "dashboard") {
		t.Fatalf("expected dashboard help, got %q", out)
	}
}

func TestDoctorJSONOutput(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "doctor", "--output-format", "json", "--quick")
	if err != nil {
		t.Fatalf("doctor json failed: %v", err)
	}
	if !strings.Contains(out, "status") {
		t.Fatalf("expected JSON status output, got %q", out)
	}
}

func TestDoctorJSONAllFail(t *testing.T) {
	root := newRootCommand()
	results := []checkResult{
		{Name: "a", Status: "fail"},
		{Name: "b", Status: "fail"},
	}
	err := renderDoctorJSON(root, results)
	if err != nil {
		t.Fatalf("renderDoctorJSON failed: %v", err)
	}
}

func TestWorkflowHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "workflow", "--help")
	if err != nil {
		t.Fatalf("workflow --help failed: %v", err)
	}
	if !strings.Contains(out, "workflow") {
		t.Fatalf("expected workflow help, got %q", out)
	}
}

func TestWorkflowCreateHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "workflow", "create", "--help")
	if err != nil {
		t.Fatalf("workflow create --help failed: %v", err)
	}
	if !strings.Contains(out, "create") {
		t.Fatalf("expected create help, got %q", out)
	}
}

func TestWorkflowExecuteHelp(t *testing.T) {
	root := newRootCommand()
	out, err := executeCommand(root, "workflow", "execute", "--help")
	if err != nil {
		t.Fatalf("workflow execute --help failed: %v", err)
	}
	if !strings.Contains(out, "execute") {
		t.Fatalf("expected execute help, got %q", out)
	}
}
