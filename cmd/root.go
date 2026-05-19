package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nbc",
	Short: "NetByCode CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nbc - netbycode toolkit")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}
