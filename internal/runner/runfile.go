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
	if len(lines) == 0 {
		return runFile{}, fmt.Errorf("[runner] invalid .run format\nmissing header")
	}

	// 先頭の空行を除去
	var cleaned []string
	for _, line := range lines {
		if len(cleaned) == 0 && strings.TrimSpace(line) == "" {
			continue
		}
		cleaned = append(cleaned, line)
	}

	if len(cleaned) == 0 {
		return runFile{}, fmt.Errorf("[runner] invalid .run format\nmissing header")
	}

	first := strings.TrimSpace(cleaned[0])

	if strings.HasPrefix(first, "@") {
		script, err := parseScriptRunFile(cleaned)
		if err != nil {
			return runFile{}, err
		}
		return runFile{
			kind:   runFileKindScript,
			script: script,
		}, nil
	}

	h, err := parseHeader(cleaned[0], true)
	if err != nil {
		return runFile{}, err
	}

	bodyLines := []string{}
	if len(cleaned) > 1 {
		bodyLines = cleaned[1:]
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
	result := scriptRunFile{
		blocks: map[string]scriptBlock{},
	}

	var currentOS string
	var currentRuntime string
	var body []string

	flush := func() error {
		if currentOS == "" {
			return nil
		}
		if currentRuntime == "" {
			return fmt.Errorf("[runner] runtime not specified for os block: %s", currentOS)
		}
		result.blocks[currentOS] = scriptBlock{
			runtimeName: currentRuntime,
			body:        strings.Join(body, "\n"),
		}
		return nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "@") {
			if err := flush(); err != nil {
				return scriptRunFile{}, err
			}
			currentOS = strings.TrimPrefix(trimmed, "@")
			currentRuntime = ""
			body = nil
			continue
		}

		if currentRuntime == "" {
			if trimmed == "" {
				continue
			}
			if !strings.HasPrefix(trimmed, "#") {
				return scriptRunFile{}, fmt.Errorf("[runner] runtime header required after @%s", currentOS)
			}
			currentRuntime = strings.TrimPrefix(trimmed, "#")
			continue
		}

		body = append(body, line)
	}

	if err := flush(); err != nil {
		return scriptRunFile{}, err
	}

	return result, nil
}

func parseOSMarker(line string) (string, error) {
	switch strings.TrimPrefix(line, "@") {
	case "windows", "linux", "macos":
		return strings.TrimPrefix(line, "@"), nil
	default:
		return "", fmt.Errorf("[runner] unknown os block: %s", strings.TrimPrefix(line, "@"))
	}
}
