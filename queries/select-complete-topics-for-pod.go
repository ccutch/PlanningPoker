package db

import (
	_ "embed"
	"log"
	"planner"

	"github.com/pkg/errors"
)

//go:embed select-complete-topics-for-pod.sql
var selectCompleteTopicsForPodSQL string

func GetLastTopic(podID string) *planner.Topic {
	t := planner.Topic{PodID: podID, Status: "complete"}
	row := db.QueryRow(selectCompleteTopicsForPodSQL, t.PodID)
	err := row.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		log.Println("error getting last topic: ", err)
		return nil
	}
	return &t
}

func GetCompleteTopicsForPod(podID string) ([]*planner.Topic, error) {
	var topics []*planner.Topic
	rows, err := db.Query(selectCompleteTopicsForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query topics for pod: "+podID)
	}
	for rows.Next() {
		t := planner.Topic{PodID: podID, Status: "complete"}
		err := rows.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse topics in: "+podID)
		}
		topics = append(topics, &t)
	}
	return topics, errors.Wrap(err, "failed to get topics for: "+podID)
}
