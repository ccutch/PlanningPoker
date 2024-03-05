package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed select-player-by-id.sql
var selectPlayerByIDSQL string

func GetPlayer(id string) (*planner.Player, error) {
	player := planner.Player{ID: id}
	row := db.QueryRow(selectPlayerByIDSQL, player.ID)
	err := row.Scan(&player.PodID, &player.Name, &player.Owner, &player.CreatedAt, &player.UpdatedAt)
	return &player, errors.Wrap(err, "failed to get player: "+id)
}
