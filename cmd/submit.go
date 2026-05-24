package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/LinkStateDev/linkstate-cli/internal/taskrunner"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Run tests and submit result to server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: lst auth")
		}

		configData, err := os.ReadFile("test_config.json")
		if err != nil {
			return fmt.Errorf("not in a task directory? Run: lst fetch <id>")
		}

		metaData, err := os.ReadFile(".linkstate-task.json")
		if err != nil {
			return fmt.Errorf(".linkstate-task.json not found. Run: lst fetch <id>")
		}
		var meta struct {
			TaskID    int    `json:"task_id"`
			LessonID  int    `json:"lesson_id"`
			Title     string `json:"title"`
			TaskType  string `json:"task_type"`
		}
		if err := json.Unmarshal(metaData, &meta); err != nil {
			return fmt.Errorf("parse .linkstate-task.json: %w", err)
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
				fmt.Printf("  %s %s: %s\n", color.Green("✅"), r.Name, color.Green("PASS"))
			} else {
				fmt.Printf("  %s %s: %s\n", color.Red("❌"), r.Name, color.Red("FAIL"))
				fmt.Printf("     %s %s\n", color.Faint("expected:"), color.Yellow(r.Expected))
				fmt.Printf("     %s %s\n", color.Faint("actual:"), color.Yellow(r.Actual))
			}
		}

		status := "fail"
		if report.AllPass {
			status = "pass"
		}

		resp, err := cliClient.Submit(meta.TaskID, status)
		if err != nil {
			return fmt.Errorf("submit failed: %w", err)
		}

		fmt.Println()
		if resp.LessonCompleted {
			fmt.Printf("%s %s\n", color.Green("✅"), color.Bold("Task completed!"))
			if resp.NextLessonID != nil {
				fmt.Printf("Next lesson: %s\n", color.Yellow(fmt.Sprintf("%s/lessons/%d", cfg.Server, *resp.NextLessonID)))
			} else {
				fmt.Println("Course completed! 🎉")
			}
		} else {
			fmt.Printf("%s %s\n", color.Red("❌"), color.Faint("Task not yet passed. Fix errors and run 'lst submit' again."))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
