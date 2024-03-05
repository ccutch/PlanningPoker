package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed select-players-for-pod.sql
var selectPlayersForPodSQL string

func GetPlayersForPod(podID string) ([]*planner.Player, error) {
	var players []*planner.Player
	rows, err := db.Query(selectPlayersForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query players for pod: "+podID)
	}
	for rows.Next() {
		p := planner.Player{PodID: podID}
		err := rows.Scan(&p.ID, &p.Name, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse players in: "+podID)
		}
		players = append(players, &p)
	}
	return players, errors.Wrap(err, "failed to get players for: "+podID)
}
