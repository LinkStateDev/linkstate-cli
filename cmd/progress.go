package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show your learning progress",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}

		var items []client.ProgressItem
		err := withSpinner("Loading progress…", func() error {
			var err error
			items, err = cliClient.GetProgress()
			return err
		})
		if err != nil {
			return fmt.Errorf("get progress: %w", err)
		}

		if len(items) == 0 {
			fmt.Println(ui.Muted.Render("No progress yet — pick a course on the web and run lst fetch <slug>."))
			return nil
		}

		// Group by course while preserving the (course, sort_order) ordering the
		// server returns.
		var (
			courseOrder []string
			byCourse    = map[string][]client.ProgressItem{}
			titleByID   = map[string]string{}
		)
		for _, p := range items {
			if _, ok := byCourse[p.CourseSlug]; !ok {
				courseOrder = append(courseOrder, p.CourseSlug)
				titleByID[p.CourseSlug] = p.CourseTitle
			}
			byCourse[p.CourseSlug] = append(byCourse[p.CourseSlug], p)
		}

		for i, slug := range courseOrder {
			if i > 0 {
				fmt.Println()
			}
			title := titleByID[slug]
			if title == "" {
				title = slug
			}
			fmt.Println(ui.Bold.Render(title) + " " + ui.Muted.Render("("+slug+")"))

			t := table.New().
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
				Headers("Lesson", "Status", "Completed").
				StyleFunc(func(row, col int) lipgloss.Style {
					base := lipgloss.NewStyle().Padding(0, 1)
					if row == table.HeaderRow {
						return base.Bold(true).Foreground(lipgloss.Color("8"))
					}
					return base
				})

			for _, p := range byCourse[slug] {
				name := p.LessonTitle
				if name == "" {
					name = p.LessonSlug
				}
				t.Row(name, statusCell(p.Status), completedCell(p.CompletedAt))
			}
			fmt.Println(t.Render())
		}
		return nil
	},
}

func statusCell(status string) string {
	switch status {
	case "completed":
		return ui.Success.Render(ui.GlyphPass + " completed")
	case "available":
		return ui.Hint.Render("○ available")
	default:
		return ui.Muted.Render(ui.GlyphLock + " " + status)
	}
}

func completedCell(t *string) string {
	if t == nil {
		return ui.Muted.Render("—")
	}
	return ui.Muted.Render(*t)
}

func init() {
	rootCmd.AddCommand(progressCmd)
}
