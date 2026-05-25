package cmd

import (
	"fmt"
	"strconv"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/lesson"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var hintCmd = &cobra.Command{
	Use:   "hint [level]",
	Short: "Show a hint for the current task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}

		meta, err := lesson.LoadMeta()
		if err != nil {
			return err
		}
		if meta.Slug == "" {
			return errorWithHint("no slug in .linkstate.json", "re-fetch the lesson with: lst fetch <slug>")
		}

		level := 1
		if len(args) > 0 {
			level, _ = strconv.Atoi(args[0])
		}
		if level < 1 {
			level = 1
		}

		var resp *client.HintResponse
		err = withSpinner("Loading hint…", func() error {
			r, err := cliClient.GetHint(meta.Slug, level)
			if err != nil {
				return err
			}
			resp = r
			return nil
		})
		if err != nil {
			return fmt.Errorf("get hint: %w", err)
		}

		if resp.Hint == "" {
			return fmt.Errorf("no hint %d available", level)
		}

		fmt.Printf("%s %s\n", ui.Bold.Render(ui.GlyphHint+" Hint"), ui.Muted.Render(fmt.Sprintf("(%d/%d):", resp.Level, resp.Total)))
		fmt.Println()

		r, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(90))
		out, err := r.Render(resp.Hint)
		if err != nil {
			fmt.Println(resp.Hint)
		} else {
			fmt.Print(out)
		}

		if resp.Level < resp.Total {
			fmt.Printf("\n%s\n", ui.Muted.Render(fmt.Sprintf("Need more details? Run: lst hint %d", resp.Level+1)))
		} else {
			fmt.Println(ui.Muted.Render("\nThat was the last hint. Good luck!"))
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(hintCmd) }
