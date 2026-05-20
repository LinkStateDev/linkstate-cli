package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <task-id>",
	Short: "Download a task to solve locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Run: linkstate-cli login")
		}

		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task id: %s", args[0])
		}

		task, err := cliClient.GetTask(taskID)
		if err != nil {
			return fmt.Errorf("fetch task: %w", err)
		}

		dir := fmt.Sprintf("task-%d", taskID)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create dir: %w", err)
		}

		solutionFile := fmt.Sprintf("%s/solution.py", dir)
		if err := os.WriteFile(solutionFile, []byte(task.Template), 0644); err != nil {
			return fmt.Errorf("write solution.py: %w", err)
		}

		testConfigFile := fmt.Sprintf("%s/test_config.json", dir)
		if err := os.WriteFile(testConfigFile, []byte(task.TestConfig), 0644); err != nil {
			return fmt.Errorf("write test_config.json: %w", err)
		}

		meta := map[string]any{
			"task_id":    task.ID,
			"lesson_id":  task.LessonID,
			"title":      task.Title,
			"task_type":  task.TaskType,
		}
		metaData, _ := json.MarshalIndent(meta, "", "  ")
		metaFile := fmt.Sprintf("%s/.linkstate-task.json", dir)
		if err := os.WriteFile(metaFile, metaData, 0644); err != nil {
			return fmt.Errorf("write .linkstate-task.json: %w", err)
		}

		fmt.Printf("Created %s/\n", dir)
		fmt.Printf("  solution.py          → your code goes here\n")
		fmt.Printf("  test_config.json     → validation rules\n")
		fmt.Printf("  .linkstate-task.json → metadata\n")
		fmt.Println()
		fmt.Println("Next: edit solution.py, then run:")
		fmt.Printf("  cd %s && linkstate-cli test\n", dir)
		fmt.Println("  linkstate-cli submit")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
