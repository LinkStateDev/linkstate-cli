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

type TestConfig struct {
	Type           string     `json:"type"`
	Run            string     `json:"run"`
	Stdin          string     `json:"stdin,omitempty"`
	ExpectedStdout string     `json:"expected_stdout,omitempty"`
	Port           int        `json:"port,omitempty"`
	StartupWaitMs  int        `json:"startup_wait_ms,omitempty"`
	Tests          []HTTPTest `json:"tests,omitempty"`
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
	args := strings.Fields(cfg.Run)
	for i := range args {
		if i+1 < len(args) && (args[i+1] == "solution.py" || args[i+1] == "solution.go") {
			args[i+1] = solutionFile
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(cfg.Stdin)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	actual := stdout.String()

	r := TestResult{
		Name:     "output",
		Expected: cfg.ExpectedStdout,
		Actual:   actual,
	}

	if err != nil {
		r.Passed = false
		if r.Actual == "" {
			r.Actual = fmt.Sprintf("process error: %v", err)
		}
		return &Report{
			Failed:  1,
			Results: []TestResult{r},
		}, nil
	}

	r.Passed = actual == cfg.ExpectedStdout
	report := &Report{Results: []TestResult{r}}
	if r.Passed {
		report.Passed = 1
		report.AllPass = true
	} else {
		report.Failed = 1
	}
	return report, nil
}

func runServer(cfg *TestConfig, solutionFile string) (*Report, error) {
	args := strings.Fields(cfg.Run)
	for i := range args {
		if i+1 < len(args) && (args[i+1] == "solution.py" || args[i+1] == "solution.go") {
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
