package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed select-upcoming-topics-for-pod.sql
var selectUpcomingTopicsForPodSQL string

func GetNextTopic(podID string) *planner.Topic {
	t := planner.Topic{PodID: podID, Status: "upcoming"}
	row := db.QueryRow(selectUpcomingTopicsForPodSQL, t.PodID)
	err := row.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil
	}
	return &t
}

func GetUpcomingTopicsForPod(podID string) ([]*planner.Topic, error) {
	var topics []*planner.Topic
	rows, err := db.Query(selectUpcomingTopicsForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query topics for pod: "+podID)
	}
	for rows.Next() {
		t := planner.Topic{PodID: podID, Status: "upcoming"}
		err := rows.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse topics in: "+podID)
		}
		topics = append(topics, &t)
	}
	return topics, errors.Wrap(err, "failed to get topics for: "+podID)
}
