package cmd

import (
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
		level := 1
		if len(args) > 0 { level, _ = strconv.Atoi(args[0]) }
		if level < 1 { level = 1 }
		if level > 3 { level = 3 }

		hintFile := fmt.Sprintf("hint%d.md", level)
		data, err := os.ReadFile(hintFile)
		if err != nil {
			return fmt.Errorf("hint %d not available", level)
		}

		fmt.Printf("%s %s\n", color.Bold(color.Yellow("💡 Hint")), color.Faint(fmt.Sprintf("(%d/3):", level)))
		fmt.Println()

		r, _ := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(90),
		)
		out, err := r.Render(string(data))
		if err != nil {
			fmt.Println(string(data))
		} else {
			fmt.Print(out)
		}

		if level < 3 {
			if _, err := os.Stat(fmt.Sprintf("hint%d.md", level+1)); err == nil {
				fmt.Printf("\nNeed more details? Run: lst hint %d\n", level+1)
			}
		} else {
			fmt.Println("\nThat was the last hint. Good luck!")
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(hintCmd) }
