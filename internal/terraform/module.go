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

// GetChangedModules finds all Terraform modules and their dependencies
func GetChangedModules(rootDir string) ([]string, error) {
	modules := make(map[string]*Module)
	processedDirs := make(map[string]bool)

	// First pass: identify all modules and their direct dependencies
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

			// Create or get module
			if _, exists := modules[dir]; !exists {
				modules[dir] = &Module{
					Path:        dir,
					IsStateless: isStatelessModule(dir),
				}
			}

			// Parse file to find module dependencies
			content, err := os.ReadFile(path)
			if err == nil {
				dependencies := findModuleDependencies(string(content), rootDir)
				modules[dir].DependsOn = append(modules[dir].DependsOn, dependencies...)
			}
		}
		return nil
	})

	// Second pass: build reverse dependencies (UsedBy)
	for _, module := range modules {
		for _, dep := range module.DependsOn {
			if depModule, exists := modules[dep]; exists {
				depModule.UsedBy = append(depModule.UsedBy, module.Path)
			}
		}
	}

	// Get git changes
	changedFiles, err := getGitChangedFiles(rootDir)
	if err != nil {
		return nil, err
	}

	// Mark changed modules and their dependents
	changedModules := make([]string, 0)
	for _, changedFile := range changedFiles {
		dir := filepath.Dir(changedFile)
		if module, exists := modules[dir]; exists {
			markModuleChanged(module, modules, &changedModules)
		}
	}

	printModuleSummary(changedModules)
	return changedModules, nil
}

// markModuleChanged marks a module and all modules that depend on it as changed
func markModuleChanged(module *Module, allModules map[string]*Module, changedModules *[]string) {
	if module.Changed {
		return
	}

	module.Changed = true
	if !module.IsStateless {
		*changedModules = append(*changedModules, module.Path)
	}

	// Mark all modules that use this module as changed
	for _, userPath := range module.UsedBy {
		if userModule, exists := allModules[userPath]; exists {
			markModuleChanged(userModule, allModules, changedModules)
		}
	}
}

// findModuleDependencies parses terraform files to find module dependencies
func findModuleDependencies(content, rootDir string) []string {
	var dependencies []string
	// Simple regex or string matching to find module sources
	// This is a basic implementation - you might want to use HCL parser for more accuracy
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "source") && strings.Contains(line, "module") {
			// Extract module source path and resolve it to absolute path
			// This is a simplified example - you'll need more robust parsing
			sourcePath := extractModuleSource(line)
			if absPath, err := filepath.Abs(filepath.Join(rootDir, sourcePath)); err == nil {
				dependencies = append(dependencies, absPath)
			}
		}
	}
	return dependencies
}

// extractModuleSource extracts the source path from a module block
func extractModuleSource(line string) string {
	// This is a simplified implementation
	// You might want to use proper HCL parsing for more accuracy
	parts := strings.Split(line, "source")
	if len(parts) < 2 {
		return ""
	}
	source := strings.Trim(parts[1], " \t\"=")
	return source
}

// getGitChangedFiles returns a list of changed files from git
func getGitChangedFiles(rootDir string) ([]string, error) {
	// Get the merge base between HEAD and main branch
	mergeBaseCmd := exec.Command("git", "merge-base", "HEAD", "main")
	mergeBaseCmd.Dir = rootDir
	mergeBase, err := mergeBaseCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to find merge base: %w", err)
	}

	// Get changed files between merge base and HEAD
	diffCmd := exec.Command("git", "diff", "--name-only", strings.TrimSpace(string(mergeBase)), "HEAD")
	diffCmd.Dir = rootDir
	output, err := diffCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	// Convert output to slice of absolute file paths
	var changedFiles []string
	for _, file := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if file == "" {
			continue
		}
		absPath, err := filepath.Abs(filepath.Join(rootDir, file))
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for %s: %w", file, err)
		}
		changedFiles = append(changedFiles, absPath)
	}

	return changedFiles, nil
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
