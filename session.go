package planner

import (
	"log"
	"net/http"
	"time"
)

func (player *Player) Attach(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     player.PodID,
		Value:    player.ID,
		Expires:  time.Now().Add(72 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})
}

func CurrentPlayer(r *http.Request, podID string) *Player {
	if cookie, err := r.Cookie(podID); err == nil {
		if player, err := GetPlayer(cookie.Value); err == nil {
			return player
		} else {
			log.Println("Error hetting", err)
		}
	} else {
		log.Println("failed", err)
	}
	return nil
}
