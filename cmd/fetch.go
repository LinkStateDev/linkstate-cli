package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <slug>",
	Short: "Download a lesson to solve locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: lst auth")
		}
		slug := args[0]
		l, err := cliClient.GetLessonBySlug(slug)
		if err != nil { return fmt.Errorf("fetch: %w", err) }

		dir := filepath.Join(cfg.Path, l.CourseSlug, slug)
		if err := os.MkdirAll(dir, 0755); err != nil { return fmt.Errorf("create dir: %w", err) }

		zipURL := fmt.Sprintf("%s/api/download/%s", cfg.Server, slug)
		if err := downloadAndUnzip(zipURL, dir); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		// Select the right test binary for this platform
		testName := fmt.Sprintf("test-%s-%s", runtime.GOOS, runtime.GOARCH)
		if runtime.GOOS == "windows" { testName += ".exe" }
		testSrc := filepath.Join(dir, testName)
		testDst := filepath.Join(dir, "test")
		if runtime.GOOS == "windows" { testDst += ".exe" }
		if _, err := os.Stat(testSrc); err == nil {
			os.Rename(testSrc, testDst)
			os.Chmod(testDst, 0755)
		}
		// Clean up unused test binaries
		for _, name := range []string{"test-linux-amd64", "test-darwin-amd64", "test-darwin-arm64", "test-windows-amd64.exe"} {
			if name != testName {
				os.Remove(filepath.Join(dir, name))
			}
		}

		// Rename template.go to main.go if needed
		tmplFile := filepath.Join(dir, "template.go")
		mainFile := filepath.Join(dir, "main.go")
		if _, err := os.Stat(tmplFile); err == nil {
			os.Rename(tmplFile, mainFile)
		}

		meta, _ := json.MarshalIndent(map[string]any{"lesson_id": l.ID, "slug": slug, "title": l.Title}, "", "  ")
		if err := os.WriteFile(filepath.Join(dir, ".linkstate.json"), meta, 0644); err != nil {
			return fmt.Errorf("write .linkstate.json: %w", err)
		}
		fmt.Printf("Created %s/\n", dir)
		fmt.Println("  main.go              → your code")
		fmt.Println("  test                 → local test runner")
		fmt.Println("  .linkstate.json      → metadata")
		fmt.Println("  .linkstate.json      → metadata")
		fmt.Printf("\nNext: cd %s && lst test\n", dir)
		fmt.Println("      lst submit")
		return nil
	},
}

func downloadAndUnzip(url, dir string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return err }
	if cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	tmpFile := filepath.Join(os.TempDir(), "lst_download.zip")
	f, err := os.Create(tmpFile)
	if err != nil { return err }
	if _, err := io.Copy(f, resp.Body); err != nil { return err }
	f.Close()
	defer os.Remove(tmpFile)

	zr, err := zip.OpenReader(tmpFile)
	if err != nil { return err }
	defer zr.Close()

	for _, zf := range zr.File {
		dest := filepath.Join(dir, zf.Name)
		if zf.FileInfo().IsDir() {
			os.MkdirAll(dest, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(dest), 0755)
		out, err := os.Create(dest)
		if err != nil { return err }
		rc, err := zf.Open()
		if err != nil { return err }
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

func init() { rootCmd.AddCommand(fetchCmd) }

func findSolutionFile() string {
	if _, err := os.Stat("main.go"); err == nil { return "main.go" }
	return ""
}
