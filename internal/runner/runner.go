package runner

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const version = "0.1.0"

type runPlan struct {
	Command  []string
	TempPath string
	UseTemp  bool
	Script   string
}

type envConfig struct {
	runtime map[string]string
	ext     map[string]string
	vars    map[string]string
}

func Main(args []string, stdout io.Writer, stderr io.Writer) int {
	opts, err := parseArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	if opts.help {
		printHelp(stdout)
		return 0
	}
	if opts.version {
		fmt.Fprintf(stdout, "runner %s\n", version)
		return 0
	}
	if opts.list {
		return runList(stdout, stderr)
	}

	if opts.check {
		if err := runCheck(opts); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}

	plan, err := buildRunPlan(opts)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	if plan.UseTemp {
		_ = os.WriteFile(plan.TempPath, []byte(plan.Script), 0o600)
	}

	fmt.Fprintf(stdout, "[runner] %s\n", strings.Join(plan.Command, " "))
	if opts.dryRun {

		if plan.UseTemp {
			data, err := os.ReadFile(plan.TempPath)
			if err == nil {
				fmt.Fprintln(stdout, "--- script ---")
				fmt.Fprint(stdout, string(data))
				if len(data) == 0 || data[len(data)-1] != '\n' {
					fmt.Fprintln(stdout)
				}
				fmt.Fprintln(stdout, "--- end ---")
			}
		}

		return 0
	}
	cmd := exec.Command(plan.Command[0], plan.Command[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if plan.UseTemp {
		_ = os.Remove(plan.TempPath)
	}
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	fmt.Fprintf(stderr, "[runner] %s\n", err.Error())
	return 1
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "usage: runner [options] [target]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "examples:")
	fmt.Fprintln(w, "  runner hello.py")
	fmt.Fprintln(w, "  runner build")
	fmt.Fprintln(w, "  runner")
}

func runList(stdout, stderr io.Writer) int {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(stderr, "[runner] %s\n", err.Error())
		return 1
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".run") {
			names = append(names, strings.TrimSuffix(name, ".run"))
		}
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(stdout, n)
	}
	return 0
}

func runCheck(opts options) error {
	target, err := resolveTargetPath(opts)
	if err != nil {
		return err
	}

	cfg, err := loadEnv(opts.envPath)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(target, ".run") {
		_, err := buildRunPlan(opts)
		return err
	}

	content, err := os.ReadFile(target)
	if err != nil {
		return err
	}

	rf, err := parseRunFile(string(content))
	if err != nil {
		return err
	}

	switch rf.kind {
	case runFileKindNormal:
		_, _, _, err := resolveRunFileTarget(rf, cfg, options{})
		return err

	case runFileKindScript:
		for osName, block := range rf.script.blocks {
			if _, ok := cfg.runtime[block.runtimeName]; !ok {
				return fmt.Errorf("[runner] runtime not defined: %s (os: %s)", block.runtimeName, osName)
			}

			if _, err := expandVars(block.body, cfg.vars); err != nil {
				return fmt.Errorf("%w (os: %s)", err, osName)
			}
		}
		return nil

	default:
		return fmt.Errorf("[runner] invalid .run format")
	}
}

func resolveTargetPath(opts options) (string, error) {
	target := opts.target
	if target == "" {
		target = "runfile.run"
		if _, err := os.Stat(target); err != nil {
			return "", fmt.Errorf("[runner] runfile.run not found")
		}
	} else {
		target = resolveTarget(target)
		if _, err := os.Stat(target); err != nil {
			if opts.target == target {
				return "", fmt.Errorf("[runner] file not found: %s", target)
			}
			return "", fmt.Errorf("[runner] target not found: %s", target)
		}
	}
	return target, nil
}

func buildRunPlan(opts options) (runPlan, error) {
	target, err := resolveTargetPath(opts)
	if err != nil {
		return runPlan{}, err
	}

	cfg, err := loadEnv(opts.envPath)
	if err != nil {
		return runPlan{}, err
	}

	if strings.HasSuffix(target, ".run") {
		return buildRunPlanFromRun(target, cfg, opts)
	}
	return buildRunPlanFromFile(target, cfg)
}

func resolveTarget(target string) string {
	if ext := filepath.Ext(target); ext != "" {
		return target
	}
	return target + ".run"
}

func buildRunPlanFromFile(target string, cfg envConfig) (runPlan, error) {
	ext := strings.TrimPrefix(filepath.Ext(target), ".")
	if ext == "" {
		return runPlan{}, fmt.Errorf("[runner] file not found: %s", target)
	}
	runtimeName, ok := cfg.ext[ext]
	if !ok {
		return runPlan{}, fmt.Errorf("[runner] extension not mapped: .%s", ext)
	}
	cmdStr, ok := cfg.runtime[runtimeName]
	if !ok {
		return runPlan{}, fmt.Errorf("[runner] runtime not defined: %s", runtimeName)
	}
	args, err := splitCommand(cmdStr)
	if err != nil {
		return runPlan{}, err
	}
	args = append(args, target)
	return runPlan{Command: args}, nil
}

func buildRunPlanFromRun(path string, cfg envConfig, opts options) (runPlan, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return runPlan{}, err
	}
	rf, err := parseRunFile(string(content))
	if err != nil {
		return runPlan{}, err
	}

	runtimeName, tempExt, body, err := resolveRunFileTarget(rf, cfg, opts)
	if err != nil {
		return runPlan{}, err
	}

	cmdStr, ok := cfg.runtime[runtimeName]
	if !ok {
		return runPlan{}, fmt.Errorf("[runner] runtime not defined: %s", runtimeName)
	}
	args, err := splitCommand(cmdStr)
	if err != nil {
		return runPlan{}, err
	}

	tempPath, err := makeTempPath(tempExt)
	if err != nil {
		return runPlan{}, err
	}

	scriptPath := runtimeScriptPath(runtimeName, tempPath, cfg)
	args = append(args, scriptPath)

	plan := runPlan{
		Command:  args,
		TempPath: tempPath,
		UseTemp:  true,
		Script:   body,
	}

	if opts.dryRun {
		return plan, nil
	}

	if err := os.WriteFile(tempPath, []byte(body), 0o600); err != nil {
		return runPlan{}, err
	}

	return plan, nil
}

func resolveRunFileTarget(rf runFile, cfg envConfig, opts options) (string, string, string, error) {
	switch rf.kind {
	case runFileKindNormal:
		runtimeName, tempExt, err := resolveNormalHeader(rf.normal.header, cfg)
		if err != nil {
			return "", "", "", err
		}

		body, err := expandVars(rf.normal.body, cfg.vars)
		if err != nil {
			return "", "", "", err
		}

		return runtimeName, tempExt, body, nil
	case runFileKindScript:
		osName := currentRunnerOS()
		if opts.dryRunOS != "" && opts.dryRunOS != "all" {
			osName = opts.dryRunOS
		}
		block, ok := rf.script.blocks[osName]
		if !ok {
			return "", "", "", fmt.Errorf("[runner] os block not found: %s", osName)
		}

		body, err := expandVars(block.body, cfg.vars)
		if err != nil {
			return "", "", "", err
		}

		return block.runtimeName, tempExtForRuntime(block.runtimeName), body, nil
	default:
		return "", "", "", fmt.Errorf("[runner] invalid .run header")
	}
}

func resolveNormalHeader(h header, cfg envConfig) (string, string, error) {
	switch h.kind {
	case headerKindRuntime:
		return h.runtimeName, tempExtForRuntime(h.runtimeName), nil
	case headerKindFilename, headerKindExt:
		runtimeName, ok := cfg.ext[h.extension]
		if !ok {
			return "", "", fmt.Errorf("[runner] extension not mapped: .%s", h.extension)
		}
		return runtimeName, "." + h.extension, nil
	default:
		return "", "", fmt.Errorf("[runner] invalid .run header")
	}
}

func tempExtForRuntime(runtimeName string) string {
	if runtimeName == "pwsh" {
		return ".ps1"
	}
	return ""
}

var varPattern = regexp.MustCompile(`\$\{var\.([a-zA-Z0-9_]+)\}`)

func expandVars(s string, vars map[string]string) (string, error) {
	var expandErr error

	out := varPattern.ReplaceAllStringFunc(s, func(m string) string {
		sub := varPattern.FindStringSubmatch(m)
		if len(sub) != 2 {
			expandErr = fmt.Errorf("[runner] invalid variable syntax: %s", m)
			return ""
		}

		key := sub[1]
		val, ok := vars[key]
		if !ok {
			expandErr = fmt.Errorf("[runner] undefined variable: %s", key)
			return ""
		}
		return val
	})

	if expandErr != nil {
		return "", expandErr
	}
	return out, nil
}

func runtimeScriptPath(runtimeName, path string, cfg envConfig) string {
	if currentRunnerOS() != "windows" {
		return path
	}

	switch runtimeName {
	case "bash", "sh", "zsh":
		cmd := cfg.runtime[runtimeName]
		return toWindowsShellPath(path, cmd)
	default:
		return path
	}
}

func toWindowsShellPath(path, cmd string) string {
	slash := filepath.ToSlash(path)

	if len(slash) >= 3 && slash[1] == ':' && slash[2] == '/' {
		drive := strings.ToLower(slash[:1])
		rest := slash[2:]

		lowerCmd := strings.ToLower(cmd)
		if strings.Contains(lowerCmd, "wsl") {
			return "/mnt/" + drive + rest
		}

		return "/" + drive + rest
	}

	return slash
}

func toMSYSPath(path string) string {
	slash := filepath.ToSlash(path)

	if len(slash) >= 3 && slash[1] == ':' && slash[2] == '/' {
		drive := slash[:1]
		rest := slash[2:]
		return "/" + strings.ToLower(drive) + rest
	}

	return slash
}

func currentRunnerOS() string {
	switch os.Getenv("RUNNER_TEST_OS_OVERRIDE") {
	case "windows", "linux", "macos":
		return os.Getenv("RUNNER_TEST_OS_OVERRIDE")
	}
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "darwin":
		return "macos"
	default:
		return "linux"
	}
}

func splitCommand(s string) ([]string, error) {
	var out []string
	var cur strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			if i+1 < len(s) {
				n := s[i+1]
				if n == '\\' || n == '"' {
					cur.WriteByte(n)
					i++
					continue
				}
			}
			cur.WriteByte(c)
		case '"':
			inQuote = !inQuote
		case ' ', '\t':
			if inQuote {
				cur.WriteByte(c)
				continue
			}
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(c)
		}
	}
	if inQuote {
		return nil, fmt.Errorf("[runner] invalid runtime command")
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("[runner] runtime not defined: empty")
	}
	return out, nil
}

func makeTempPath(ext string) (string, error) {
	suffix, err := randomHex(8)
	if err != nil {
		return "", err
	}
	name := "runner-" + suffix + ext
	return filepath.Join(os.TempDir(), name), nil
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
