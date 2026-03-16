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
		for _, rest := range args[i+1:] {
			if strings.HasPrefix(rest, "-") {
				return options{}, fmt.Errorf("[runner] unknown option: %s", rest)
			}
			return options{}, fmt.Errorf("[runner] unexpected argument: %s", rest)
		}
	}

	major := 0
	if opts.help {
		major++
	}
	if opts.version {
		major++
	}
	if opts.list {
		major++
	}
	if major > 1 {
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
		if strings.HasSuffix(e.Name(), ".run") {
			names = append(names, strings.TrimSuffix(e.Name(), ".run"))
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
	if filepath.Ext(target) != "" {
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
	cmd, err := splitCommand(cmdStr)
	if err != nil {
		return runPlan{}, err
	}
	cmd = append(cmd, target)
	return runPlan{Command: cmd}, nil
}

func buildRunPlanFromRun(path string, cfg envConfig, dryRun bool) (runPlan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return runPlan{}, err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.TrimPrefix(text, "\ufeff")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || lines[0] == "" {
		return runPlan{}, fmt.Errorf("[runner] invalid .run format\nmissing header")
	}
	header := lines[0]
	if !strings.HasPrefix(header, "#") {
		return runPlan{}, fmt.Errorf("[runner] invalid .run format\nmissing header")
	}
	if header == "#" {
		return runPlan{}, fmt.Errorf("[runner] invalid .run header")
	}
	body := ""
	if len(lines) > 1 {
		body = strings.Join(lines[1:], "\n")
	}

	runtimeName, tempName, err := resolveRunHeader(strings.TrimPrefix(header, "#"), cfg)
	if err != nil {
		return runPlan{}, err
	}

	cmdStr, ok := cfg.runtime[runtimeName]
	if !ok {
		return runPlan{}, fmt.Errorf("[runner] runtime not defined: %s", runtimeName)
	}
	cmd, err := splitCommand(cmdStr)
	if err != nil {
		return runPlan{}, err
	}
	tempPath := filepath.Join(os.TempDir(), tempName)
	cmd = append(cmd, tempPath)

	if dryRun {
		return runPlan{Command: cmd, TempPath: tempPath, UseTemp: true}, nil
	}

	if err := os.WriteFile(tempPath, []byte(body), 0o600); err != nil {
		return runPlan{}, err
	}
	return runPlan{Command: cmd, TempPath: tempPath, UseTemp: true}, nil
}

func resolveRunHeader(h string, cfg envConfig) (runtimeName, tempName string, err error) {
	if strings.HasPrefix(h, ".") {
		ext := strings.TrimPrefix(h, ".")
		if ext == "" {
			return "", "", fmt.Errorf("[runner] invalid .run header")
		}
		rn, ok := cfg.ext[ext]
		if !ok {
			return "", "", fmt.Errorf("[runner] extension not mapped: .%s", ext)
		}
		return rn, "runner_tmp." + ext, nil
	}

	if strings.Contains(h, ".") {
		parts := strings.Split(h, ".")
		ext := parts[len(parts)-1]
		if ext == "" {
			return "", "", fmt.Errorf("[runner] invalid .run header")
		}
		rn, ok := cfg.ext[ext]
		if !ok {
			return "", "", fmt.Errorf("[runner] extension not mapped: .%s", ext)
		}
		return rn, h, nil
	}

	return h, "runner_tmp", nil
}

func loadEnv(path string) (envConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return envConfig{}, fmt.Errorf("[runner] file not found: %s", path)
	}
	defer f.Close()

	cfg := envConfig{runtime: map[string]string{}, ext: map[string]string{}}
	s := bufio.NewScanner(f)
	first := true
	for s.Scan() {
		line := s.Text()
		if first {
			line = strings.TrimPrefix(line, "\ufeff")
			first = false
		}
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if strings.HasPrefix(key, "runtime.") {
			cfg.runtime[strings.TrimPrefix(key, "runtime.")] = val
		}
		if strings.HasPrefix(key, "ext.") {
			cfg.ext[strings.TrimPrefix(key, "ext.")] = val
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
		return nil, fmt.Errorf("[runner] invalid runtime command")
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("[runner] invalid runtime command")
	}
	return out, nil
}
