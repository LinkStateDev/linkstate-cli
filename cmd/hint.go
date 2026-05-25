package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var hintCmd = &cobra.Command{
	Use:   "hint [level]",
	Short: "Show a hint for the current task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" { return fmt.Errorf("not logged in. Run: lst auth") }

		metaData, err := os.ReadFile(".linkstate.json")
		if err != nil { return fmt.Errorf(".linkstate.json not found. Run: lst fetch <slug>") }
		var meta struct {
			Slug string `json:"slug"`
		}
		if err := json.Unmarshal(metaData, &meta); err != nil { return fmt.Errorf("parse .linkstate.json: %w", err) }
		if meta.Slug == "" { return fmt.Errorf("no slug in .linkstate.json. Re-fetch the lesson") }

		level := 1
		if len(args) > 0 { level, _ = strconv.Atoi(args[0]) }
		if level < 1 { level = 1 }

		resp, err := cliClient.GetHint(meta.Slug, level)
		if err != nil { return fmt.Errorf("get hint: %w", err) }

		lvl, _ := resp["level"].(float64)
		total, _ := resp["total"].(float64)
		hint, _ := resp["hint"].(string)

		if hint == "" { return fmt.Errorf("no hint %d available", level) }

		fmt.Printf("%s %s\n", color.Bold(color.Yellow("💡 Hint")), color.Faint(fmt.Sprintf("(%d/%d):", int(lvl), int(total))))
		fmt.Println()

		r, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(90))
		out, err := r.Render(hint)
		if err != nil { fmt.Println(hint) } else { fmt.Print(out) }

		if int(lvl) < int(total) {
			fmt.Printf("\nNeed more details? Run: lst hint %d\n", int(lvl)+1)
		} else {
			fmt.Println("\nThat was the last hint. Good luck!")
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(hintCmd) }
