package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/LinkStateDev/linkstate-cli/internal/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests against your solution",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTests(false)
	},
}

type testLine struct {
	passed bool
	fail   bool
	name   string
	detail string
}

var testLineRe = regexp.MustCompile(`^\s*(--- (PASS|FAIL|SKIP):?\s*(.*?)(?:\s+\(\d+\.\d+s\))?\s*)$`)
var testDetailRe = regexp.MustCompile(`^\s*test_test.go:\d+:\s*`)

func runTests(submitting bool) error {
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		return fmt.Errorf("main.go not found in current directory")
	}
	testBin := "./test"
	if _, err := os.Stat(testBin); os.IsNotExist(err) {
		return fmt.Errorf("test binary not found. Run: lst fetch <slug>")
	}

	exec.Command("go", "build", "-o", "solution", "main.go").Run()
	defer os.Remove("solution")

	cmd := exec.Command(testBin, "-test.v")
	stdout, err := cmd.StdoutPipe()
	if err != nil { return err }
	cmd.Start()

	var lines []testLine
	failCount := 0
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "=== ") { continue } // skip === RUN/PASS

		m := testLineRe.FindStringSubmatch(text)
		if m == nil {
			// detail line (test output)
			if len(lines) > 0 && strings.TrimSpace(text) != "" {
				clean := testDetailRe.ReplaceAllString(strings.TrimSpace(text), "")
				last := &lines[len(lines)-1]
				if last.detail != "" && clean != "" { last.detail += "\n" }
				last.detail += clean
			}
			continue
		}
		name := strings.TrimSpace(m[3])
		if name == "" { continue }
		l := testLine{name: name}
		switch m[2] {
		case "PASS":
			l.passed = true
		case "FAIL":
			l.fail = true
			failCount++
		case "SKIP":
		}
		lines = append(lines, l)
	}
	cmd.Wait()

	passed := 0
	for _, l := range lines {
		if l.fail {
			fmt.Printf("  %s %s: %s\n", color.Red("❌"), l.name, color.Red("FAIL"))
			if l.detail != "" {
				fmt.Printf("     %s\n", color.Yellow(l.detail))
			}
		} else if l.passed {
			fmt.Printf("  %s %s: %s\n", color.Green("✅"), l.name, color.Green("PASS"))
			passed++
		}
	}

	fmt.Println()
	if failCount == 0 {
		if !submitting { fmt.Printf("All %d tests passed! Run: lst submit\n", passed) }
		return nil
	}
	fmt.Printf("%d passed, %d failed.\n", passed, failCount)
	return fmt.Errorf("tests failed")
}

func init() { rootCmd.AddCommand(testCmd) }
