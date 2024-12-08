package main

import (
	"fmt"
	"os"

	"github.com/deplotix/solar-boat-cli/cmd"
	"github.com/deplotix/solar-boat-cli/internal/version"
)

func main() {
	cmd.SetVersion(version.GetVersion())
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
