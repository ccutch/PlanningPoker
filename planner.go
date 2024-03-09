package planner

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	db "planner/database"
	"planner/templates"
	"strings"
)

func init() {
	Routes.HandleFunc("GET /pod/{id}", plannerPage)
	Routes.HandleFunc("GET /pod/{id}/events", streamPodEvents)
	Routes.HandleFunc("POST /pod/{id}/register", handleRegister)
	Routes.HandleFunc("POST /pod/{id}/topics", handleNewTopic)
}

type PlannerProps struct {
	Pod            *db.Pod
	CurrentPlayer  *db.Player
	NextTopic      *db.Topic
	LastTopic      *db.Topic
	CurrentChoice  int
	Players        []*db.Player
	UpcomingTopics []*db.Topic
	CompleteTopics []*db.Topic
	Error          error
}

func getPlannerProps(r *http.Request) (props PlannerProps) {
	id := r.PathValue("id")
	props.Pod, props.Error = db.GetPod(id)
	if props.Error != nil {
		return
	}
	props.NextTopic = db.GetNextTopic(props.Pod.ID)
	props.LastTopic = db.GetLastTopic(props.Pod.ID)
	props.Players, props.Error = db.GetPlayersForPod(props.Pod.ID)
	if props.Error != nil {
		return
	}
	props.UpcomingTopics, props.Error = db.GetUpcomingTopicsForPod(props.Pod.ID)
	if props.Error != nil {
		return
	}
	props.CompleteTopics, props.Error = db.GetCompleteTopicsForPod(props.Pod.ID)
	if props.Error != nil {
		return
	}
	props.CurrentPlayer = db.CurrentPlayer(r, props.Pod.ID)
	if props.CurrentPlayer != nil && props.NextTopic != nil {
		vote, err := db.SelectVoteForPlayer(
			props.NextTopic.ID,
			props.CurrentPlayer.ID,
		)
		if err == nil {
			props.CurrentChoice = vote.Choice
		}
	}
	return
}

// planner page renders whole html with normal planner props
func plannerPage(w http.ResponseWriter, r *http.Request) {
	props := getPlannerProps(r)
	log.Println("Props =", props.NextTopic)
	templates.Render(w, "planner.html", props)
}

// streamEvents stream Server Side Events to the client with html data
func streamPodEvents(w http.ResponseWriter, r *http.Request) {
	if player := db.CurrentPlayer(r, r.PathValue("id")); player == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "event: ping\ndata: \n\n")
	flusher.Flush()
	events, done := make(chan db.Event), make(chan bool)
	l, err := db.Subscribe(r.PathValue("id"), events, done)
	if err != nil {
		http.Error(w, "Failed to subscribe!", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case e := <-events:
			props := getPlannerProps(r)

			var buf bytes.Buffer
			templates.Render(&buf, e.Tmpl, props)
			data := strings.ReplaceAll(buf.String(), "\n", "")

			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Name, data)
			flusher.Flush()
		case <-r.Context().Done():
			// Unsubscribe(r.PathValue("id"), events)
			db.Unsubscribe(l, r.PathValue("id"))
			done <- true
			close(events)
			return
		}
	}
}

// register will create a new player for the game, attach to session, and push state pod
func handleRegister(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")
	player, err := db.CreatePlayer(podID, r.FormValue("name"), false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	player.StartSession(w)
	db.Publish(podID, "players", "pod-info/players")
	w.Header().Add("Hx-Refresh", "true")
	w.WriteHeader(http.StatusNoContent)
}

// handleNewTopic handles a new topic posted by the owner of a pod
func handleNewTopic(w http.ResponseWriter, r *http.Request) {
	log.Println("form value prompt", r.FormValue("prompt"))
	t, err := db.CreateTopic(r.PathValue("id"), r.FormValue("prompt"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := db.Publish(t.PodID, "topics", "voting-topics"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	go func() {
		pod, _ := db.GetPod(t.PodID)
		topics, _ := db.GetUpcomingTopicsForPod(t.PodID)
		if pod == nil || topics == nil {
			return
		}

		if pod.Status != "voting" && len(topics) == 1 {
			db.Publish(t.PodID, "content", "voting-content")
		}
	}()
}
