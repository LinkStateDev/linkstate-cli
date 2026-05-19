package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linkstate-cli",
	Short: "LinkStateDev CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("linkstate-cli - LinkStateDev toolkit")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
