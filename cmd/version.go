package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the version number
	Version = "0.1.2"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of the CLI tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Solar Boat CLI v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
