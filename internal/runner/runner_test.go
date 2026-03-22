package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runMain(t *testing.T, args []string, out, errOut *bytes.Buffer) int {
	t.Helper()
	allArgs := append([]string{"-e", "runner.env"}, args...)
	return Main(allArgs, out, errOut)
}

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

func writeFile(t *testing.T, dir, name, body string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
func TestListShowsRunWithoutExtension(t *testing.T) {
	withDir(t)
	write(t, "build.run", "#python\nprint(1)")
	write(t, "runner.env", "runtime.python=python\next.py=python\n")
	var out, err bytes.Buffer
	code := Main([]string{"-e", "runner.env", "--list"}, &out, &err)
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
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "hello.go", "package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"OK\")}\n")
	var out, err bytes.Buffer
	code := Main([]string{"-e", "runner.env", "hello.go"}, &out, &err)
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
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "build.run", "#.go\npackage main\nimport \"fmt\"\nfunc main(){fmt.Println(\"TASK\")}\n")
	var out, err bytes.Buffer
	code := Main([]string{"-e", "runner.env", "build"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "TASK") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestDryRunDoesNotCreateTempFile(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "build.run", "#.go\npackage main\nfunc main(){}\n")
	var out, err bytes.Buffer
	code := Main([]string{"-e", "runner.env", "-n", "build"}, &out, &err)
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
	write(t, "runner.env", "runtime.go=go run\r\next.go=go\r\n")
	content := "\ufeff#.go\r\npackage main\r\nimport \"fmt\"\r\nfunc main(){fmt.Println(\"BOM\")}\r\n"
	write(t, "runfile.run", content)
	var out, err bytes.Buffer
	code := Main([]string{"-e", "runner.env"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "BOM") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestExtensionNotMapped(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\n")
	write(t, "hello.go", "package main\nfunc main(){}")
	var out, err bytes.Buffer
	code := Main([]string{"hello.go"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "extension not mapped: .go") {
		t.Fatalf("err=%q", err.String())
	}
}

func TestTempExtForRuntime(t *testing.T) {
	if got := tempExtForRuntime("pwsh"); got != ".ps1" {
		t.Fatalf("pwsh: got %q, want %q", got, ".ps1")
	}
	if got := tempExtForRuntime("bash"); got != "" {
		t.Fatalf("bash: got %q, want empty", got)
	}
}

func TestResolveNormalHeaderRuntime(t *testing.T) {
	cfg := envConfig{
		ext: map[string]string{
			"py": "python",
		},
	}

	runtimeName, tempExt, err := resolveNormalHeader(header{
		kind:        headerKindRuntime,
		runtimeName: "pwsh",
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runtimeName != "pwsh" {
		t.Fatalf("runtimeName: got %q, want %q", runtimeName, "pwsh")
	}
	if tempExt != ".ps1" {
		t.Fatalf("tempExt: got %q, want %q", tempExt, ".ps1")
	}
}

func TestResolveNormalHeaderExt(t *testing.T) {
	cfg := envConfig{
		ext: map[string]string{
			"py": "python",
		},
	}

	runtimeName, tempExt, err := resolveNormalHeader(header{
		kind:      headerKindExt,
		extension: "py",
	}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runtimeName != "python" {
		t.Fatalf("runtimeName: got %q, want %q", runtimeName, "python")
	}
	if tempExt != ".py" {
		t.Fatalf("tempExt: got %q, want %q", tempExt, ".py")
	}
}

func TestToMSYSPath(t *testing.T) {
	got := toMSYSPath(`C:\Users\jun\AppData\Local\Temp\runner-123`)
	want := `/c/Users/jun/AppData/Local/Temp/runner-123`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestExpandVars(t *testing.T) {
	vars := map[string]string{
		"name": "jun",
	}

	got, err := expandVars("hello ${var.name}", vars)
	if err != nil {
		t.Fatal(err)
	}

	want := "hello jun"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestExpandVarsUndefined(t *testing.T) {
	_, err := expandVars("hello ${var.xxx}", map[string]string{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckScriptAllOS(t *testing.T) {
	// env
	write(t, "runner.env", `
runtime.pwsh=pwsh
runtime.bash=bash
var.name=jun
`)

	// run
	write(t, "test.run", `#script

@windows
#pwsh
echo ${var.name}

@linux
#bash
echo ${var.name}
`)

	// 実行
	var out, errBuf bytes.Buffer
	code := Main([]string{"-e", "runner.env", "--check", "test.run"}, &out, &errBuf)

	if code != 0 {
		t.Fatalf("check failed: %s", errBuf.String())
	}
}
