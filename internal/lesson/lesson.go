package lesson

import (
	"encoding/json"
	"fmt"
	"os"
)

const MetaFile = ".linkstate.json"

type Meta struct {
	LessonID   int            `json:"lesson_id"`
	TrackSlug  string         `json:"track_slug,omitempty"`
	Slug       string         `json:"slug"`
	Title      string         `json:"title"`
	HintLevels map[string]int `json:"hint_levels,omitempty"`
}

func LoadMeta() (Meta, error) {
	var m Meta
	data, err := os.ReadFile(MetaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return m, fmt.Errorf("%s not found. Run: lst start <lesson-slug>", MetaFile)
		}
		return m, fmt.Errorf("read %s: %w", MetaFile, err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("parse %s: %w", MetaFile, err)
	}
	if m.HintLevels == nil {
		m.HintLevels = map[string]int{}
	}
	return m, nil
}

func SaveMeta(m Meta) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(MetaFile, data, 0644)
}
