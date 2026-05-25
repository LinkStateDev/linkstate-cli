package color

func Green(msg string) string  { return "\033[32m" + msg + "\033[0m" }
func Red(msg string) string    { return "\033[31m" + msg + "\033[0m" }
func Yellow(msg string) string { return "\033[33m" + msg + "\033[0m" }
func Bold(msg string) string   { return "\033[1m" + msg + "\033[0m" }
func Faint(msg string) string  { return "\033[2m" + msg + "\033[0m" }
