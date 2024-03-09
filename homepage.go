package planner

import (
	"fmt"
	"net/http"
	db "planner/database"
	"planner/templates"
)

var Routes = http.NewServeMux()

func init() {
	Routes.HandleFunc("GET /{$}", homepage)
	Routes.HandleFunc("POST /create-pod", handleCreatePod)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "homepage.html", nil)
}

func handleCreatePod(w http.ResponseWriter, r *http.Request) {
	n, s, p := r.FormValue("name"), r.FormValue("strat"), r.FormValue("private")
	pod, player, err := db.CreatePod(n, s, p == "on")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player.StartSession(w)
	w.Header().Add("Hx-Redirect", fmt.Sprintf("/pod/%s", pod.ID))
	w.WriteHeader(202)
}
