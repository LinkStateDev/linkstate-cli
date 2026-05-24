package cmd

import (
	"fmt"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear saved authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.Token = ""
		cfg.Email = ""
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Println(color.Faint("Logged out."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
