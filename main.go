package main

import (
	"fmt"
	"os"
	"solarboat/cmd"
	"solarboat/internal/version"
)

func main() {
	cmd.SetVersion(version.GetVersion())
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
