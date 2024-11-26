package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/deplotix/solar-boat-cli/internal/terraform"
	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Manage Terraform operations",
	Long:  `Execute Terraform operations on changed modules and their dependencies`,
}

var planOutputDir string = "terraform-plans"

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan on affected modules",
	Long:  `Detect changed Terraform modules and run terraform plan on them and their dependencies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printHeader("Terraform Plan")

		fmt.Println("🔍 Analyzing changes in Terraform modules...")
		modules, err := terraform.GetChangedModules(".")
		if err != nil {
			return fmt.Errorf("❌ Failed to get changed modules: %w", err)
		}

		if len(modules) == 0 {
			fmt.Println("\n✨ No changes detected:")
			fmt.Println("  • No stateful modules were changed")
			fmt.Println("  • No stateful modules were affected by changes in stateless modules")
			return nil
		}

		if err := os.MkdirAll(planOutputDir, 0755); err != nil {
			return fmt.Errorf("❌ Failed to create output directory: %w", err)
		}

		fmt.Printf("\n📋 Found %d affected module(s):\n", len(modules))
		for _, module := range modules {
			fmt.Printf("  • %s\n", shortenPath(module))
		}

		fmt.Printf("\n🚀 Starting Terraform operations...\n")
		if err := terraform.RunTerraformCommand(modules, "plan", planOutputDir); err != nil {
			fmt.Println("\n⚠️  Some operations failed. Check the errors above.")
			return err
		}

		fmt.Println("\n✅ All operations completed successfully!")
		return nil
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply on affected modules",
	Long:  `Detect changed Terraform modules and run terraform apply on them and their dependencies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printHeader("Terraform Apply")

		fmt.Println("🔍 Analyzing changes in Terraform modules...")
		modules, err := terraform.GetChangedModules(".")
		if err != nil {
			return fmt.Errorf("❌ Failed to get changed modules: %w", err)
		}

		if len(modules) == 0 {
			fmt.Println("\n✨ No changes detected:")
			fmt.Println("  • No stateful modules were changed")
			fmt.Println("  • No stateful modules were affected by changes in stateless modules")
			return nil
		}

		fmt.Printf("\n📋 Found %d affected module(s):\n", len(modules))
		for _, module := range modules {
			fmt.Printf("  • %s\n", shortenPath(module))
		}

		fmt.Printf("\n⚠️  About to apply changes to the above modules\n")
		fmt.Printf("🚀 Starting Terraform operations...\n")

		if err := terraform.RunTerraformCommand(modules, "apply", ""); err != nil {
			fmt.Println("\n⚠️  Some operations failed. Check the errors above.")
			return err
		}

		fmt.Println("\n✅ All operations completed successfully!")
		return nil
	},
}

func init() {
	planCmd.Flags().StringVar(&planOutputDir, "output-dir", "terraform-plans", "Directory to store plan files (default: terraform-plans)")
	terraformCmd.AddCommand(planCmd)
	terraformCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(terraformCmd)
}

// Helper functions for better output formatting
func printHeader(operation string) {
	fmt.Printf("\n%s\n%s\n", strings.Repeat("=", 50), operation)
	fmt.Printf("%s\n\n", strings.Repeat("=", 50))
}

func shortenPath(path string) string {
	// Get the last two components of the path for cleaner output
	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) <= 2 {
		return path
	}
	return "..." + string(os.PathSeparator) + strings.Join(parts[len(parts)-2:], string(os.PathSeparator))
}
