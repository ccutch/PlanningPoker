package database

import (
	_ "embed"
	"time"

	"github.com/pkg/errors"
)

type Topic struct {
	ID        string
	PodID     string
	Prompt    string
	Status    string
	Result    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Topic) ResultString() string {
	p, _ := GetPod(t.PodID)
	return p.VotingChoices()[t.Result]
}

//go:embed queries/insert-topic.sql
var insertTopicSQL string

func CreateTopic(podID, prompt string) (*Topic, error) {
	t := Topic{PodID: podID, Prompt: prompt}
	row := db.QueryRow(insertTopicSQL, genID(10), t.PodID, t.Prompt)
	err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	return &t, errors.Wrap(err, "failed to create topic for: "+podID)
}

//go:embed queries/select-topic-by-id.sql
var selectTopicByIDSQL string

func GetTopic(id string) (*Topic, error) {
	t := Topic{ID: id}
	row := db.QueryRow(selectTopicByIDSQL, id)
	err := row.Scan(&t.PodID, &t.Prompt, &t.Status,
		&t.Result, &t.CreatedAt, &t.UpdatedAt)
	return &t, errors.Wrap(err, "failed to get topic: "+id)
}

//go:embed queries/select-upcoming-topics-for-pod.sql
var selectUpcomingTopicsForPodSQL string

func GetNextTopic(podID string) *Topic {
	t := Topic{PodID: podID, Status: "upcoming"}
	row := db.QueryRow(selectUpcomingTopicsForPodSQL, t.PodID)
	err := row.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil
	}
	return &t
}

func GetUpcomingTopicsForPod(podID string) ([]*Topic, error) {
	var topics []*Topic
	rows, err := db.Query(selectUpcomingTopicsForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query topics for pod: "+podID)
	}
	for rows.Next() {
		t := Topic{PodID: podID, Status: "upcoming"}
		err := rows.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse topics in: "+podID)
		}
		topics = append(topics, &t)
	}
	return topics, errors.Wrap(err, "failed to get topics for: "+podID)
}

//go:embed queries/select-complete-topics-for-pod.sql
var selectCompleteTopicsForPodSQL string

func GetLastTopic(podID string) *Topic {
	t := Topic{PodID: podID, Status: "complete"}
	row := db.QueryRow(selectCompleteTopicsForPodSQL, t.PodID)
	err := row.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil
	}
	return &t
}

func GetCompleteTopicsForPod(podID string) ([]*Topic, error) {
	var topics []*Topic
	rows, err := db.Query(selectCompleteTopicsForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query topics for pod: "+podID)
	}
	for rows.Next() {
		t := Topic{PodID: podID, Status: "complete"}
		err := rows.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse topics in: "+podID)
		}
		topics = append(topics, &t)
	}
	return topics, errors.Wrap(err, "failed to get topics for: "+podID)
}

//go:embed queries/update-topic-with-result.sql
var updateTopicWithResultsSQL string

func CompleteTopicWithResults(id string, result int) error {
	_, err := db.Exec(updateTopicWithResultsSQL, id, "complete", result)
	return errors.Wrap(err, "failed to complete topic: "+id)
}

func RequeueTopicWithoutResults(id string) error {
	_, err := db.Exec(updateTopicWithResultsSQL, id, "upcoming", 0)
	return errors.Wrap(err, "failed to complete topic: "+id)
}
