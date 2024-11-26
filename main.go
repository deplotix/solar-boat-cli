package main

import (
	"fmt"
	"os"

	"github.com/deplotix/solar-boat-cli/cmd"
)

// These variables will be set during build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
