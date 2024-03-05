package db

import (
	_ "embed"

	"github.com/pkg/errors"
)

//go:embed update-pod-status-by-id.sql
var updatePodStatusByIDSQL string

func UpdatePodStatus(id, status string) error {
	_, err := db.Exec(updatePodStatusByIDSQL, id, status)
	return errors.Wrap(err, "failed to update pod: "+id)
}
