package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withDir(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(d); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	return d
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "runner.env")
	write(t, path, content)
	return path
}

func TestListShowsRunWithoutExtension(t *testing.T) {
	withDir(t)
	write(t, "build.run", "#python\nprint(1)")
	write(t, "runner.env", "runtime.python=python\next.py=python\n")
	var out, err bytes.Buffer
	code := Main([]string{"--list"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if strings.TrimSpace(out.String()) != "build" {
		t.Fatalf("out=%q", out.String())
	}
}

func TestTargetOptionAfterIsRejected(t *testing.T) {
	withDir(t)
	var out, err bytes.Buffer
	code := Main([]string{"hello.py", "-n"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "unknown option: -n") {
		t.Fatalf("err=%q", err.String())
	}
}

func TestRunPythonFile(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\next.go=go\n")
	write(t, "hello.go", "package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"OK\")}\n")
	var out, err bytes.Buffer
	code := Main([]string{"--env", envPath, "hello.go"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "[runner] command: go run hello.go") {
		t.Fatalf("out=%q", out.String())
	}
	if !strings.Contains(out.String(), "OK") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestRunNamedTask(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\next.go=go\n")
	write(t, "build.run", "#.go\npackage main\nimport \"fmt\"\nfunc main(){fmt.Println(\"TASK\")}\n")
	var out, err bytes.Buffer
	code := Main([]string{"--env", envPath, "build"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "TASK") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestDryRunDoesNotCreateTempFile(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\next.go=go\n")
	write(t, "build.run", "#.go\npackage main\nfunc main(){}\n")
	var out, err bytes.Buffer
	code := Main([]string{"-n", "--env", envPath, "build"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if strings.Contains(out.String(), "runner_tmp.go") {
		t.Fatalf("dry-run should use unique temp file path: %q", out.String())
	}
}

func TestCommandQuoteSplit(t *testing.T) {
	args, err := splitCommand("\"C:\\\\Program Files\\\\Python\\\\python.exe\" -u")
	if err != nil {
		t.Fatal(err)
	}
	if len(args) != 2 || args[0] != "C:\\Program Files\\Python\\python.exe" || args[1] != "-u" {
		t.Fatalf("args=%v", args)
	}
}

func TestCommandInvalidQuote(t *testing.T) {
	_, err := splitCommand("\"python -u")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBOMAndCRLFRunFile(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\r\next.go=go\r\n")
	content := "\ufeff#.go\r\npackage main\r\nimport \"fmt\"\r\nfunc main(){fmt.Println(\"BOM\")}\r\n"
	write(t, "runfile.run", content)
	var out, err bytes.Buffer
	code := Main([]string{"--env", envPath}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "BOM") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestExtensionNotMapped(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\n")
	write(t, "hello.go", "package main\nfunc main(){}")
	var out, err bytes.Buffer
	code := Main([]string{"--env", envPath, "hello.go"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "extension not mapped: .go") {
		t.Fatalf("err=%q", err.String())
	}
}

func TestEnvOptionUsesExplicitFile(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\next.go=go\n")
	write(t, "hello.go", "package main\nfunc main(){}")
	var out, err bytes.Buffer
	code := Main([]string{"--env", envPath, "-n", "hello.go"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "go run hello.go") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestCheckReturnsSuccessWithoutOutput(t *testing.T) {
	withDir(t)
	envPath := writeEnvFile(t, "runtime.go=go run\next.go=go\n")
	write(t, "build.run", "#.go\npackage main\nfunc main(){}\n")
	var out, err bytes.Buffer
	code := Main([]string{"--env", envPath, "--check", "build"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if out.Len() != 0 {
		t.Fatalf("out=%q", out.String())
	}
	if err.Len() != 0 {
		t.Fatalf("err=%q", err.String())
	}
}

func TestCheckAndDryRunCombinationIsRejected(t *testing.T) {
	withDir(t)
	var out, err bytes.Buffer
	code := Main([]string{"--check", "--dry-run", "build"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "invalid option combination") {
		t.Fatalf("err=%q", err.String())
	}
}

func TestDryRunOSOptionIsAccepted(t *testing.T) {
	opts, err := parseArgs([]string{"--dry-run=linux", "build"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.dryRun || opts.dryRunOS != "linux" {
		t.Fatalf("opts=%+v", opts)
	}
}

func TestDryRunUnknownOSIsRejected(t *testing.T) {
	_, err := parseArgs([]string{"--dry-run=freebsd", "build"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown os: freebsd") {
		t.Fatalf("err=%q", err.Error())
	}
}

func TestResolveEnvPathDoesNotUseCurrentDirectory(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	path, err := resolveEnvPath("")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "runner.env" {
		t.Fatalf("path=%q", path)
	}
	if path == "runner.env" || path == filepath.Join(".", "runner.env") {
		t.Fatalf("unexpected current-directory env path: %q", path)
	}
}
