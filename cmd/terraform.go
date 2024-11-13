package cmd

import (
	"fmt"
	"os"

	"github.com/deplotix/solar-boat-cli/internal/terraform"
	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Manage Terraform operations",
	Long:  `Execute Terraform operations on changed modules and their dependencies`,
}

var planOutputDir string

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan on affected modules",
	Long:  `Detect changed Terraform modules and run terraform plan on them and their dependencies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modules, err := terraform.GetChangedModules(".")
		if err != nil {
			return fmt.Errorf("failed to get changed modules: %w", err)
		}

		if len(modules) == 0 {
			fmt.Println("No Terraform modules were changed")
			return nil
		}

		if planOutputDir != "" {
			// Create output directory if it doesn't exist
			if err := os.MkdirAll(planOutputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		fmt.Printf("Running terraform plan on %d modules\n", len(modules))
		return terraform.RunTerraformCommand(modules, "plan", planOutputDir)
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply on affected modules",
	Long:  `Detect changed Terraform modules and run terraform apply on them and their dependencies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modules, err := terraform.GetChangedModules(".")
		if err != nil {
			return fmt.Errorf("failed to get changed modules: %w", err)
		}

		if len(modules) == 0 {
			fmt.Println("No Terraform modules were changed")
			return nil
		}

		fmt.Printf("Running terraform apply on %d modules\n", len(modules))
		return terraform.RunTerraformCommand(modules, "apply")
	},
}

func init() {
	planCmd.Flags().StringVar(&planOutputDir, "output-dir", "", "Directory to save plan outputs")
	terraformCmd.AddCommand(planCmd)
	terraformCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(terraformCmd)
}
