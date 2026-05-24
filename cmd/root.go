package cmd

import (
	"fmt"
	"os"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	serverURL string
	cliClient *client.Client
	cfg       *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "lst",
	Short: "LinkStateDev — network automation learning platform",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if serverURL != "" {
			cfg.Server = serverURL
		}
		cliClient = client.New(cfg.Server, cfg.Token)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(color.Bold("LinkStateDev CLI — network automation courses"))
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  auth      Authenticate via browser")
		fmt.Println("  fetch     Download a task to solve locally")
		fmt.Println("  test      Run local tests against your solution")
		fmt.Println("  submit    Submit your solution result")
		fmt.Println("  progress  Show your learning progress")
		fmt.Println("  hint      Get a hint for the current task")
		fmt.Println("  config    Show or change settings")
		fmt.Println("  version   Print version")
		fmt.Println("  logout    Clear saved authentication")
		fmt.Println()
		fmt.Printf("Server: %s\n", color.Yellow(cfg.Server))
		if cfg.Email != "" {
			fmt.Printf("Logged in: %s\n", cfg.Email)
		} else {
			fmt.Println(color.Faint("Not logged in. Run: lst auth"))
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "", "Server URL (default http://localhost)")
}
