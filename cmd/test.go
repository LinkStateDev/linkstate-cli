package cmd

import (
	"fmt"
	"os"

	"github.com/LinkStateDev/linkstate-cli/internal/taskrunner"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run local tests against your solution",
	RunE: func(cmd *cobra.Command, args []string) error {
		configData, err := os.ReadFile("test_config.json")
		if err != nil {
			return fmt.Errorf("not in a task directory? (test_config.json not found). Run: lst fetch <id>")
		}

		solutionFile := findSolutionFile()
		if solutionFile == "" {
			return fmt.Errorf("main.go not found in current directory")
		}

		report, err := taskrunner.Run(string(configData), solutionFile)
		if err != nil {
			return fmt.Errorf("test error: %w", err)
		}

		fmt.Println()
		for _, r := range report.Results {
			if r.Passed {
				fmt.Printf("  ✅ %s: PASS\n", r.Name)
			} else {
				fmt.Printf("  ❌ %s: FAIL\n", r.Name)
				fmt.Printf("     expected: %s\n", r.Expected)
				fmt.Printf("     actual:   %s\n", r.Actual)
			}
		}

		fmt.Println()
		if report.AllPass {
			fmt.Printf("All %d tests passed! Ready to submit.\n", report.Passed)
			fmt.Println("Run: lst submit")
		} else {
			fmt.Printf("%d passed, %d failed.\n", report.Passed, report.Failed)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
