package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func readPassword() ([]byte, error) {
	if term.IsTerminal(int(syscall.Stdin)) {
		return term.ReadPassword(int(syscall.Stdin))
	}
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(strings.TrimSpace(line)), nil
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to LinkStateDev",
	RunE: func(cmd *cobra.Command, args []string) error {
		var email, password string

		fmt.Print("Email: ")
		fmt.Scanln(&email)
		email = strings.TrimSpace(email)
		if email == "" {
			return fmt.Errorf("email required")
		}

		fmt.Print("Password: ")
		bytepw, err := readPassword()
		fmt.Println()
		if err != nil {
			return fmt.Errorf("read password: %w", err)
		}
		password = strings.TrimSpace(string(bytepw))
		if password == "" {
			return fmt.Errorf("password required")
		}

		token, err := cliClient.Login(email, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		cfg.Token = token
		cfg.Email = email
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}

		fmt.Printf("Logged in as %s\n", email)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
