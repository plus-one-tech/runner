package main

import (
	"os"

	"runner/internal/runner"
)

func main() {
	os.Exit(runner.Main(os.Args[1:], os.Stdout, os.Stderr))
}
