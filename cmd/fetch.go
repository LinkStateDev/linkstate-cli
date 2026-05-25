package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

		zipURL := fmt.Sprintf("%s/static/zips/%s.zip", cfg.Server, slug)
		if err := downloadAndUnzip(zipURL, dir); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		// Rename template.go to main.go if needed
		tmplFile := filepath.Join(dir, "template.go")
		mainFile := filepath.Join(dir, "main.go")
		if _, err := os.Stat(tmplFile); err == nil {
			os.Rename(tmplFile, mainFile)
		}

		// Write test_config.json from API
		if err := os.WriteFile(filepath.Join(dir, "test_config.json"), []byte(l.TestConfig), 0644); err != nil {
			return fmt.Errorf("write test_config.json: %w", err)
		}

		meta, _ := json.MarshalIndent(map[string]any{"lesson_id": l.ID, "slug": slug, "title": l.Title}, "", "  ")
		if err := os.WriteFile(filepath.Join(dir, ".linkstate.json"), meta, 0644); err != nil {
			return fmt.Errorf("write .linkstate.json: %w", err)
		}
		fmt.Printf("Created %s/\n", dir)
		fmt.Println("  main.go              → your code")
		fmt.Println("  test_config.json     → validation rules")
		fmt.Println("  .linkstate.json      → metadata")
		fmt.Printf("\nNext: cd %s && lst test\n", dir)
		fmt.Println("      lst submit")
		return nil
	},
}

func downloadAndUnzip(url, dir string) error {
	resp, err := http.Get(url)
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
