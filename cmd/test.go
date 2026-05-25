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

var resultRe = regexp.MustCompile(`^\s*--- (PASS|FAIL|SKIP):?\s*(.*?)(?:\s+\(\d+\.\d+s\))?\s*$`)
var runRe = regexp.MustCompile(`^\s*=== RUN\s+(.*)$`)
var detailStripRe = regexp.MustCompile(`^\s*test_test\.go:\d+:\s*`)

type testResult struct {
	name    string
	failed  bool
	details []string
}

func runTests(submitting bool) error {
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		return fmt.Errorf("main.go not found in current directory")
	}
	testBin := "./test"
	if _, err := os.Stat(testBin); os.IsNotExist(err) {
		return fmt.Errorf("test binary not found. Run: lst fetch <slug>")
	}

	cmd := exec.Command(testBin, "-test.v")
	stdout, err := cmd.StdoutPipe()
	if err != nil { return err }
	cmd.Start()

	// First pass: collect details by test name (from === RUN lines)
	detailsByTest := make(map[string][]string)
	currentTest := ""
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "=== RUN ") {
			if m := runRe.FindStringSubmatch(text); m != nil {
				currentTest = strings.TrimSpace(m[1])
				if detailsByTest[currentTest] == nil {
					detailsByTest[currentTest] = nil
				}
			}
			continue
		}
		if resultRe.MatchString(text) || strings.TrimSpace(text) == "FAIL" {
			continue
		}
		if currentTest != "" && strings.TrimSpace(text) != "" {
			clean := detailStripRe.ReplaceAllString(strings.TrimSpace(text), "")
			if clean != "" {
				detailsByTest[currentTest] = append(detailsByTest[currentTest], clean)
			}
		}
	}
	cmd.Wait()

	// Second pass: collect result lines only
	cmd2 := exec.Command(testBin, "-test.v")
	stdout2, _ := cmd2.StdoutPipe()
	cmd2.Start()

	// Preserve test order from output
	var results []testResult
	seen := make(map[string]bool)
	scanner2 := bufio.NewScanner(stdout2)
	for scanner2.Scan() {
		text := scanner2.Text()
		m := resultRe.FindStringSubmatch(text)
		if m == nil { continue }
		name := strings.TrimSpace(m[2])
		if name == "" || seen[name] { continue }
		seen[name] = true

		r := testResult{name: name}
		if m[1] == "FAIL" { r.failed = true }
		// Merge details from first pass
		if d, ok := detailsByTest[name]; ok {
			r.details = d
		}
		results = append(results, r)
	}
	cmd2.Wait()

	// Filter: skip empty parent tests that have sub-tests
	filtered := make([]testResult, 0)
	for _, r := range results {
		hasSub := false
		for _, other := range results {
			if other.name != r.name && strings.HasPrefix(other.name, r.name+"/") {
				hasSub = true
				break
			}
		}
		if hasSub && len(r.details) == 0 { continue }
		filtered = append(filtered, r)
	}

	passed, failed := 0, 0
	for _, r := range filtered {
		if r.failed {
			fmt.Printf("  %s %s\n", color.Red("❌"), r.name)
			for _, d := range r.details {
				fmt.Printf("     %s\n", color.Yellow(d))
			}
			failed++
		} else {
			fmt.Printf("  %s %s\n", color.Green("✅"), r.name)
			passed++
		}
	}

	fmt.Println()
	if failed == 0 {
		if !submitting { fmt.Printf("All %d tests passed! Run: lst submit\n", passed) }
		return nil
	}
	fmt.Printf("%d passed, %d failed.\n", passed, failed)
	return nil
}

func init() { rootCmd.AddCommand(testCmd) }
