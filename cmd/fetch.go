package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:    "fetch <slug>",
	Short:  "Download a lesson to solve locally",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return errorWithHint("not logged in", "run: lst auth")
		}
		_, err := fetchAndPrepareLesson(args[0], false)
		return err
	},
}

// fetchAndPrepareLesson downloads the lesson identified by slug into
// cfg.Path/<course>/<slug>, sets up go.mod and .linkstate.json, and prints
// the "Next:" hint. Shared by `lst fetch`, `lst next`, and `lst resume`.
// If force is false and the target directory already contains files, the
// download is refused. Returns the directory where the lesson was placed.
func fetchAndPrepareLesson(slug string, force bool) (string, error) {
	var l *client.Lesson
	err := withSpinner("Fetching lesson…", func() error {
		var err error
		l, err = cliClient.GetLessonBySlug(slug)
		return err
	})
	if err != nil {
		return "", errorWithHint(
			fmt.Sprintf("lesson %q not found", slug),
			fmt.Sprintf("Browse all lessons at %s/courses", cfg.Server),
		)
	}

	dir := filepath.Join(cfg.Path, l.TrackSlug, slug)
	if !force {
		if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
			return "", errorWithHint(
				fmt.Sprintf("directory %s already exists and is not empty", dir),
				"use --force to overwrite, or run 'lst next' to continue",
			)
		}
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create dir: %w", err)
	}

	zipURL := fmt.Sprintf("%s/api/download/%s/%s-%s", cfg.Server, slug, runtime.GOOS, runtime.GOARCH)
	if err := downloadAndUnzip(zipURL, dir, slug); err != nil {
		return "", fmt.Errorf("download: %w", err)
	}

	os.Chmod(filepath.Join(dir, "test"), 0755)

	if out, err := exec.Command("go", "version").Output(); err == nil {
		ver := strings.TrimPrefix(strings.Fields(string(out))[2], "go")
		modContent := fmt.Sprintf("module solution\n\ngo %s\n", ver)
		os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644)
	}

	metaMap := map[string]any{"lesson_id": l.ID, "course_slug": l.TrackSlug, "slug": slug, "title": l.Title}
	meta, _ := json.MarshalIndent(metaMap, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, ".linkstate.json"), meta, 0644); err != nil {
		return "", fmt.Errorf("write .linkstate.json: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render(l.Title))
	fmt.Println(ui.Muted.Render("  " + dir))
	fmt.Println()
	fmt.Println(ui.Muted.Render("Next:"))
	fmt.Printf("  %s\n", ui.Hint.Render("lst test"))
	fmt.Printf("  %s\n", ui.Hint.Render("lst submit"))
	return dir, nil
}

func downloadAndUnzip(url, dir, slug string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return errorWithHint(
			fmt.Sprintf("no downloadable content for %q", slug),
			fmt.Sprintf("This lesson is available on the web at %s/lessons/%s", cfg.Server, slug),
		)
	}
	if resp.StatusCode == http.StatusForbidden {
		return errorWithHint("subscription required", "Upgrade your account on the web to download this lesson")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	tmpFile := filepath.Join(os.TempDir(), "lst_download.zip")
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription(ui.Muted.Render("Downloading "+slug+"…")),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(80*1e6),
		progressbar.OptionOnCompletion(func() { fmt.Fprint(os.Stderr, "\n") }),
		progressbar.OptionSetWidth(30),
		progressbar.OptionClearOnFinish(),
	)
	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()
	defer os.Remove(tmpFile)

	zr, err := zip.OpenReader(tmpFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, zf := range zr.File {
		cleanName := filepath.Clean(zf.Name)
		if strings.HasPrefix(cleanName, "..") {
			return fmt.Errorf("invalid zip entry: %s", zf.Name)
		}
		dest := filepath.Join(dir, cleanName)
		if !strings.HasPrefix(filepath.Clean(dest)+string(filepath.Separator), filepath.Clean(dir)+string(filepath.Separator)) {
			return fmt.Errorf("zip entry escapes directory: %s", zf.Name)
		}
		if zf.UncompressedSize64 > 10<<20 {
			return fmt.Errorf("zip entry too large (%d bytes): %s", zf.UncompressedSize64, zf.Name)
		}
		if zf.FileInfo().IsDir() {
			os.MkdirAll(dest, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(dest), 0755)
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		rc, err := zf.Open()
		if err != nil {
			return err
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

func init() { rootCmd.AddCommand(fetchCmd) }
