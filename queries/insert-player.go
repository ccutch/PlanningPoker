package db

import (
	_ "embed"
	"planner"

	"github.com/pkg/errors"
)

//go:embed insert-player.sql
var insertPlayerSQL string

func InsertPlayer(podID, name string, owner bool) (*planner.Player, error) {
	p := planner.Player{PodID: podID, Name: name, Owner: owner}
	row := db.QueryRow(insertPlayerSQL, genID(10), p.PodID, p.Name, p.Owner)
	err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return &p, errors.Wrap(err, "failed to create player in db")
}
