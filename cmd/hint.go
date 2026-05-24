package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/spf13/cobra"
)

var hintCmd = &cobra.Command{
	Use:   "hint [level]",
	Short: "Show a hint for the current task",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile("test_config.json")
		if err != nil {
			return fmt.Errorf("not in a task directory? Run: lst fetch <id>")
		}

		var cfg struct {
			Hints []string `json:"hints"`
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parse test_config.json: %w", err)
		}

		if len(cfg.Hints) == 0 {
			fmt.Println("No hints available for this task.")
			return nil
		}

		level := 1
		if len(args) > 0 {
			level, _ = strconv.Atoi(args[0])
		}
		if level < 1 {
			level = 1
		}
		if level > len(cfg.Hints) {
			level = len(cfg.Hints)
		}

		hint := cfg.Hints[level-1]
		fmt.Printf("%s %s\n", color.Bold(color.Yellow("💡 Hint")), color.Faint(fmt.Sprintf("(%d/%d):", level, len(cfg.Hints))))
		fmt.Println(hint)

		if level < len(cfg.Hints) {
			fmt.Printf("\nNeed more details? Run: lst hint %d\n", level+1)
		} else if level > 1 {
			fmt.Println("\nThat was the last hint. Good luck!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hintCmd)
}
