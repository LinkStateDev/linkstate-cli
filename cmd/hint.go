package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/LinkStateDev/linkstate-cli/internal/lesson"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

type taskHintEntry struct {
	ID    int      `json:"id"`
	Slug  string   `json:"slug"`
	Title string   `json:"title"`
	Hints []string `json:"hints"`
}

var hintCmd = &cobra.Command{
	Use:   "hint [task-slug]",
	Short: "Show hints for a task",
	Long: `Show progressive hints for a task.

  lst hint                    — pick a task from the list
  lst hint parse-inventory    — next hint for this task
  lst hint parse-inventory 2  — specific hint level`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tasks, err := loadTasksJSON()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return hintDialog(tasks)
		}

		slug := args[0]
		level := 0
		if len(args) > 1 {
			level, _ = strconv.Atoi(args[1])
		}
		return showHint(tasks, slug, level)
	},
}

func loadTasksJSON() ([]taskHintEntry, error) {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return nil, errorWithHint(
			"tasks.json not found",
			"run lst start <slug> to download the lesson",
		)
	}
	var tasks []taskHintEntry
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("parse tasks.json: %w", err)
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks in this lesson")
	}
	return tasks, nil
}

func hintDialog(tasks []taskHintEntry) error {
	fmt.Println(ui.Bold.Render("Which task do you need help with?"))
	fmt.Println()
	for _, t := range tasks {
		fmt.Printf("  %s  %s\n", ui.Hint.Render(t.Slug), ui.Muted.Render(t.Title))
	}
	fmt.Println()
	fmt.Println(ui.Muted.Render("Type the slug to get a hint, e.g:"))
	fmt.Printf("  %s\n", ui.Hint.Render(fmt.Sprintf("lst hint %s", tasks[0].Slug)))
	return nil
}

func showHint(tasks []taskHintEntry, slug string, explicitLevel int) error {
	task := findTaskBySlug(tasks, slug)
	if task == nil {
		var slugs []string
		for _, t := range tasks {
			slugs = append(slugs, t.Slug)
		}
		return errorWithHint(
			fmt.Sprintf("task %q not found", slug),
			fmt.Sprintf("available: %s", fmt.Sprintf("%v", slugs)),
		)
	}
	if len(task.Hints) == 0 {
		fmt.Printf("%s %s\n", ui.Muted.Render("→"), "No hints available for this task.")
		return nil
	}

	meta, err := lesson.LoadMeta()
	if err != nil {
		meta = lesson.Meta{HintLevels: map[string]int{}}
	}

	currentLevel := meta.HintLevels[slug]
	level := currentLevel + 1
	if explicitLevel > 0 {
		level = explicitLevel
	}
	if level > len(task.Hints) {
		level = len(task.Hints)
	}
	if level < 1 {
		level = 1
	}

	hint := task.Hints[level-1]

	fmt.Printf("%s %s\n", ui.Bold.Render(ui.GlyphHint+" Hint"), ui.Muted.Render(fmt.Sprintf("(%d/%d):", level, len(task.Hints))))
	fmt.Println()
	fmt.Printf("%s  %s\n", ui.Muted.Render("→"), ui.Muted.Render(task.Title))
	fmt.Println()

	r, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(90))
	out, err := r.Render(hint)
	if err != nil {
		fmt.Println(hint)
	} else {
		fmt.Print(out)
	}

	if explicitLevel == 0 {
		meta.HintLevels[slug] = level
		_ = lesson.SaveMeta(meta)
	}

	if level < len(task.Hints) {
		fmt.Printf("\n%s\n", ui.Muted.Render(fmt.Sprintf("Need more details? Run: lst hint %s", slug)))
	} else {
		fmt.Println(ui.Muted.Render("\nThat was the last hint. Good luck!"))
	}
	return nil
}

func findTaskBySlug(tasks []taskHintEntry, slug string) *taskHintEntry {
	for i := range tasks {
		if tasks[i].Slug == slug {
			return &tasks[i]
		}
	}
	return nil
}

func init() { rootCmd.AddCommand(hintCmd) }
