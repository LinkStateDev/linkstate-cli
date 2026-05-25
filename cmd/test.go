package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests against your solution",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTests(false)
	},
}

var (
	resultRe      = regexp.MustCompile(`^\s*--- (PASS|FAIL|SKIP):?\s*(.*?)(?:\s+\(\d+\.\d+s\))?\s*$`)
	runRe         = regexp.MustCompile(`^\s*=== RUN\s+(.*)$`)
	detailStripRe = regexp.MustCompile(`^\s*test_test\.go:\d+:\s*`)
)

type testResult struct {
	name    string
	failed  bool
	skipped bool
	details []string
}

func runTests(submitting bool) error {
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		return errorWithHint("main.go not found in current directory", "run lst from a fetched lesson directory")
	}
	testBin := "./test"
	if _, err := os.Stat(testBin); os.IsNotExist(err) {
		return errorWithHint("test binary not found", "fetch the lesson first: lst fetch <slug>")
	}

	cmd := exec.Command(testBin, "-test.v")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	var (
		results       []testResult
		seen          = map[string]bool{}
		detailsByTest = map[string][]string{}
		currentTest   string
	)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()

		if m := runRe.FindStringSubmatch(text); m != nil {
			currentTest = strings.TrimSpace(m[1])
			continue
		}

		if m := resultRe.FindStringSubmatch(text); m != nil {
			name := strings.TrimSpace(m[2])
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			r := testResult{name: name, details: detailsByTest[name]}
			switch m[1] {
			case "FAIL":
				r.failed = true
			case "SKIP":
				r.skipped = true
			}
			results = append(results, r)
			continue
		}

		if strings.TrimSpace(text) == "FAIL" {
			continue
		}

		if currentTest != "" && strings.TrimSpace(text) != "" {
			clean := detailStripRe.ReplaceAllString(strings.TrimSpace(text), "")
			if clean != "" {
				detailsByTest[currentTest] = append(detailsByTest[currentTest], clean)
				for i := range results {
					if results[i].name == currentTest {
						results[i].details = detailsByTest[currentTest]
					}
				}
			}
		}
	}
	cmd.Wait()

	// Hide a parent test row when its sub-tests already cover the failure detail.
	filtered := make([]testResult, 0, len(results))
	for _, r := range results {
		hasSub := false
		for _, other := range results {
			if other.name != r.name && strings.HasPrefix(other.name, r.name+"/") {
				hasSub = true
				break
			}
		}
		if hasSub && len(r.details) == 0 {
			continue
		}
		filtered = append(filtered, r)
	}

	passed, failed, skipped := 0, 0, 0
	fmt.Println()
	for i, r := range filtered {
		if i > 0 {
			fmt.Println()
		}
		switch {
		case r.skipped:
			skipped++
			fmt.Printf("  %s  %s\n", ui.Muted.Render("○"), ui.Muted.Render(humanName(r.name)))
		case r.failed:
			failed++
			fmt.Printf("  %s  %s\n", ui.Error.Render(ui.GlyphFail), ui.Bold.Render(humanName(r.name)))
			if len(r.details) > 0 {
				fmt.Println(ui.FailBlock.Render(formatDetails(r.details)))
			}
		default:
			passed++
			fmt.Printf("  %s  %s\n", ui.Success.Render(ui.GlyphPass), humanName(r.name))
		}
	}

	fmt.Println()
	summary := fmt.Sprintf("%d passed", passed)
	if failed > 0 {
		summary += fmt.Sprintf(" · %d failed", failed)
	}
	if skipped > 0 {
		summary += fmt.Sprintf(" · %d skipped", skipped)
	}
	if failed == 0 && passed > 0 {
		fmt.Println(ui.SummaryPass.Render(summary))
		if !submitting {
			fmt.Println()
			fmt.Println("  " + ui.Muted.Render("All green. Run: ") + ui.Hint.Render("lst submit"))
		}
		return nil
	}
	fmt.Println(ui.SummaryFail.Render(summary))
	return fmt.Errorf("%d test(s) failed", failed)
}

var expectedGotRe = regexp.MustCompile(`(?i)^\s*(expected|got|want|have|actual)\b[\s:=-]*(.*)$`)

// formatDetails groups raw "expected: X / got: Y"–style lines into a tidy
// block. Lines that don't match the pattern are kept verbatim.
func formatDetails(lines []string) string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		l = cleanDetail(l)
		if l == "" {
			continue
		}
		if m := expectedGotRe.FindStringSubmatch(l); m != nil {
			label := strings.ToLower(m[1])
			value := strings.TrimSpace(m[2])
			switch label {
			case "got", "actual", "have":
				out = append(out, fmt.Sprintf("%s  %s", ui.Error.Render(pad(label)), value))
			default:
				out = append(out, fmt.Sprintf("%s  %s", ui.Hint.Render(pad(label)), value))
			}
			continue
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

func pad(s string) string {
	const w = 8
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}

func humanName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx > 0 {
		name = name[idx+1:]
	}
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.TrimPrefix(name, "Test")
	return name
}

func cleanDetail(d string) string {
	d = strings.ReplaceAll(d, "\\n", "")
	return d
}

func init() { rootCmd.AddCommand(testCmd) }
