package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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

		testErr := runTests(true)
		status := "fail"
		if testErr == nil { status = "pass" }

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
