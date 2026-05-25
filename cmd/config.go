package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or change settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Printf("%s %s\n", ui.Bold.Render("Server:"), ui.Hint.Render(cfg.Server))
			fmt.Printf("%s %s\n", ui.Bold.Render("Path:"), ui.Hint.Render(cfg.Path))
			if cfg.Email != "" {
				fmt.Printf("Logged in: %s\n", cfg.Email)
			} else {
				fmt.Println("Logged in: no (run lst auth)")
			}
			return nil
		}
		if args[0] != "set" || len(args) < 3 {
			return fmt.Errorf("usage: lst config set <key> <value>")
		}
		key, val := args[1], args[2]
		switch key {
		case "server":
			cfg.Server = val
		case "path":
			cfg.Path = val
		default:
			return fmt.Errorf("unknown key: %s (available: server, path)", key)
		}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save: %w", err)
		}
		fmt.Printf("%s = %v\n", key, val)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
