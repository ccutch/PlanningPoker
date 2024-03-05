package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed insert-pod.sql
var insertPodSQL string

func InsertPod(name, strategy string, private bool) (*planner.Pod, *planner.Player, error) {
	if _, ok := VotingStrategies[strategy]; !ok {
		return nil, nil, errors.Errorf("Invalid strategy: %s", strategy)
	}
	pod := planner.Pod{Name: name, Strategy: strategy, Private: private}
	row := db.QueryRow(insertPodSQL, genID(10), pod.Name, pod.Strategy, pod.Private)
	if err := row.Scan(&pod.ID, &pod.CreatedAt, &pod.UpdatedAt); err != nil {
		return nil, nil, errors.Wrap(err, "failed to create pod in db")
	}
	player, err := InsertPlayer(pod.ID, "Pod Leader", true)
	return &pod, player, errors.Wrap(err, "failed to owner for pod")
}
