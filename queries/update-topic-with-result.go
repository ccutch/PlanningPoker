package db

import (
	_ "embed"

	"github.com/pkg/errors"
)

//go:embed update-topic-with-result.sql
var updateTopicWithResultsSQL string

func CompleteTopicWithResults(id string, result int) error {
	_, err := db.Exec(updateTopicWithResultsSQL, id, result)
	return errors.Wrap(err, "failed to complete topic: "+id)
}
