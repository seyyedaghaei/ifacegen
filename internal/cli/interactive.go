package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// InteractiveConfig holds the configuration from interactive mode
type InteractiveConfig struct {
	MatchPatterns []string
	OutputFile    string
	NamePattern   string
	DryRun        bool
	Verbose       bool
	Watch         bool
	Progress      bool
}

// RunInteractive starts the interactive configuration mode
func RunInteractive() (*InteractiveConfig, error) {
	config := &InteractiveConfig{
		OutputFile:  "iface_gen.go",
		NamePattern: "I{}",
		Progress:    true,
	}

	fmt.Println("🎯 ifacegen Interactive Configuration")
	fmt.Println(strings.Repeat("=", 40))

	// Get match patterns
	fmt.Println("\n1️⃣  Match Patterns")
	fmt.Println("Enter comma-separated glob patterns to match struct names")
	fmt.Println("Examples: *Service, *Repository, *Manager")
	fmt.Print("Patterns: ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			config.MatchPatterns = parseCommaSeparated(input)
		}
	}

	// Get output file name
	fmt.Println("\n2️⃣  Output File")
	fmt.Println("Name of the generated interface file")
	fmt.Printf("Output file (default: %s): ", config.OutputFile)

	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			config.OutputFile = input
		}
	}

	// Get name pattern
	fmt.Println("\n3️⃣  Interface Naming")
	fmt.Println("Pattern for naming interfaces. Use {} as placeholder")
	fmt.Printf("Name pattern (default: %s): ", config.NamePattern)

	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			config.NamePattern = input
		}
	}

	// Get options
	fmt.Println("\n4️⃣  Options")

	// Dry run
	fmt.Print("Enable dry-run mode? (y/N): ")
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		config.DryRun = input == "y" || input == "yes"
	}

	// Verbose
	fmt.Print("Enable verbose output? (y/N): ")
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		config.Verbose = input == "y" || input == "yes"
	}

	// Watch mode
	fmt.Print("Enable watch mode? (y/N): ")
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		config.Watch = input == "y" || input == "yes"
	}

	// Progress
	fmt.Print("Show progress indicators? (Y/n): ")
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		config.Progress = input != "n" && input != "no"
	}

	// Show summary
	fmt.Println("\n📋 Configuration Summary")
	fmt.Println(strings.Repeat("=", 30))
	fmt.Printf("Match patterns: %v\n", config.MatchPatterns)
	fmt.Printf("Output file: %s\n", config.OutputFile)
	fmt.Printf("Name pattern: %s\n", config.NamePattern)
	fmt.Printf("Dry run: %t\n", config.DryRun)
	fmt.Printf("Verbose: %t\n", config.Verbose)
	fmt.Printf("Watch mode: %t\n", config.Watch)
	fmt.Printf("Progress: %t\n", config.Progress)

	// Confirm
	fmt.Print("\nProceed with this configuration? (Y/n): ")
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if input == "n" || input == "no" {
			fmt.Println("Configuration cancelled.")
			return nil, fmt.Errorf("user cancelled")
		}
	}

	return config, nil
}

// parseCommaSeparated parses comma-separated values
func parseCommaSeparated(input string) []string {
	parts := strings.Split(input, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GenerateCommand generates the command line equivalent
func (c *InteractiveConfig) GenerateCommand() string {
	var parts []string
	parts = append(parts, "ifacegen")

	if len(c.MatchPatterns) > 0 {
		parts = append(parts, "-match", strings.Join(c.MatchPatterns, ","))
	}

	if c.OutputFile != "iface_gen.go" {
		parts = append(parts, "-output", c.OutputFile)
	}

	if c.NamePattern != "I{}" {
		parts = append(parts, "-name", c.NamePattern)
	}

	if c.DryRun {
		parts = append(parts, "-dry-run")
	}

	if c.Verbose {
		parts = append(parts, "-verbose")
	}

	if c.Watch {
		parts = append(parts, "-watch")
	}

	if !c.Progress {
		parts = append(parts, "-progress=false")
	}

	parts = append(parts, "./...")

	return strings.Join(parts, " ")
}
