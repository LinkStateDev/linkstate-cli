package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Run tests and submit result",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" { return fmt.Errorf("not logged in. Run: lst auth") }
		metaData, err := os.ReadFile(".linkstate.json")
		if err != nil { return fmt.Errorf(".linkstate.json not found. Run: lst fetch <slug>") }
		var meta struct {
			LessonID int    `json:"lesson_id"`
			Title    string `json:"title"`
		}
		if err := json.Unmarshal(metaData, &meta); err != nil { return fmt.Errorf("parse .linkstate.json: %w", err) }

		if _, err := os.Stat("main.go"); os.IsNotExist(err) {
			return fmt.Errorf("main.go not found in current directory")
		}
		testBin := "./test"
		if _, err := os.Stat(testBin); os.IsNotExist(err) {
			return fmt.Errorf("test binary not found. Run: lst fetch <slug>")
		}

		exec.Command("go", "build", "-o", "solution", "main.go").Run()
		defer os.Remove("solution")

		out, err := exec.Command(testBin).Output()
		exitOK := err == nil

		var results []testOutput
		if json.Unmarshal(out, &results) != nil {
			fmt.Println(string(out))
		} else {
			for _, r := range results {
				if r.Passed {
					fmt.Printf("  %s %s: %s\n", color.Green("✅"), r.Name, color.Green("PASS"))
				} else {
					fmt.Printf("  %s %s: %s\n", color.Red("❌"), r.Name, color.Red("FAIL"))
					if r.Expected != "" {
						fmt.Printf("     %s %s\n", color.Faint("expected:"), color.Yellow(r.Expected))
						fmt.Printf("     %s %s\n", color.Faint("actual:"), color.Yellow(r.Actual))
					}
					if r.Hint != "" {
						fmt.Printf("     %s %s\n", color.Bold("💡 Hint:"), r.Hint)
					}
				}
			}
		}

		status := "fail"
		if exitOK { status = "pass" }
		resp, err := cliClient.Submit(meta.LessonID, status)
		if err != nil { return fmt.Errorf("submit: %w", err) }
		fmt.Println()
		if resp.LessonCompleted {
			fmt.Printf("%s %s\n", color.Green("✅"), color.Bold("Lesson completed!"))
			if resp.NextLessonID != nil {
				fmt.Printf("Next lesson: %s\n", color.Yellow(fmt.Sprintf("%s/lessons/%d", cfg.Server, *resp.NextLessonID)))
			} else {
				fmt.Println("Course completed! 🎉")
			}
		} else {
			fmt.Printf("%s %s\n", color.Red("❌"), color.Faint("Not yet passed. Fix errors and run 'lst submit' again."))
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(submitCmd) }
