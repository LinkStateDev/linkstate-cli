// Package ui defines semantic terminal styles used across lst commands.
//
// All styles route through lipgloss, which automatically degrades to plain
// text when stdout is not a TTY or NO_COLOR is set, so callers don't need to
// branch on terminal capability.
package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

func init() {
	if os.Getenv("NO_COLOR") != "" || !isatty.IsTerminal(os.Stdout.Fd()) {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
}

var (
	successColor = lipgloss.Color("10")
	errorColor   = lipgloss.Color("9")
	warnColor    = lipgloss.Color("11")
	hintColor    = lipgloss.Color("14")
	mutedColor   = lipgloss.Color("8")

	Success = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	Error   = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	Warn    = lipgloss.NewStyle().Foreground(warnColor).Bold(true)
	Hint    = lipgloss.NewStyle().Foreground(hintColor)
	Muted   = lipgloss.NewStyle().Foreground(mutedColor)
	Bold    = lipgloss.NewStyle().Bold(true)
	Code    = lipgloss.NewStyle().Foreground(hintColor)

	FailBlock = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(errorColor).
			Padding(0, 0, 0, 1).
			MarginLeft(2)

	SummaryPass = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(0, 1).
			MarginLeft(2)

	SummaryFail = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(0, 1).
			MarginLeft(2)
)

// Glyphs that look right in monospace; use these instead of inline emoji
// to keep widths predictable across terminals.
const (
	GlyphPass = "✓"
	GlyphFail = "✘"
	GlyphInfo = "ℹ"
	GlyphHint = "💡"
	GlyphLock = "🔒"
)
