package planner

import (
	"fmt"
	"net/http"
)

var Routes = http.NewServeMux()

func init() {
	Routes.HandleFunc("GET /{$}", homepage)
	Routes.HandleFunc("POST /create-pod", handleCreatePod)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	render(w, "homepage.html", nil)
}

func handleCreatePod(w http.ResponseWriter, r *http.Request) {
	n, s, p := r.FormValue("name"), r.FormValue("strat"), r.FormValue("private")
	pod, player, err := CreatePod(n, s, p == "on")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player.Attach(w)
	w.Header().Add("Hx-Redirect", fmt.Sprintf("/pod/%s", pod.ID))
	w.WriteHeader(202)
}
