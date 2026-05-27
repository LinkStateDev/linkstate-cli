package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Fetch the first unsolved lesson",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}

		var items []client.ProgressItem
		err := withSpinner("Finding next unsolved lesson…", func() error {
			var err error
			items, err = cliClient.GetProgress()
			return err
		})
		if err != nil {
			return fmt.Errorf("get progress: %w", err)
		}

		slug, err := firstUnsolved(items)
		if err != nil {
			return err
		}
		return fetchAndPrepareLesson(slug)
	},
}

// firstUnsolved returns the slug of the first lesson whose status is not
// "completed" in the server-provided ordering.
func firstUnsolved(items []client.ProgressItem) (string, error) {
	if len(items) == 0 {
		return "", errorWithHint(
			"no progress yet",
			"pick a course on the web and run lst fetch <slug>",
		)
	}
	for _, p := range items {
		if p.Status != "completed" {
			return p.LessonSlug, nil
		}
	}
	return "", fmt.Errorf("all lessons completed — congratulations")
}

func init() { rootCmd.AddCommand(resumeCmd) }
