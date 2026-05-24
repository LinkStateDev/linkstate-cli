package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <lesson-id>",
	Short: "Download a lesson to solve locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: lst auth")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil { return fmt.Errorf("invalid lesson id: %s", args[0]) }
		lesson, err := cliClient.GetLesson(id)
		if err != nil { return fmt.Errorf("fetch: %w", err) }

		dir := filepath.Join(cfg.Workspace, fmt.Sprintf("lesson-%d", id))
		if err := os.MkdirAll(dir, 0755); err != nil { return fmt.Errorf("create dir: %w", err) }

		if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(lesson.Template), 0644); err != nil {
			return fmt.Errorf("write main.go: %w", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "test_config.json"), []byte(lesson.TestConfig), 0644); err != nil {
			return fmt.Errorf("write test_config.json: %w", err)
		}
		meta, _ := json.MarshalIndent(map[string]any{"lesson_id": id, "title": lesson.Title}, "", "  ")
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
