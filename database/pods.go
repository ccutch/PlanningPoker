package database

import (
	_ "embed"
	"time"

	"github.com/pkg/errors"
)

type Pod struct {
	ID             string
	Name           string
	Private        bool
	Strategy       string
	Status         string
	UpcomingTopics []*Topic
	ResolvedTopics []*Topic
	Votes          map[string]string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (p Pod) VotingChoices() []string {
	return VotingStrategies[p.Strategy]
}

//go:embed queries/insert-pod.sql
var insertPodSQL string

func CreatePod(name, strategy string, private bool) (*Pod, *Player, error) {
	if _, ok := VotingStrategies[strategy]; !ok {
		return nil, nil, errors.Errorf("Invalid strategy: %s", strategy)
	}
	pod := Pod{Name: name, Strategy: strategy, Private: private}
	row := db.QueryRow(insertPodSQL, genID(10), pod.Name, pod.Strategy, pod.Private)
	if err := row.Scan(&pod.ID, &pod.CreatedAt, &pod.UpdatedAt); err != nil {
		return nil, nil, errors.Wrap(err, "failed to create pod in db")
	}
	player, err := CreatePlayer(pod.ID, "Pod Leader", true)
	return &pod, player, errors.Wrap(err, "failed to owner for pod")
}

//go:embed queries/select-pod-by-id.sql
var selectPodByIDSQL string

func GetPod(id string) (*Pod, error) {
	pod := Pod{ID: id}
	row := db.QueryRow(selectPodByIDSQL, pod.ID)
	err := row.Scan(&pod.Name, &pod.Strategy, &pod.Private, &pod.Status, &pod.CreatedAt, &pod.UpdatedAt)
	return &pod, errors.Wrap(err, "failed to get pod: "+id)
}

//go:embed queries/update-pod-status-by-id.sql
var updatePodStatusByIDSQL string

func UpdatePodStatus(id, status string) error {
	_, err := db.Exec(updatePodStatusByIDSQL, id, status)
	return errors.Wrap(err, "failed to update pod: "+id)
}
