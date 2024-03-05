package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed insert-topic.sql
var insertTopicSQL string

func InsertTopic(podID, prompt string) (*planner.Topic, error) {
	t := planner.Topic{PodID: podID, Prompt: prompt}
	row := db.QueryRow(insertTopicSQL, genID(10), t.PodID, t.Prompt)
	err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	return &t, errors.Wrap(err, "failed to create topic for: "+podID)
}
