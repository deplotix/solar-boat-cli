package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "solarboat",
	Version: "0.1.1",
	Short:   "Solar Boat - A CLI tool for GitOps and Developer Experience",
	Long: `Solar Boat is a command-line interface tool designed for Infrastructure as Code (IaC) 
and GitOps workflows. It provides a wide range of Developer Experience (DX) capabilities including:

- Self-service ephemeral environments on Kubernetes
- Infrastructure management and deployment
- GitOps-based operations and workflows

Use "solarboat [command] --help" for more information about a command.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add version template
	rootCmd.SetVersionTemplate(`Solar Boat v{{.Version}}
`)
}
