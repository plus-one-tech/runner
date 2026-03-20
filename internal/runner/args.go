package runner

import (
	"fmt"
	"strings"
)

type options struct {
	dryRun   bool
	dryRunOS string
	check    bool
	help     bool
	version  bool
	list     bool
	envPath  string
	target   string
}

func parseArgs(args []string) (options, error) {
	var opts options
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if opts.target != "" {
				return options{}, unexpectedOptionAfterTarget(arg)
			}
			switch {
			case arg == "-n" || arg == "--dry-run":
				opts.dryRun = true
			case strings.HasPrefix(arg, "--dry-run="):
				osName := strings.TrimPrefix(arg, "--dry-run=")
				if err := validateDryRunOS(osName); err != nil {
					return options{}, err
				}
				opts.dryRun = true
				opts.dryRunOS = osName
			case arg == "-c" || arg == "--check":
				opts.check = true
			case arg == "-e" || arg == "--env":
				if i+1 >= len(args) {
					return options{}, fmt.Errorf("[runner] missing option value: %s", arg)
				}
				i++
				opts.envPath = args[i]
			case arg == "-h" || arg == "--help":
				opts.help = true
			case arg == "--version":
				opts.version = true
			case arg == "--list":
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
					return options{}, unexpectedOptionAfterTarget(rest)
				}
				return options{}, fmt.Errorf("[runner] unexpected argument: %s", rest)
			}
		}
	}

	if opts.check && opts.dryRun {
		return options{}, fmt.Errorf("[runner] invalid option combination")
	}

	modeCount := 0
	if opts.help {
		modeCount++
	}
	if opts.version {
		modeCount++
	}
	if opts.list {
		modeCount++
	}
	if modeCount > 1 {
		return options{}, fmt.Errorf("[runner] unknown option: conflicting options")
	}
	if (opts.help || opts.version || opts.list) && opts.target != "" {
		return options{}, fmt.Errorf("[runner] unexpected argument: %s", opts.target)
	}
	return opts, nil
}

func validateDryRunOS(osName string) error {
	switch osName {
	case "windows", "linux", "macos", "all":
		return nil
	default:
		return fmt.Errorf("[runner] unknown os: %s", osName)
	}
}

func unexpectedOptionAfterTarget(option string) error {
	return fmt.Errorf("[runner] unknown option: %s", option)
}
