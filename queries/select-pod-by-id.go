package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed select-pod-by-id.sql
var selectPodByIDSQL string

func GetPod(id string) (*planner.Pod, error) {
	pod := planner.Pod{ID: id}
	row := db.QueryRow(selectPodByIDSQL, pod.ID)
	err := row.Scan(&pod.Name, &pod.Strategy, &pod.Private, &pod.Status, &pod.CreatedAt, &pod.UpdatedAt)
	return &pod, errors.Wrap(err, "failed to get pod: "+id)
}
