package runner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func loadEnv(explicitPath string) (envConfig, error) {
	path, err := resolveEnvPath(explicitPath)
	if err != nil {
		return envConfig{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return envConfig{}, fmt.Errorf("[runner] file not found: %s", path)
	}
	defer f.Close()

	cfg := envConfig{
		runtime: map[string]string{},
		ext:     map[string]string{},
		vars:    map[string]string{},
	}

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

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch {
		case strings.HasPrefix(key, "runtime."):
			cfg.runtime[strings.TrimPrefix(key, "runtime.")] = value
		case strings.HasPrefix(key, "ext."):
			cfg.ext[strings.TrimPrefix(key, "ext.")] = value
		case strings.HasPrefix(key, "var."):
			cfg.vars[key] = value
		}
	}

	if err := s.Err(); err != nil {
		return envConfig{}, err
	}
	return cfg, nil
}

func resolveEnvPath(explicitPath string) (string, error) {
	if explicitPath != "" {
		return explicitPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("[runner] file not found: runner.env")
	}

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("[runner] file not found: runner.env")
		}
		return filepath.Join(appData, "runner", "runner.env"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "runner", "runner.env"), nil
	default:
		return filepath.Join(home, ".config", "runner", "runner.env"), nil
	}
}
