package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withDir(t *testing.T) {
	t.Helper()
	d := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(d); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestHelp(t *testing.T) {
	withDir(t)
	var out, err bytes.Buffer
	code := Main([]string{"--help"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "usage: runner [options] [target]") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestVersion(t *testing.T) {
	withDir(t)
	var out, err bytes.Buffer
	code := Main([]string{"--version"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if strings.TrimSpace(out.String()) != "runner 0.1.0" {
		t.Fatalf("out=%q", out.String())
	}
}

func TestListShowsRunWithoutExtension(t *testing.T) {
	withDir(t)
	write(t, "build.run", "#python\nprint(1)")
	if err := os.MkdirAll("sub", 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, "sub/deploy.run", "#python\nprint(1)")
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

func TestRunfileMissingWhenNoArgs(t *testing.T) {
	withDir(t)
	var out, err bytes.Buffer
	code := Main(nil, &out, &err)
	if code == 0 {
		t.Fatal("expected non-zero")
	}
	if !strings.Contains(err.String(), "runfile.run not found") {
		t.Fatalf("err=%q", err.String())
	}
}

func TestRunSourceFile(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "hello.go", "package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"OK\")}\n")
	var out, err bytes.Buffer
	code := Main([]string{"hello.go"}, &out, &err)
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

func TestNamedTaskRun(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "build.run", "#.go\npackage main\nimport \"fmt\"\nfunc main(){fmt.Println(\"TASK\")}\n")
	var out, err bytes.Buffer
	code := Main([]string{"build"}, &out, &err)
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
	code := Main([]string{"-n", "build"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	p := filepath.Join(os.TempDir(), "runner_tmp.go")
	if _, statErr := os.Stat(p); statErr == nil {
		t.Fatalf("temp file exists: %s", p)
	}
}

func TestMissingHeaderFormatError(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "bad.run", "print('x')\n")
	var out, err bytes.Buffer
	code := Main([]string{"bad.run"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "invalid .run format") || !strings.Contains(err.String(), "missing header") {
		t.Fatalf("err=%q", err.String())
	}
}

func TestInvalidHeaderError(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\next.go=go\n")
	write(t, "bad.run", "#\n")
	var out, err bytes.Buffer
	code := Main([]string{"bad.run"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "invalid .run header") {
		t.Fatalf("err=%q", err.String())
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
	if !strings.Contains(err.Error(), "invalid runtime command") {
		t.Fatalf("err=%v", err)
	}
}

func TestBOMAndCRLFRunFile(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.go=go run\r\next.go=go\r\n")
	content := "\ufeff#.go\r\npackage main\r\nimport \"fmt\"\r\nfunc main(){fmt.Println(\"BOM\")}\r\n"
	write(t, "runfile.run", content)
	var out, err bytes.Buffer
	code := Main(nil, &out, &err)
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

func TestExitCodePropagation(t *testing.T) {
	withDir(t)
	write(t, "runner.env", "runtime.sh=sh\next.sh=sh\n")
	write(t, "fail.sh", "exit 5\n")
	var out, err bytes.Buffer
	code := Main([]string{"fail.sh"}, &out, &err)
	if code != 5 {
		t.Fatalf("expected exit code 5, got %d, stderr=%s", code, err.String())
	}
}
