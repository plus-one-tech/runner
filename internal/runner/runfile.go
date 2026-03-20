package runner

import (
	"fmt"
	"strings"
)

type runFileKind string

const (
	runFileKindNormal runFileKind = "normal"
	runFileKindScript runFileKind = "script"
)

type headerKind string

const (
	headerKindRuntime  headerKind = "runtime"
	headerKindFilename headerKind = "filename"
	headerKindExt      headerKind = "ext"
)

type runFile struct {
	kind   runFileKind
	normal normalRunFile
	script scriptRunFile
}

type normalRunFile struct {
	header header
	body   string
}

type header struct {
	kind        headerKind
	runtimeName string
	filename    string
	extension   string
}

type scriptRunFile struct {
	blocks map[string]scriptBlock
}

type scriptBlock struct {
	osName      string
	runtimeName string
	body        string
}

func parseRunFile(text string) (runFile, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimPrefix(text, "\ufeff")

	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return runFile{}, fmt.Errorf("[runner] invalid .run format\nmissing header")
	}

	h, err := parseHeader(lines[0], true)
	if err != nil {
		return runFile{}, err
	}

	bodyLines := []string{}
	if len(lines) > 1 {
		bodyLines = lines[1:]
	}

	if h.kind == "" && h.runtimeName == "script" {
		script, err := parseScriptRunFile(bodyLines)
		if err != nil {
			return runFile{}, err
		}
		return runFile{
			kind:   runFileKindScript,
			script: script,
		}, nil
	}

	return runFile{
		kind: runFileKindNormal,
		normal: normalRunFile{
			header: h,
			body:   strings.Join(bodyLines, "\n"),
		},
	}, nil
}

func parseHeader(line string, allowScript bool) (header, error) {
	if !strings.HasPrefix(line, "#") || line == "#" {
		return header{}, fmt.Errorf("[runner] invalid .run header")
	}

	value := strings.TrimPrefix(line, "#")
	switch {
	case allowScript && value == "script":
		return header{runtimeName: "script"}, nil
	case strings.HasPrefix(value, "."):
		ext := strings.TrimPrefix(value, ".")
		if ext == "" {
			return header{}, fmt.Errorf("[runner] invalid .run header")
		}
		return header{kind: headerKindExt, extension: ext}, nil
	case strings.Contains(value, "."):
		parts := strings.Split(value, ".")
		ext := parts[len(parts)-1]
		if ext == "" {
			return header{}, fmt.Errorf("[runner] invalid .run header")
		}
		return header{kind: headerKindFilename, filename: value, extension: ext}, nil
	case value == "":
		return header{}, fmt.Errorf("[runner] invalid .run header")
	default:
		return header{kind: headerKindRuntime, runtimeName: value}, nil
	}
}

func parseScriptRunFile(lines []string) (scriptRunFile, error) {
	blocks := map[string]scriptBlock{}
	var currentOS string
	var currentRuntime string
	var body []string
	runtimeRequiredOS := ""

	flush := func() error {
		if currentOS == "" {
			return nil
		}
		if currentRuntime == "" {
			return fmt.Errorf("[runner] runtime header required in os block: %s", runtimeRequiredOS)
		}
		blocks[currentOS] = scriptBlock{
			osName:      currentOS,
			runtimeName: currentRuntime,
			body:        strings.Join(body, "\n"),
		}
		return nil
	}

	for _, raw := range lines {
		trimmed := strings.TrimSpace(raw)
		if currentOS == "" {
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if strings.HasPrefix(trimmed, "@") {
				osName, err := parseOSMarker(trimmed)
				if err != nil {
					return scriptRunFile{}, err
				}
				if _, exists := blocks[osName]; exists {
					return scriptRunFile{}, fmt.Errorf("[runner] duplicate os block: %s", osName)
				}
				currentOS = osName
				runtimeRequiredOS = osName
				currentRuntime = ""
				body = nil
				continue
			}
			return scriptRunFile{}, fmt.Errorf("[runner] invalid script block")
		}

		if currentRuntime == "" {
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, "@") {
				return scriptRunFile{}, fmt.Errorf("[runner] runtime header required in os block: %s", runtimeRequiredOS)
			}
			h, err := parseHeader(trimmed, false)
			if err != nil || h.kind != headerKindRuntime {
				return scriptRunFile{}, fmt.Errorf("[runner] runtime header required in os block: %s", runtimeRequiredOS)
			}
			currentRuntime = h.runtimeName
			continue
		}

		if strings.HasPrefix(trimmed, "@") {
			if err := flush(); err != nil {
				return scriptRunFile{}, err
			}
			osName, err := parseOSMarker(trimmed)
			if err != nil {
				return scriptRunFile{}, err
			}
			if _, exists := blocks[osName]; exists {
				return scriptRunFile{}, fmt.Errorf("[runner] duplicate os block: %s", osName)
			}
			currentOS = osName
			runtimeRequiredOS = osName
			currentRuntime = ""
			body = nil
			continue
		}

		body = append(body, raw)
	}

	if err := flush(); err != nil {
		return scriptRunFile{}, err
	}
	if len(blocks) == 0 {
		return scriptRunFile{}, nil
	}
	return scriptRunFile{blocks: blocks}, nil
}

func parseOSMarker(line string) (string, error) {
	switch strings.TrimPrefix(line, "@") {
	case "windows", "linux", "macos":
		return strings.TrimPrefix(line, "@"), nil
	default:
		return "", fmt.Errorf("[runner] unknown os block: %s", strings.TrimPrefix(line, "@"))
	}
}
