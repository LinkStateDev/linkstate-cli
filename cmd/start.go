package cmd

import (
	"fmt"
	"strings"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var startForce bool

var startCmd = &cobra.Command{
	Use:   "start <lesson-slug | course-slug/module-slug>",
	Short: "Start a new lesson or module from scratch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}
		if len(args) == 0 {
			return errorWithHint(
				"specify a lesson to start",
				"lst start hello-engineer\nlst start netops-framework/state-tester\nOr use 'lst resume' to continue your last lesson.",
			)
		}

		arg := args[0]
		var dir string
		var err error
		if strings.Contains(arg, "/") {
			dir, err = startModule(arg)
		} else {
			dir, err = startLesson(arg)
		}
		if err != nil {
			return err
		}
		return cdInto(dir)
	},
}

func startLesson(slug string) (string, error) {
	fmt.Printf("%s Starting lesson %s\n", ui.Muted.Render("→"), ui.Bold.Render(slug))
	return fetchAndPrepareLesson(slug, startForce)
}

func startModule(arg string) (string, error) {
	parts := strings.SplitN(arg, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", errorWithHint(
			fmt.Sprintf("invalid argument %q", arg),
			"use <course-slug>/<module-slug>, e.g. lst start netops-framework/state-tester",
		)
	}
	trackSlug, moduleSlug := parts[0], parts[1]

	var resp *client.StartResponse
	err := withSpinner("Looking up module…", func() error {
		var err error
		resp, err = cliClient.GetStartLesson(trackSlug, moduleSlug)
		return err
	})
	if err != nil {
		return "", fmt.Errorf("start module: %w", err)
	}

	if resp.ModuleCompleted {
		fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render("Module already completed — opening last lesson for review"))
	}
	return fetchAndPrepareLesson(resp.Lesson.Slug, startForce)
}

func init() {
	startCmd.Flags().BoolVarP(&startForce, "force", "f", false, "Overwrite existing lesson directory")
	rootCmd.AddCommand(startCmd)
}
