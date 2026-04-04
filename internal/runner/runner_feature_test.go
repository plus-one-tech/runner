// runner_test.go

package runner

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runMain(t *testing.T, args []string, out, errOut *bytes.Buffer) int {
	t.Helper()
	allArgs := append([]string{"-e", "runner.test.env"}, args...)
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
func TestListShowsRunWithoutExtension(t *testing.T) {
	withDir(t)
	copyTestFile(t, "runner.test.env")
	copyTestFileAs(t, "list-build.run", "build.run")

	var out, err bytes.Buffer
	code := Main([]string{"-e", "runner.test.env", "--list"}, &out, &err)
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
func TestRunNamedTask(t *testing.T) {
	withDir(t)
	copyTestFile(t, "runner.test.env")
	copyTestFileAs(t, "build-task.run", "build.run")

	var out, err bytes.Buffer
	code := runMain(t, []string{"build"}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "TASK") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestDryRunDoesNotCreateTempFile(t *testing.T) {
	withDir(t)
	copyTestFile(t, "runner.test.env")

	copyTestFileAs(t, "build-dryrun.run", "build.run")
	var out, err bytes.Buffer
	code := runMain(t, []string{"-n", "build"}, &out, &err)
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
	copyTestFile(t, "runner.test.env")
	copyTestFile(t, "runfile.run")

	var out, err bytes.Buffer
	code := runMain(t, []string{}, &out, &err)
	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "runfile default") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestExtensionNotMapped(t *testing.T) {
	withDir(t)
	copyTestFile(t, "runner.test.env")
	copyTestFile(t, "hello.txt")

	var out, err bytes.Buffer
	code := runMain(t, []string{"hello.txt"}, &out, &err)
	if code == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.String(), "extension not mapped: .txt") {
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
	withDir(t)
	copyTestFile(t, "runner.test.env")
	copyTestFileAs(t, "test-script.run", "test.run")

	var out, errBuf bytes.Buffer
	code := runMain(t, []string{"--check", "test.run"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("check failed: %s", errBuf.String())
	}
}

func TestToWindowsShellPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows only")
	}

	path := `C:\Users\jun\AppData\Local\Temp\runner-123`

	// Git Bash 想定
	got := toWindowsShellPath(path, "bash")
	want := "/c/Users/jun/AppData/Local/Temp/runner-123"

	if got != want {
		t.Fatalf("bash: got %q, want %q", got, want)
	}

	// WSL 想定
	got = toWindowsShellPath(path, "wsl bash")
	want = "/mnt/c/Users/jun/AppData/Local/Temp/runner-123"

	if got != want {
		t.Fatalf("wsl: got %q, want %q", got, want)
	}
}

func TestDryRunPrintsCommand(t *testing.T) {
	withDir(t)
	copyTestFile(t, "runner.test.env")
	copyTestFileAs(t, "build-dryrun.run", "build.run")

	var out, err bytes.Buffer
	code := runMain(t, []string{"--dry-run", "build"}, &out, &err)

	if code != 0 {
		t.Fatalf("expected 0, got %d, err=%s", code, err.String())
	}
	if !strings.Contains(out.String(), "[runner] go run") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestDryRunDoesNotExecute(t *testing.T) {
	withDir(t)
	copyTestFile(t, "runner.test.env")
	copyTestFileAs(t, "dryrun-noexecute.run", "test.run")

	_ = os.Remove("test_output.txt")

	code := Main([]string{
		"-e", "runner.test.env",
		"--dry-run",
		"test.run",
	}, io.Discard, io.Discard)

	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}

	if _, err := os.Stat("test_output.txt"); err == nil {
		t.Fatal("file should not be created in dry-run")
	}
}

func copyTestFile(t *testing.T, name string) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	src := filepath.Join(filepath.Dir(thisFile), "..", "..", "test", name)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(name, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func copyTestFileAs(t *testing.T, srcName, dstName string) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	src := filepath.Join(filepath.Dir(thisFile), "..", "..", "test", srcName)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(dstName, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", dstName, err)
	}
}
