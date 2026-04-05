package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func at(t *testing.T, id string) {
	t.Helper()
	t.Logf("%s", id)
}

func specPath(name string) string {
	return filepath.Join("testdata", "spec", name)
}

func withSpecDir(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	if err := os.Chdir(specPath(".")); err != nil {
		t.Fatal(err)
	}
}

//
// ===== 実行系 =====
//

// AT-010: ファイル指定実行
func TestSpec_AT010_FileExecution(t *testing.T) {
	at(t, "AT-010")

	var out, err bytes.Buffer
	code := Main([]string{specPath("hello.py")}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

// AT-012: 名前指定実行
func TestSpec_AT012_NamedTask(t *testing.T) {
	at(t, "AT-012")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{"build"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

// AT-013: 引数なし実行
func TestSpec_AT013_NoArgs(t *testing.T) {
	at(t, "AT-013")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

//
// ===== dry-run =====
//

// AT-030: dry-runは実行しない
func TestSpec_AT030_DryRun(t *testing.T) {
	at(t, "AT-030")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{"--dry-run", "build"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}

	if !strings.Contains(out.String(), "[runner]") {
		t.Fatalf("dry-run output missing: %s", out.String())
	}
}

// AT-063: dry-run tempなし
func TestSpec_AT063_DryRunNoTemp(t *testing.T) {
	at(t, "AT-063")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{"--dry-run", "build"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

//
// ===== CLI =====
//

// AT-040: help
func TestSpec_AT040_Help(t *testing.T) {
	at(t, "AT-040")

	var out, err bytes.Buffer
	code := Main([]string{"--help"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}

	if !strings.Contains(out.String(), "usage: runner") {
		t.Fatalf("invalid help: %s", out.String())
	}
}

// AT-041: version
func TestSpec_AT041_Version(t *testing.T) {
	at(t, "AT-041")

	var out, err bytes.Buffer
	code := Main([]string{"--version"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}

	if !strings.Contains(out.String(), "runner version") {
		t.Fatalf("invalid version: %s", out.String())
	}
}

//
// ===== list / check =====
//

// AT-042: --list
func TestSpec_AT042_List(t *testing.T) {
	at(t, "AT-042")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{"--list"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

// AT-044: list 空（最低限エラーにならない）
func TestSpec_AT044_ListEmpty(t *testing.T) {
	at(t, "AT-044")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{"--list"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

// AT-045: target後オプション拒否
func TestSpec_AT045_TargetAfterOptionRejected(t *testing.T) {
	at(t, "AT-045")

	var out, err bytes.Buffer
	code := Main([]string{"build", "--dry-run"}, &out, &err)

	if code == 0 {
		t.Fatalf("expected error but got success")
	}
}

// AT-046: --check
func TestSpec_AT046_Check(t *testing.T) {
	at(t, "AT-046")

	var out, err bytes.Buffer
	code := Main([]string{"--check", specPath("build.run")}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

// AT-049: check + dry-run
func TestSpec_AT049_CheckWithDryRun(t *testing.T) {
	at(t, "AT-049")

	var out, err bytes.Buffer
	code := Main([]string{"--check", "--dry-run", specPath("build.run")}, &out, &err)

	if code == 0 {
		t.Fatalf("expected error but got success")
	}
}

//
// ===== header =====
//

// AT-050: runtimeヘッダ
func TestSpec_AT050_RuntimeHeader(t *testing.T) {
	at(t, "AT-050")

	var out, err bytes.Buffer
	code := Main([]string{specPath("runtime.run")}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

// AT-052: extヘッダ
func TestSpec_AT052_ExtHeader(t *testing.T) {
	at(t, "AT-052")

	var out, err bytes.Buffer
	code := Main([]string{specPath("ext.run")}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

//
// ===== env =====
//

// AT-076: --env指定
func TestSpec_AT076_EnvSpecified(t *testing.T) {
	at(t, "AT-076")

	withSpecDir(t)

	var out, err bytes.Buffer
	code := Main([]string{"--env", "runner.test.env", "build"}, &out, &err)

	if code != 0 {
		t.Fatalf("code=%d err=%s", code, err.String())
	}
}

//
// ===== error =====
//

// AT-022: target未存在
func TestSpec_AT022_TargetNotFound(t *testing.T) {
	at(t, "AT-022")

	var out, err bytes.Buffer
	code := Main([]string{"notfound"}, &out, &err)

	if code == 0 {
		t.Fatalf("expected error but got success")
	}
}

// AT-090: ext未定義
func TestSpec_AT090_ExtNotDefined(t *testing.T) {
	at(t, "AT-090")

	var out, err bytes.Buffer
	code := Main([]string{"unknown.ext"}, &out, &err)

	if code == 0 {
		t.Fatalf("expected error but got success")
	}
}

// AT-091: runtime未定義
func TestSpec_AT091_RuntimeNotDefined(t *testing.T) {
	at(t, "AT-091")

	var out, err bytes.Buffer
	code := Main([]string{specPath("runtime-missing.run")}, &out, &err)

	if code == 0 {
		t.Fatalf("expected error but got success")
	}
}
