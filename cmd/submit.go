package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/lesson"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Run tests and submit result",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}
		meta, err := lesson.LoadMeta()
		if err != nil {
			return err
		}

		testErr := runTests(true)
		status := "fail"
		if testErr == nil {
			status = "pass"
		}

		var resp *client.SubmitResponse
		err = withSpinner("Submitting…", func() error {
			r, err := cliClient.Submit(meta.LessonID, status)
			if err != nil {
				return err
			}
			resp = r
			return nil
		})
		if err != nil {
			return fmt.Errorf("submit: %w", err)
		}

		fmt.Println()
		if resp.LessonCompleted {
			fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render("Lesson completed!"))
			if resp.NextLessonSlug != nil && *resp.NextLessonSlug != "" {
				fmt.Printf("Next lesson: %s\n", ui.Hint.Render(fmt.Sprintf("%s/lessons/%s", cfg.Server, *resp.NextLessonSlug)))
			} else if resp.NextLessonID != nil {
				fmt.Printf("Next lesson: %s\n", ui.Hint.Render(fmt.Sprintf("%s/lessons/%d", cfg.Server, *resp.NextLessonID)))
			} else {
				fmt.Println(ui.Bold.Render("Course completed! 🎉"))
			}
		} else {
			fmt.Printf("%s %s\n", ui.Error.Render(ui.GlyphFail), ui.Muted.Render("Not yet passed. Fix errors and run 'lst submit' again."))
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(submitCmd) }
