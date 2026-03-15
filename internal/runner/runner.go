package runner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const version = "0.1.0"

type options struct {
	dryRun  bool
	help    bool
	version bool
	list    bool
	target  string
}

type runPlan struct {
	Command  []string
	TempPath string
	UseTemp  bool
}

type envConfig struct {
	runtime map[string]string
	ext     map[string]string
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

	plan, err := buildRunPlan(opts)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	fmt.Fprintf(stdout, "[runner] command: %s\n", strings.Join(plan.Command, " "))
	if opts.dryRun {
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

func parseArgs(args []string) (options, error) {
	var opts options
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if opts.target != "" {
				return options{}, fmt.Errorf("[runner] unknown option: %s", arg)
			}
			switch arg {
			case "-n", "--dry-run":
				opts.dryRun = true
			case "-h", "--help":
				opts.help = true
			case "--version":
				opts.version = true
			case "--list":
				opts.list = true
			default:
				return options{}, fmt.Errorf("[runner] unknown option: %s", arg)
			}
			continue
		}
		if opts.target != "" {
			return options{}, fmt.Errorf("[runner] unexpected argument: %s", arg)
		}
		opts.target = arg
		if i != len(args)-1 {
			for _, rest := range args[i+1:] {
				if strings.HasPrefix(rest, "-") {
					return options{}, fmt.Errorf("[runner] unknown option: %s", rest)
				}
				return options{}, fmt.Errorf("[runner] unexpected argument: %s", rest)
			}
		}
	}

	flagCount := 0
	if opts.help {
		flagCount++
	}
	if opts.version {
		flagCount++
	}
	if opts.list {
		flagCount++
	}
	if flagCount > 1 {
		return options{}, fmt.Errorf("[runner] unknown option: conflicting options")
	}
	if (opts.help || opts.version || opts.list) && opts.target != "" {
		return options{}, fmt.Errorf("[runner] unexpected argument: %s", opts.target)
	}
	return opts, nil
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

func buildRunPlan(opts options) (runPlan, error) {
	target := opts.target
	if target == "" {
		target = "runfile.run"
		if _, err := os.Stat(target); err != nil {
			return runPlan{}, fmt.Errorf("[runner] runfile.run not found")
		}
	} else {
		target = resolveTarget(target)
		if _, err := os.Stat(target); err != nil {
			if opts.target == target {
				return runPlan{}, fmt.Errorf("[runner] file not found: %s", target)
			}
			return runPlan{}, fmt.Errorf("[runner] target not found: %s", target)
		}
	}

	cfg, err := loadEnv("runner.env")
	if err != nil {
		return runPlan{}, err
	}

	if strings.HasSuffix(target, ".run") {
		return buildRunPlanFromRun(target, cfg, opts.dryRun)
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

func buildRunPlanFromRun(path string, cfg envConfig, dryRun bool) (runPlan, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return runPlan{}, err
	}
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	text = strings.TrimPrefix(text, "\ufeff")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return runPlan{}, fmt.Errorf("[runner] invalid .run format\nmissing header")
	}
	header := lines[0]
	if !strings.HasPrefix(header, "#") || header == "#" {
		return runPlan{}, fmt.Errorf("[runner] invalid .run header")
	}
	body := ""
	if len(lines) > 1 {
		body = strings.Join(lines[1:], "\n")
	}

	var runtimeName, tempName string
	h := strings.TrimPrefix(header, "#")
	switch {
	case strings.HasPrefix(h, "."):
		ext := strings.TrimPrefix(h, ".")
		if ext == "" {
			return runPlan{}, fmt.Errorf("[runner] invalid .run header")
		}
		rn, ok := cfg.ext[ext]
		if !ok {
			return runPlan{}, fmt.Errorf("[runner] extension not mapped: .%s", ext)
		}
		runtimeName = rn
		tempName = "runner_tmp." + ext
	case strings.Contains(h, "."):
		parts := strings.Split(h, ".")
		ext := parts[len(parts)-1]
		if ext == "" {
			return runPlan{}, fmt.Errorf("[runner] invalid .run header")
		}
		rn, ok := cfg.ext[ext]
		if !ok {
			return runPlan{}, fmt.Errorf("[runner] extension not mapped: .%s", ext)
		}
		runtimeName = rn
		tempName = h
	default:
		runtimeName = h
		tempName = "runner_tmp"
	}

	cmdStr, ok := cfg.runtime[runtimeName]
	if !ok {
		return runPlan{}, fmt.Errorf("[runner] runtime not defined: %s", runtimeName)
	}
	args, err := splitCommand(cmdStr)
	if err != nil {
		return runPlan{}, err
	}
	tempPath := filepath.Join(os.TempDir(), tempName)
	args = append(args, tempPath)
	if dryRun {
		return runPlan{Command: args, TempPath: tempPath, UseTemp: true}, nil
	}
	if err := os.WriteFile(tempPath, []byte(body), 0o600); err != nil {
		return runPlan{}, err
	}
	return runPlan{Command: args, TempPath: tempPath, UseTemp: true}, nil
}

func loadEnv(path string) (envConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return envConfig{}, fmt.Errorf("[runner] file not found: %s", path)
	}
	defer f.Close()
	cfg := envConfig{runtime: map[string]string{}, ext: map[string]string{}}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(strings.TrimPrefix(s.Text(), "\ufeff"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if strings.HasPrefix(k, "runtime.") {
			cfg.runtime[strings.TrimPrefix(k, "runtime.")] = v
		}
		if strings.HasPrefix(k, "ext.") {
			cfg.ext[strings.TrimPrefix(k, "ext.")] = v
		}
	}
	if err := s.Err(); err != nil {
		return envConfig{}, err
	}
	return cfg, nil
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
		return nil, fmt.Errorf("[runner] invalid .run header")
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("[runner] runtime not defined: empty")
	}
	return out, nil
}
