package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

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

	cmd2 := exec.Command(testBin, "-test.v")
	stdout2, _ := cmd2.StdoutPipe()
	cmd2.Start()

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
		if d, ok := detailsByTest[name]; ok {
			r.details = d
		}
		results = append(results, r)
	}
	cmd2.Wait()

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
	fmt.Println()
	for _, r := range filtered {
		name := cleanTestName(r.name)
		status := "PASS"
		if r.failed { status = "FAIL"; failed++ } else { passed++ }

		fmt.Printf("  %-44s %s\n", name, status)
		if r.failed && len(r.details) > 0 {
			for _, d := range r.details {
				fmt.Printf("    %s\n", d)
			}
		}
	}

	fmt.Println()
	if failed == 0 {
		if !submitting { fmt.Printf("All %d tests passed. Run: lst submit\n", passed) }
		return nil
	}
	fmt.Printf("%d passed, %d failed.\n", passed, failed)
	return nil
}

func cleanTestName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx > 0 {
		return name[idx+1:]
	}
	return name
}

func init() { rootCmd.AddCommand(testCmd) }
