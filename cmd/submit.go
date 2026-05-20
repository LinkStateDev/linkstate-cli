package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LinkStateDev/linkstate-cli/internal/taskrunner"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Run tests and submit result to server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: linkstate-cli login")
		}

		configData, err := os.ReadFile("test_config.json")
		if err != nil {
			return fmt.Errorf("not in a task directory? Run: linkstate-cli fetch <id>")
		}

		metaData, err := os.ReadFile(".linkstate-task.json")
		if err != nil {
			return fmt.Errorf(".linkstate-task.json not found. Run: linkstate-cli fetch <id>")
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

		solutionFile := "solution.py"
		if _, err := os.Stat(solutionFile); os.IsNotExist(err) {
			return fmt.Errorf("solution.py not found in current directory")
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
			fmt.Printf("✅ Task %d completed!\n", meta.TaskID)
			if resp.NextLessonID != nil {
				fmt.Printf("Next lesson: %s/lessons/%d\n", cfg.Server, *resp.NextLessonID)
			} else {
				fmt.Println("Course completed! 🎉")
			}
		} else {
			fmt.Printf("❌ Task %d not yet passed. Fix errors and run 'linkstate-cli submit' again.\n", meta.TaskID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
