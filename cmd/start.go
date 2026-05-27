package cmd

import (
	"fmt"
	"strings"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <course-slug>[/<module-slug>]",
	Short: "Start a course (or specific module) at the first unsolved lesson",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}
		arg := args[0]

		if strings.Contains(arg, "/") {
			return startModule(arg)
		}
		return startCourse(arg)
	},
}

func startModule(arg string) error {
	parts := strings.SplitN(arg, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errorWithHint(
			fmt.Sprintf("invalid argument %q", arg),
			"use <course-slug>/<module-slug>, e.g. lst start netops-framework/state-tester",
		)
	}
	courseSlug, moduleSlug := parts[0], parts[1]

	var resp *client.StartResponse
	err := withSpinner("Looking up module…", func() error {
		var err error
		resp, err = cliClient.GetStartLesson(courseSlug, moduleSlug)
		return err
	})
	if err != nil {
		return fmt.Errorf("start module: %w", err)
	}

	if resp.ModuleCompleted {
		fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render("Module already completed — opening last lesson for review"))
	}
	return fetchAndPrepareLesson(resp.Lesson.Slug)
}

func startCourse(course string) error {
	var items []client.ProgressItem
	err := withSpinner("Looking up course…", func() error {
		var err error
		items, err = cliClient.GetProgress()
		return err
	})
	if err != nil {
		return fmt.Errorf("get progress: %w", err)
	}

	slug, err := firstUnsolvedInCourse(items, course)
	if err != nil {
		return err
	}
	return fetchAndPrepareLesson(slug)
}

// firstUnsolvedInCourse returns the slug of the first non-completed lesson in
// the given course, in the server-provided ordering.
func firstUnsolvedInCourse(items []client.ProgressItem, course string) (string, error) {
	var inCourse []client.ProgressItem
	for _, p := range items {
		if p.CourseSlug == course {
			inCourse = append(inCourse, p)
		}
	}
	if len(inCourse) == 0 {
		seen := map[string]bool{}
		var available []string
		for _, p := range items {
			if !seen[p.CourseSlug] {
				seen[p.CourseSlug] = true
				available = append(available, p.CourseSlug)
			}
		}
		hint := "run: lst progress"
		if len(available) > 0 {
			hint = "available courses: " + strings.Join(available, ", ")
		}
		return "", errorWithHint(
			fmt.Sprintf("course %q not found in your progress", course),
			hint,
		)
	}
	for _, p := range inCourse {
		if p.Status != "completed" {
			return p.LessonSlug, nil
		}
	}
	return "", fmt.Errorf("course %q is fully completed — congratulations", course)
}

func init() { rootCmd.AddCommand(startCmd) }
