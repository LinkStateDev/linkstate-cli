package cmd

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/mattn/go-isatty"
)

// withSpinner runs fn while showing a spinner with the given message. The
// spinner writes to stderr so it does not pollute parsable stdout, and is
// suppressed entirely when stderr is not a TTY.
func withSpinner(msg string, fn func() error) error {
	if !isatty.IsTerminal(os.Stderr.Fd()) {
		return fn()
	}
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " " + ui.Muted.Render(msg)
	s.HideCursor = true
	s.Start()
	defer s.Stop()
	return fn()
}
