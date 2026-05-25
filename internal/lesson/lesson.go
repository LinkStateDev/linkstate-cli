// Package lesson reads the .linkstate.json metadata file that lst writes into
// each fetched lesson directory.
package lesson

import (
	"encoding/json"
	"fmt"
	"os"
)

const MetaFile = ".linkstate.json"

type Meta struct {
	LessonID int    `json:"lesson_id"`
	Slug     string `json:"slug"`
	Title    string `json:"title"`
}

// LoadMeta reads .linkstate.json from the current directory.
func LoadMeta() (Meta, error) {
	var m Meta
	data, err := os.ReadFile(MetaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return m, fmt.Errorf("%s not found. Run: lst fetch <slug>", MetaFile)
		}
		return m, fmt.Errorf("read %s: %w", MetaFile, err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("parse %s: %w", MetaFile, err)
	}
	return m, nil
}
