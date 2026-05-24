package cmd

import (
	"encoding/json"
	"fmt"
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

		mainFile := filepath.Join(dir, "main.go")
		if err := os.WriteFile(mainFile, []byte(l.Template), 0644); err != nil {
			return fmt.Errorf("write main.go: %w", err)
		}
		tcFile := filepath.Join(dir, "test_config.json")
		if err := os.WriteFile(tcFile, []byte(l.TestConfig), 0644); err != nil {
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

func init() { rootCmd.AddCommand(fetchCmd) }

func findSolutionFile() string {
	if _, err := os.Stat("main.go"); err == nil { return "main.go" }
	return ""
}
