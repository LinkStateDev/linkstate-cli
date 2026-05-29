package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LinkStateDev/linkstate-cli/internal/client"
	"github.com/LinkStateDev/linkstate-cli/internal/lesson"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	submitNext bool
	submitCmd  = &cobra.Command{
		Use:   "submit",
		Short: "Submit your solution result",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Token == "" {
				return errorWithHint("not logged in", "run: lst auth")
			}
			meta, err := lesson.LoadMeta()
			if err != nil {
				return err
			}

			testErr := runTests(true)
			if testErr != nil {
				fmt.Println()
				fmt.Printf("%s %s\n", ui.Error.Render(ui.GlyphFail), ui.Muted.Render("Fix errors and run 'lst submit' again."))
				return nil
			}

			tasks := loadSubmitTasks()
			if len(tasks) == 0 {
				return errorWithHint("no tasks found", "run lst start <slug> to download the lesson")
			}

			allDone, lastResp, err := submitAllTasks(tasks, meta)
			if err != nil {
				return fmt.Errorf("submit: %w", err)
			}

			fmt.Println()
			if allDone {
				fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render("Lesson completed!"))
				if lastResp != nil && lastResp.NextLessonSlug != nil && *lastResp.NextLessonSlug != "" {
					fmt.Printf("  %s %s\n", ui.Muted.Render("Next lesson:"), ui.Hint.Render(*lastResp.NextLessonSlug))
					if submitNext {
						dir, err := fetchAndPrepareLesson(*lastResp.NextLessonSlug, false)
						if err != nil {
							return err
						}
						return cdInto(dir)
					}
				}
			} else {
				fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render("All passing tasks submitted!"))
			}

			if !allDone {
				fmt.Printf("%s %s %s\n", ui.Muted.Render("→"), ui.Muted.Render("Keep coding and run:"), ui.Hint.Render("lst test && lst submit"))
			} else if !submitNext {
				fmt.Printf("%s %s %s\n", ui.Muted.Render("→"), ui.Muted.Render("Continue with:"), ui.Hint.Render("lst next"))
			}
			return nil
		},
	}
)

type submitTaskEntry struct {
	ID int `json:"id"`
}

func loadSubmitTasks() []submitTaskEntry {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return nil
	}
	var tasks []submitTaskEntry
	json.Unmarshal(data, &tasks)
	return tasks
}

func submitAllTasks(tasks []submitTaskEntry, meta lesson.Meta) (bool, *client.SubmitResponse, error) {
	progress, err := cliClient.GetProgress()
	if err != nil {
		return false, nil, fmt.Errorf("get progress: %w", err)
	}
	completed := map[int]bool{}
	for _, p := range progress {
		if p.LessonID == meta.LessonID && p.Status == "completed" {
			completed[p.TaskID] = true
		}
	}

	allDone := true
	var lastResp *client.SubmitResponse
	foundPending := false
	for _, t := range tasks {
		if completed[t.ID] {
			continue
		}
		foundPending = true
		resp, err := cliClient.Submit(t.ID, "pass")
		if err != nil {
			return false, nil, fmt.Errorf("submit task %d: %w", t.ID, err)
		}
		lastResp = resp
		allDone = resp.LessonCompleted
	}
	if !foundPending {
		return true, nil, nil
	}
	return allDone, lastResp, nil
}

func init() {
	submitCmd.Flags().BoolVar(&submitNext, "next", false, "After pass, automatically fetch the next lesson")
	rootCmd.AddCommand(submitCmd)
}
