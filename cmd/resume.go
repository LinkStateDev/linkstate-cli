package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Continue working on your last lesson",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.LastLessonSlug == "" {
			return errorWithHint(
				"no last lesson — run 'lst test' from a lesson first",
				"Or start a new lesson: lst start hello-engineer",
			)
		}

		dir := filepath.Join(cfg.Path, cfg.LastTrackSlug, cfg.LastLessonSlug)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return errorWithHint(
				fmt.Sprintf("lesson directory not found: %s", dir),
				"Re-download with: lst start "+cfg.LastLessonSlug,
			)
		}

		fmt.Printf("%s %s\n", ui.Muted.Render("→"), ui.Muted.Render("Resuming "+cfg.LastLessonTitle))
		return cdInto(dir)
	},
}

func init() { rootCmd.AddCommand(resumeCmd) }
