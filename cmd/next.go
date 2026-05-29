package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/lesson"
	"github.com/spf13/cobra"
)

var nextForce bool

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Fetch the lesson after the current one",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}
		meta, err := lesson.LoadMeta()
		if err != nil {
			return err
		}

		var items []client.ProgressItem
		err = withSpinner("Finding next lesson…", func() error {
			var err error
			items, err = cliClient.GetProgress()
			return err
		})
		if err != nil {
			return fmt.Errorf("get progress: %w", err)
		}

		nextSlug, err := nextLessonAfter(items, meta.Slug)
		if err != nil {
			return err
		}
		dir, err := fetchAndPrepareLesson(nextSlug, nextForce)
		if err != nil {
			return err
		}
		return cdInto(dir)
	},
}

func nextLessonAfter(items []client.ProgressItem, currentSlug string) (string, error) {
	idx := -1
	for i, p := range items {
		if p.LessonSlug == currentSlug {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", errorWithHint(
			fmt.Sprintf("current lesson %q not found in your progress", currentSlug),
			"run: lst progress",
		)
	}
	trackSlug := items[idx].TrackSlug
	for i := idx + 1; i < len(items); i++ {
		if items[i].TrackSlug == trackSlug && items[i].LessonSlug != currentSlug {
			return items[i].LessonSlug, nil
		}
	}
	return "", errorWithHint(
		"no more lessons in this course",
		"run: lst resume to pick up another course",
	)
}

func init() {
	nextCmd.Flags().BoolVarP(&nextForce, "force", "f", false, "Overwrite existing lesson directory")
	rootCmd.AddCommand(nextCmd)
}
