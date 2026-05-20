package cmd

import (
	"fmt"
	"os"

	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	serverURL string
	cliClient *client.Client
	cfg       *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "linkstate-cli",
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
		fmt.Println("LinkStateDev CLI — network automation courses")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  login     Authenticate and save token")
		fmt.Println("  fetch     Download a task to solve locally")
		fmt.Println("  test      Run local tests against your solution")
		fmt.Println("  submit    Submit your solution result")
		fmt.Println("  progress  Show your learning progress")
		fmt.Println("")
		fmt.Printf("Server: %s\n", cfg.Server)
		if cfg.Email != "" {
			fmt.Printf("Logged in as: %s\n", cfg.Email)
		} else {
			fmt.Println("Not logged in. Run: linkstate-cli login")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "", "API server URL (default http://localhost:8080)")
}
