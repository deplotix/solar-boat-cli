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

var autoApprove bool

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan on affected modules",
	Long:  `Detect changed Terraform modules and run terraform plan on them and their dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		printHeader("Terraform Plan")

		fmt.Println("🔍 Analyzing changes in Terraform modules...")

		modules, err := terraform.GetChangedModules(".")
		if err != nil {
			fmt.Printf("❌ Failed to get changed modules: %v\n", err)
			os.Exit(1)
		}

		if len(modules) == 0 {
			fmt.Println("\n✨ No changes detected:")
			fmt.Println("  • No stateful modules were changed")
			fmt.Println("  • No stateful modules were affected by changes in stateless modules")
			return
		}

		if err := os.MkdirAll(planOutputDir, 0755); err != nil {
			fmt.Printf("❌ Failed to create output directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n📋 Found %d affected module(s):\n", len(modules))
		for _, module := range modules {
			fmt.Printf("  • %s\n", shortenPath(module))
		}

		fmt.Printf("\n🚀 Starting Terraform operations...\n")
		if err := terraform.RunTerraformCommand(modules, "plan", planOutputDir); err != nil {
			fmt.Println("\n⚠️  Some operations failed. Check the errors above.")
			os.Exit(1)
		}

		fmt.Println("\n✅ All operations completed successfully!")
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply on affected modules",
	Long:  `Detect changed Terraform modules and run terraform apply on them and their dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		printHeader("Terraform Apply")

		fmt.Println("🔍 Analyzing changes in Terraform modules...")
		modules, err := terraform.GetChangedModules(".")
		if err != nil {
			fmt.Printf("❌ Failed to get changed modules: %v\n", err)
			os.Exit(1)
		}

		if len(modules) == 0 {
			fmt.Println("\n✨ No changes detected:")
			fmt.Println("  • No stateful modules were changed")
			fmt.Println("  • No stateful modules were affected by changes in stateless modules")
			return
		}

		fmt.Printf("\n📋 Found %d affected module(s):\n", len(modules))
		for _, module := range modules {
			fmt.Printf("  • %s\n", shortenPath(module))
		}

		// Add confirmation unless auto-approve is set
		if !autoApprove {
			fmt.Print("\n⚠️  Do you want to apply these changes? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("❌ Apply cancelled")
				os.Exit(0)
			}
		}

		fmt.Printf("\n🚀 Starting Terraform operations...\n")
		if err := terraform.RunTerraformCommand(modules, "apply", ""); err != nil {
			fmt.Println("\n⚠️  Some operations failed. Check the errors above.")
			os.Exit(1)
		}

		fmt.Println("\n✅ All operations completed successfully!")
	},
}

func init() {
	planCmd.Flags().StringVar(&planOutputDir, "output-dir", "terraform-plans", "Directory to store plan files (default: terraform-plans)")
	applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Skip interactive approval before applying")
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
