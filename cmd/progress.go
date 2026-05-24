package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show your learning progress",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: lst auth")
		}

		items, err := cliClient.GetProgress()
		if err != nil {
			return fmt.Errorf("get progress: %w", err)
		}

		if len(items) == 0 {
			fmt.Println("No progress yet. Start a course!")
			fmt.Println("Run: lst courses")
			return nil
		}

		fmt.Println()
		fmt.Printf("%-12s %-12s %s\n", "Lesson ID", "Status", "Completed")
		fmt.Println("---------------------------------------------")
		for _, p := range items {
			icon := "🔒"
			switch p.Status {
			case "available":
				icon = color.Green("✅")
			case "completed":
				icon = color.Green("✔️")
			default:
				icon = color.Faint("🔒")
			}
			completed := "-"
			if p.CompletedAt != nil {
				completed = *p.CompletedAt
			}
			fmt.Printf("%-12d %s %-10s %s\n", p.LessonID, icon, p.Status, completed)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(progressCmd)
}
