package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Module represents a Terraform module with its properties
type Module struct {
	Path        string
	DependsOn   []string
	UsedBy      []string
	Changed     bool
	IsStateless bool
}

// GetChangedModules finds all Terraform modules in the given root directory
// and determines which ones are stateful vs stateless
func GetChangedModules(rootDir string) ([]string, error) {
	moduleMap := make(map[string]struct{})
	processedDirs := make(map[string]bool)

	// First check all .tf files in the current directory and subdirectories
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tf") {
			dir := filepath.Dir(path)
			if _, processed := processedDirs[dir]; processed {
				return nil
			}
			processedDirs[dir] = true

			fmt.Printf("🔍 Checking module: %s\n", dir)
			if !isStatelessModule(dir) {
				moduleMap[dir] = struct{}{}
				fmt.Printf("✅ Found stateful module: %s\n", dir)
			} else {
				fmt.Printf("📦 Skipping shared module: %s\n", dir)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	modules := make([]string, 0, len(moduleMap))
	for module := range moduleMap {
		modules = append(modules, module)
	}

	printModuleSummary(modules)
	return modules, nil
}

// RunTerraformCommand executes terraform commands for the given modules
func RunTerraformCommand(modules []string, command string, outputDir string) error {
	// Initialize all modules first
	if command == "plan" || command == "apply" {
		fmt.Println("🔧 Initializing all terraform modules...")
		for _, modulePath := range modules {
			fmt.Printf("  ⚡ Initializing %s\n", modulePath)
			initCmd := exec.Command("terraform", "init")
			initCmd.Dir = modulePath
			initCmd.Stdout = os.Stdout
			initCmd.Stderr = os.Stderr

			if err := initCmd.Run(); err != nil {
				return fmt.Errorf("failed to initialize terraform in %s: %w", modulePath, err)
			}
		}
		fmt.Println("✅ All modules initialized successfully")
	}

	// Run the actual command on all modules
	for _, modulePath := range modules {
		args := []string{command}

		if command == "plan" {
			if outputDir == "" {
				outputDir = "terraform-plans"
			}
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
			planFile := filepath.Join(outputDir, fmt.Sprintf("%s.tfplan", filepath.Base(modulePath)))
			args = append(args, "-out="+planFile)
		}

		fmt.Printf("🔨 Running terraform %s in %s\n", command, modulePath)
		cmd := exec.Command("terraform", args...)
		cmd.Dir = modulePath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run terraform %s in %s: %w", command, modulePath, err)
		}
	}
	return nil
}

// isStatelessModule determines if a module is stateless by checking for backend configuration
func isStatelessModule(path string) bool {
	files, err := filepath.Glob(filepath.Join(path, "*.tf"))
	if err != nil {
		// If there's an error reading the directory, assume it's not stateless
		return false
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// If we find a backend configuration, the module is not stateless
		if hasBackendConfig(string(content)) {
			return false
		}
	}

	return true
}

// hasBackendConfig checks if the given content contains a backend configuration
func hasBackendConfig(content string) bool {
	return strings.Contains(content, "terraform {") &&
		strings.Contains(content, "backend ")
}

// printModuleSummary prints a summary of found modules
func printModuleSummary(modules []string) {
	if len(modules) == 0 {
		fmt.Println("⚠️  No stateful modules were found")
		return
	}

	fmt.Printf("📝 Found %d stateful modules:\n", len(modules))
	for _, module := range modules {
		fmt.Printf("   - %s\n", module)
	}
}
