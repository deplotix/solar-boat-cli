package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Module struct {
	Path        string
	DependsOn   []string
	UsedBy      []string
	Changed     bool
	IsStateless bool
}

type ModuleError struct {
	Path    string
	Command string
	Error   error
}

// GetChangedModules returns a list of stateful module paths that need to be processed
func GetChangedModules(rootDir string) ([]string, error) {
	// Initialize module registry
	modules := make(map[string]*Module)

	// First pass: discover all modules and their stateless status
	if err := discoverModules(rootDir, modules); err != nil {
		return nil, fmt.Errorf("failed to discover modules: %w", err)
	}

	// Second pass: build dependency graph
	if err := buildDependencyGraph(modules); err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Get changed files from git
	changedFiles, err := getGitChangedFiles(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	// Process changes and return affected stateful modules
	return processChangedModules(changedFiles, modules)
}

func discoverModules(rootDir string, modules map[string]*Module) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		// Check if directory contains .tf files
		tfFiles, err := filepath.Glob(filepath.Join(path, "*.tf"))
		if err != nil || len(tfFiles) == 0 {
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}

		// Create module if it doesn't exist
		if _, exists := modules[absPath]; !exists {
			modules[absPath] = &Module{
				Path:        absPath,
				IsStateless: !hasBackendConfig(tfFiles),
			}
		}

		return nil
	})
}

func buildDependencyGraph(modules map[string]*Module) error {
	for path, module := range modules {
		// Find all .tf files in the module
		tfFiles, err := filepath.Glob(filepath.Join(path, "*.tf"))
		if err != nil {
			return fmt.Errorf("failed to glob .tf files in %s: %w", path, err)
		}

		// Parse each file for module dependencies
		for _, file := range tfFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			deps := findModuleDependencies(string(content), path)
			for _, dep := range deps {
				if depModule, exists := modules[dep]; exists {
					module.DependsOn = append(module.DependsOn, dep)
					depModule.UsedBy = append(depModule.UsedBy, path)
				}
			}
		}
	}
	return nil
}

func findModuleDependencies(content, currentDir string) []string {
	var deps []string
	lines := strings.Split(content, "\n")
	inModuleBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "module") && strings.Contains(trimmedLine, "{") {
			inModuleBlock = true
			continue
		}

		if inModuleBlock {
			if strings.Contains(trimmedLine, "source") {
				parts := strings.Split(trimmedLine, "=")
				if len(parts) == 2 {
					source := strings.Trim(parts[1], " \t\"'")
					modulePath := filepath.Clean(filepath.Join(currentDir, source))
					if absPath, err := filepath.Abs(modulePath); err == nil {
						deps = append(deps, absPath)
					}
				}
			}
			if strings.Contains(trimmedLine, "}") {
				inModuleBlock = false
			}
		}
	}
	return deps
}

func hasBackendConfig(tfFiles []string) bool {
	for _, file := range tfFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		if strings.Contains(string(content), "backend") {
			return true
		}
	}
	return false
}

func getGitChangedFiles(rootDir string) ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = rootDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	var changedFiles []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}

		file := strings.TrimSpace(line[2:])
		if strings.HasSuffix(file, ".tf") {
			absPath, err := filepath.Abs(filepath.Join(rootDir, file))
			if err == nil {
				changedFiles = append(changedFiles, absPath)
			}
		}
	}

	return changedFiles, nil
}

func processChangedModules(changedFiles []string, modules map[string]*Module) ([]string, error) {
	var affectedModules []string
	processed := make(map[string]bool)

	for _, file := range changedFiles {
		moduleDir := filepath.Dir(file)
		if module, exists := modules[moduleDir]; exists {
			markModuleChanged(module, modules, &affectedModules, processed)
		}
	}

	return affectedModules, nil
}

func markModuleChanged(module *Module, allModules map[string]*Module, affectedModules *[]string, processed map[string]bool) {
	if processed[module.Path] {
		return
	}
	processed[module.Path] = true

	module.Changed = true

	if module.IsStateless {
		// For stateless modules, propagate changes to dependent modules
		for _, userPath := range module.UsedBy {
			if userModule, exists := allModules[userPath]; exists {
				markModuleChanged(userModule, allModules, affectedModules, processed)
			}
		}
	} else {
		// For stateful modules, add to affected list
		*affectedModules = append(*affectedModules, module.Path)
	}
}

// RunTerraformCommand executes terraform commands on the specified modules
func RunTerraformCommand(modules []string, command string, planDir string) error {
	var failedModules []ModuleError

	for _, module := range modules {
		fmt.Printf("\n📦 Processing module: %s\n", module)

		// Run terraform init first
		fmt.Printf("  🔧 Initializing module...\n")
		initCmd := exec.Command("terraform", "init")
		initCmd.Dir = module
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr

		if err := initCmd.Run(); err != nil {
			fmt.Printf("  ❌ Initialization failed, skipping module\n")
			failedModules = append(failedModules, ModuleError{
				Path:    module,
				Command: "init",
				Error:   err,
			})
			continue
		}

		// Run the actual command (plan or apply)
		fmt.Printf("  🚀 Running terraform %s...\n", command)
		args := []string{command}
		if command == "plan" && planDir != "" {
			planFile := filepath.Join(planDir, filepath.Base(module)+".tfplan")
			args = append(args, "-out="+planFile)
		}

		cmd := exec.Command("terraform", args...)
		cmd.Dir = module
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			failedModules = append(failedModules, ModuleError{
				Path:    module,
				Command: command,
				Error:   err,
			})
			continue
		}

		fmt.Printf("  ✅ Module processed successfully\n")
	}

	// Report any failures at the end
	if len(failedModules) > 0 {
		fmt.Printf("\n⚠️  Some modules failed to process:\n")
		for _, failure := range failedModules {
			fmt.Printf("  ❌ %s: %s failed - %v\n",
				failure.Path,
				failure.Command,
				failure.Error)
		}
		return fmt.Errorf("failed to process %d module(s)", len(failedModules))
	}

	return nil
}
