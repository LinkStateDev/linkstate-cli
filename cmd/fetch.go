package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <slug-or-id>",
	Short: "Download a lesson to solve locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: lst auth")
		}
		arg := args[0]

		var lesson *struct {
			ID         int
			Slug       string
			Title      string
			Template   string
			TestConfig string
		}

		// Try as numeric ID first
		if id, err := strconv.Atoi(arg); err == nil {
			l, err := cliClient.GetLesson(id)
			if err != nil { return fmt.Errorf("fetch: %w", err) }
			lesson = &struct {
				ID         int
				Slug       string
				Title      string
				Template   string
				TestConfig string
			}{l.ID, l.Slug, l.Title, l.Template, l.TestConfig}
		} else {
			l, err := cliClient.GetLessonBySlug(arg)
			if err != nil { return fmt.Errorf("fetch: %w", err) }
			lesson = &struct {
				ID         int
				Slug       string
				Title      string
				Template   string
				TestConfig string
			}{l.ID, l.Slug, l.Title, l.Template, l.TestConfig}
		}

		dir := arg
		if lesson.Slug != "" { dir = lesson.Slug }
		if err := os.MkdirAll(dir, 0755); err != nil { return fmt.Errorf("create dir: %w", err) }

		if err := os.WriteFile(dir+"/main.go", []byte(lesson.Template), 0644); err != nil {
			return fmt.Errorf("write main.go: %w", err)
		}
		if err := os.WriteFile(dir+"/test_config.json", []byte(lesson.TestConfig), 0644); err != nil {
			return fmt.Errorf("write test_config.json: %w", err)
		}
		meta, _ := json.MarshalIndent(map[string]any{"lesson_id": lesson.ID, "slug": lesson.Slug, "title": lesson.Title}, "", "  ")
		if err := os.WriteFile(dir+"/.linkstate.json", meta, 0644); err != nil {
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
