package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	serverURL string
	cliClient *client.Client
	cfg       *config.Config
)

// hintedError carries an actionable suggestion alongside the error message.
// Returned from RunE so the root error handler can render it as
// "✘ <msg>\n  → <hint>".
type hintedError struct {
	msg  string
	hint string
}

func (e *hintedError) Error() string { return e.msg }

func errorWithHint(msg, hint string) error {
	return &hintedError{msg: msg, hint: hint}
}

var rootCmd = &cobra.Command{
	Use:   "lst",
	Short: "LinkState — network automation learning platform",
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
		fmt.Println(ui.Bold.Render("LinkState CLI") + ui.Muted.Render(" — network automation courses"))
		fmt.Println()
		printRow := func(name, desc string) {
			fmt.Printf("  %-10s %s\n", ui.Hint.Render(name), ui.Muted.Render(desc))
		}
		printRow("auth", "Authenticate via browser")
		printRow("start", "Start a new lesson from scratch")
		printRow("resume", "Continue working on your last lesson")
		printRow("test", "Run local tests against your solution")
		printRow("submit", "Submit your solution result (--next to advance)")
		printRow("next", "Fetch the next lesson after the current one")
		printRow("progress", "Show your learning progress")
		printRow("hint", "Get a hint for a task")
		printRow("config", "Show or change settings")
		printRow("version", "Print version")
		printRow("logout", "Clear saved authentication")
		fmt.Println()
		fmt.Printf("%s %s\n", ui.Muted.Render("Server:"), ui.Hint.Render(cfg.Server))
		if cfg.Email != "" {
			fmt.Printf("%s %s\n", ui.Muted.Render("Logged in:"), cfg.Email)
		} else {
			fmt.Println(ui.Muted.Render("Not logged in. Run: lst auth"))
		}
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func printError(err error) {
	if errors.Is(err, client.ErrUnauthorized) {
		if cfg != nil && cfg.Token != "" {
			cfg.Token = ""
			cfg.Email = ""
			_ = config.Save(cfg)
		}
		fmt.Fprintf(os.Stderr, "%s %s\n", ui.Error.Render(ui.GlyphFail), "session expired")
		fmt.Fprintf(os.Stderr, "  %s %s\n", ui.Hint.Render("→"), ui.Hint.Render("run: lst auth"))
		return
	}
	msg := strings.TrimSpace(err.Error())
	fmt.Fprintf(os.Stderr, "%s %s\n", ui.Error.Render(ui.GlyphFail), msg)
	var he *hintedError
	if errors.As(err, &he) && he.hint != "" {
		fmt.Fprintf(os.Stderr, "  %s %s\n", ui.Hint.Render("→"), ui.Hint.Render(he.hint))
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "", "Server URL (default http://localhost)")
}
