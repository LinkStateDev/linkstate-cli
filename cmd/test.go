package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/spf13/cobra"
)

type testOutput struct {
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
	Hint     string `json:"hint,omitempty"`
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests against your solution",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("main.go"); os.IsNotExist(err) {
			return fmt.Errorf("main.go not found in current directory")
		}
		testBin := "./test"
		if _, err := os.Stat(testBin); os.IsNotExist(err) {
			return fmt.Errorf("test binary not found. Run: lst fetch <slug>")
		}

		// Compile student's solution
		exec.Command("go", "build", "-o", "solution", "main.go").Run()
		defer os.Remove("solution")

		// Run test binary
		out, err := exec.Command(testBin).Output()
		exitOK := err == nil

		var results []testOutput
		if json.Unmarshal(out, &results) != nil {
			fmt.Println(string(out))
			if !exitOK { os.Exit(1) }
			return nil
		}

		passed, failed := 0, 0
		for _, r := range results {
			if r.Passed {
				fmt.Printf("  %s %s: %s\n", color.Green("✅"), r.Name, color.Green("PASS"))
				passed++
			} else {
				fmt.Printf("  %s %s: %s\n", color.Red("❌"), r.Name, color.Red("FAIL"))
				if r.Expected != "" {
					fmt.Printf("     %s %s\n", color.Faint("expected:"), color.Yellow(r.Expected))
					fmt.Printf("     %s %s\n", color.Faint("actual:"), color.Yellow(r.Actual))
				}
				if r.Hint != "" {
					fmt.Printf("     %s %s\n", color.Bold("💡 Hint:"), r.Hint)
				}
				failed++
			}
		}
		fmt.Println()
		if failed == 0 {
			fmt.Printf("All %d tests passed! Run: lst submit\n", passed)
		} else {
			fmt.Printf("%d passed, %d failed.\n", passed, failed)
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(testCmd) }
