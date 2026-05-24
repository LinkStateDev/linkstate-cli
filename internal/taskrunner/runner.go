package taskrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	Args           []string `json:"args,omitempty"`
	Stdin          string   `json:"stdin,omitempty"`
	ExpectedStdout string   `json:"expected_stdout,omitempty"`
}

type TestConfig struct {
	Type           string     `json:"type"`
	Run            string     `json:"run"`
	Args           []string   `json:"args,omitempty"`
	Stdin          string     `json:"stdin,omitempty"`
	ExpectedStdout string     `json:"expected_stdout,omitempty"`
	Cases          []TestCase `json:"cases,omitempty"`
	Port           int        `json:"port,omitempty"`
	StartupWaitMs  int        `json:"startup_wait_ms,omitempty"`
	Tests          []HTTPTest `json:"tests,omitempty"`
	Hints          []string   `json:"hints,omitempty"`
}

type HTTPTest struct {
	Method       string `json:"method"`
	Path         string `json:"path"`
	Body         string `json:"body,omitempty"`
	ExpectStatus *int   `json:"expect_status,omitempty"`
	ExpectBody   string `json:"expect_body,omitempty"`
}

type TestResult struct {
	Name     string
	Passed   bool
	Expected string
	Actual   string
}

type Report struct {
	Passed  int
	Failed  int
	Results []TestResult
	AllPass bool
}

func Run(configJSON, solutionFile string) (*Report, error) {
	var cfg TestConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("parse test_config: %w", err)
	}
	if cfg.Type == "" {
		return nil, fmt.Errorf("invalid test_config. Try: lst fetch <id> again to get the latest task config")
	}

	switch cfg.Type {
	case "output":
		return runOutput(&cfg, solutionFile)
	case "server":
		return runServer(&cfg, solutionFile)
	default:
		return nil, fmt.Errorf("unknown task_type: %q", cfg.Type)
	}
}

func runOutput(cfg *TestConfig, solutionFile string) (*Report, error) {
	if len(cfg.Cases) > 0 {
		return runCases(cfg, solutionFile)
	}
	return runSingle(cfg, solutionFile, cfg.Args, cfg.Stdin, cfg.ExpectedStdout, "output")
}

func runSingle(cfg *TestConfig, solutionFile string, args []string, stdin, expected, name string) (*Report, error) {
	baseArgs := strings.Fields(cfg.Run)
	for i := range baseArgs {
		if i+1 < len(baseArgs) && baseArgs[i+1] == "main.go" {
			baseArgs[i+1] = solutionFile
		}
	}
	allArgs := append(baseArgs, args...)

	cmd := exec.Command(allArgs[0], allArgs[1:]...)
	cmd.Stdin = strings.NewReader(stdin)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	actual := stdout.String()

	r := TestResult{
		Name:     name,
		Expected: expected,
		Actual:   actual,
	}

	if err != nil {
		r.Passed = false
		if r.Actual == "" {
			r.Actual = fmt.Sprintf("process error: %v", err)
		}
		return &Report{Failed: 1, Results: []TestResult{r}}, nil
	}

	r.Passed = actual == expected
	report := &Report{Results: []TestResult{r}}
	if r.Passed {
		report.Passed = 1
		report.AllPass = true
	} else {
		report.Failed = 1
	}
	return report, nil
}

func runCases(cfg *TestConfig, solutionFile string) (*Report, error) {
	report := &Report{AllPass: true}
	for i, c := range cfg.Cases {
		name := fmt.Sprintf("case %d", i+1)
		r, err := runSingle(cfg, solutionFile, c.Args, c.Stdin, c.ExpectedStdout, name)
		if err != nil {
			return nil, err
		}
		report.Results = append(report.Results, r.Results[0])
		if r.AllPass {
			report.Passed++
		} else {
			report.Failed++
			report.AllPass = false
		}
	}
	return report, nil
}

func runServer(cfg *TestConfig, solutionFile string) (*Report, error) {
	args := strings.Fields(cfg.Run)
	for i := range args {
		if i+1 < len(args) && args[i+1] == "main.go" {
			args[i+1] = solutionFile
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start server: %w", err)
	}
	defer cmd.Process.Kill()

	waitMs := cfg.StartupWaitMs
	if waitMs == 0 {
		waitMs = 500
	}
	time.Sleep(time.Duration(waitMs) * time.Millisecond)

	baseURL := fmt.Sprintf("http://localhost:%d", cfg.Port)
	report := &Report{AllPass: true}

	for i, test := range cfg.Tests {
		r := TestResult{Name: fmt.Sprintf("test %d (%s %s)", i+1, test.Method, test.Path)}
		r.Passed, r.Expected, r.Actual = runHTTPTest(baseURL, test)
		report.Results = append(report.Results, r)
		if r.Passed {
			report.Passed++
		} else {
			report.Failed++
			report.AllPass = false
		}
	}

	return report, nil
}

func runHTTPTest(baseURL string, test HTTPTest) (passed bool, expected, actual string) {
	req, err := http.NewRequest(test.Method, baseURL+test.Path, strings.NewReader(test.Body))
	if err != nil {
		return false, "", fmt.Sprintf("request error: %v", err)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Sprintf("connection error: %v", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	body := buf.String()

	if test.ExpectStatus != nil {
		expected = fmt.Sprintf("status %d", *test.ExpectStatus)
		actual = fmt.Sprintf("status %d", resp.StatusCode)
		if resp.StatusCode != *test.ExpectStatus {
			return false, expected, actual
		}
	}
	if test.ExpectBody != "" && !strings.Contains(body, test.ExpectBody) {
		expected = test.ExpectBody
		actual = body
		if len(actual) > 200 {
			actual = actual[:200] + "..."
		}
		return false, expected, actual
	}
	return true, "", ""
}
