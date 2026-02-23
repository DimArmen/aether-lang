package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aether-lang/aether/pkg/lexer"
	"github.com/aether-lang/aether/pkg/parser"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// Execute runs the CLI
func Execute(version string) error {
	rootCmd = &cobra.Command{
		Use:     "aether",
		Short:   "Aether - Next-generation Infrastructure as Code",
		Long:    "Aether - Infrastructure as Code with AI-powered intelligence",
		Version: version,
	}

	// Add commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(stateCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(fmtCmd)
	rootCmd.AddCommand(versionCmd(version))

	return rootCmd.Execute()
}

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new Aether project",
	Long:  "Create a new Aether project with starter configuration",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		return runInit(dir)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Aether configuration",
	Long:  "Check syntax and type correctness of Aether files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runValidate()
	},
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate an execution plan",
	Long:  "Show what changes Aether will make to your infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPlan()
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply [plan]",
	Short: "Apply infrastructure changes",
	Long:  "Create or update infrastructure according to configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApply()
	},
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy infrastructure",
	Long:  "Remove all resources defined in configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDestroy()
	},
}

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "State management commands",
	Long:  "Advanced state management operations",
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "AI agent management",
	Long:  "Manage and monitor AI agents",
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run infrastructure tests",
	Long:  "Execute unit and integration tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTest()
	},
}

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format Aether files",
	Long:  "Reformat Aether configuration files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFmt()
	},
}

func versionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Aether version %s\n", version)
		},
	}
}

// Command implementations

func runInit(dir string) error {
	color.Green("✓ Initializing Aether project in %s...", dir)
	color.Yellow("⚠  This is a placeholder - full implementation coming soon")
	return nil
}

func runValidate() error {
	color.Blue("→ Validating Aether configuration...")

	// Find all .ae files in current directory
	files, err := findAetherFiles(".")
	if err != nil {
		return fmt.Errorf("failed to find Aether files: %w", err)
	}

	if len(files) == 0 {
		color.Yellow("⚠  No .ae files found in current directory")
		return nil
	}

	hasErrors := false
	for _, file := range files {
		if err := validateFile(file); err != nil {
			color.Red("✗ %s: %v", file, err)
			hasErrors = true
		} else {
			color.Green("✓ %s", file)
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}

	color.Green("\n✓ All files validated successfully!")
	return nil
}

func findAetherFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".ae") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func validateFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors:\n  %s", strings.Join(p.Errors(), "\n  "))
	}

	// Basic validation: check that we have at least one statement
	if len(program.Statements) == 0 {
		return fmt.Errorf("empty program")
	}

	return nil
}

func runPlan() error {
	color.Blue("→ Generating execution plan...")
	color.Yellow("⚠  This is a placeholder - full implementation coming soon")
	return nil
}

func runApply() error {
	color.Blue("→ Applying infrastructure changes...")
	color.Yellow("⚠  This is a placeholder - full implementation coming soon")
	return nil
}

func runDestroy() error {
	color.Red("→ Destroying infrastructure...")
	color.Yellow("⚠  This is a placeholder - full implementation coming soon")
	return nil
}

func runTest() error {
	color.Blue("→ Running infrastructure tests...")
	color.Yellow("⚠  This is a placeholder - full implementation coming soon")
	return nil
}

func runFmt() error {
	color.Blue("→ Formatting Aether files...")
	color.Yellow("⚠  This is a placeholder - full implementation coming soon")
	return nil
}
