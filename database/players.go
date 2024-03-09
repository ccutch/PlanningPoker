package database

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Player struct {
	ID        string
	PodID     string
	Name      string
	Owner     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Player) StartSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     p.PodID,
		Value:    p.ID,
		Expires:  time.Now().Add(72 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})
}

func CurrentPlayer(r *http.Request, podID string) (p *Player) {
	cookie, err := r.Cookie(podID)
	if err != nil {
		log.Println("No session", err)
		return
	}
	p, err = GetPlayer(cookie.Value)
	if err != nil {
		log.Println("Error getting", err)
	}
	return p
}

//go:embed queries/insert-player.sql
var insertPlayerSQL string

func CreatePlayer(podID, name string, owner bool) (*Player, error) {
	p := Player{PodID: podID, Name: name, Owner: owner}
	row := db.QueryRow(insertPlayerSQL, genID(10), p.PodID, p.Name, p.Owner)
	err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return &p, errors.Wrap(err, "failed to create player in db")
}

//go:embed queries/select-player-by-id.sql
var selectPlayerByIDSQL string

func GetPlayer(id string) (*Player, error) {
	player := Player{ID: id}
	row := db.QueryRow(selectPlayerByIDSQL, player.ID)
	err := row.Scan(&player.PodID, &player.Name, &player.Owner, &player.CreatedAt, &player.UpdatedAt)
	return &player, errors.Wrap(err, "failed to get player: "+id)
}

//go:embed queries/select-players-for-pod.sql
var selectPlayersForPodSQL string

func GetPlayersForPod(podID string) ([]*Player, error) {
	var players []*Player
	rows, err := db.Query(selectPlayersForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query players for pod: "+podID)
	}
	for rows.Next() {
		p := Player{PodID: podID}
		err := rows.Scan(&p.ID, &p.Name, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse players in: "+podID)
		}
		players = append(players, &p)
	}
	return players, errors.Wrap(err, "failed to get players for: "+podID)
}
