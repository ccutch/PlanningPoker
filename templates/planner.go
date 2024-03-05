package views

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"planner"
	db "planner/queries"
	"strconv"
	"strings"
)

func init() {
	Routes.HandleFunc("GET /pod/{id}", plannerPage)
	Routes.HandleFunc("GET /pod/{id}/events", streamPodEvents)
	Routes.HandleFunc("POST /pod/{id}/register", handleRegister)
	Routes.HandleFunc("POST /pod/{id}/topics", handleNewTopic)
	Routes.HandleFunc("POST /pod/{id}/start", handlePodStart)
	Routes.HandleFunc("POST /pod/{id}/vote", handlePodVote)
}

// planner page renders whole html with normal planner props
func plannerPage(w http.ResponseWriter, r *http.Request) {
	props := getPlannerProps(r)
	log.Println("error:", props.Error)
	render(w, "planner.html", props)
}

// streamEvents stream Server Side Events to the client with html data
func streamPodEvents(w http.ResponseWriter, r *http.Request) {
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
	events, done := make(chan planner.Event), make(chan bool)
	l, err := planner.Subscribe(r.PathValue("id"), events, done)
	if err != nil {
		http.Error(w, "Failed to subscribe!", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case e := <-events:
			props := getPlannerProps(r)

			var buf bytes.Buffer
			render(&buf, e.Tmpl, props)
			data := strings.ReplaceAll(buf.String(), "\n", "")

			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Name, data)
			flusher.Flush()
		case <-r.Context().Done():
			// Unsubscribe(r.PathValue("id"), events)
			planner.Unsubscribe(l, r.PathValue("id"))
			done <- true
			close(events)
			return
		}
	}
}

// register will create a new player for the game, attach to session, and push state pod
func handleRegister(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")
	player, err := db.InsertPlayer(podID, r.FormValue("name"), false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	player.Attach(w)
	planner.Publish(podID, "players", "pod-info/players")
	w.Header().Add("Hx-Refresh", "true")
	w.WriteHeader(http.StatusNoContent)
}

// handleNewTopic handles a new topic posted by the owner of a pod
func handleNewTopic(w http.ResponseWriter, r *http.Request) {
	t, err := db.InsertTopic(r.PathValue("id"), r.FormValue("prompt"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := planner.Publish(t.PodID, "topics", "voting-topics"); err != nil {
		log.Println("Error creating topic")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	go func() {
		pod, _ := db.GetPod(t.PodID)
		topics, _ := db.GetUpcomingTopicsForPod(t.PodID)
		log.Println(pod, topics)
		if pod == nil || topics == nil {
			return
		}

		log.Println(pod.Status, topics)
		if pod.Status != "voting" && len(topics) == 1 {
			planner.Publish(t.PodID, "content", "voting-content")
		}
	}()
}

// handlePodStart handles the owner starting a new round of voting
func handlePodStart(w http.ResponseWriter, r *http.Request) {
	topic := db.GetNextTopic(r.PathValue("id"))
	if topic == nil {
		http.Error(w, "No topic: "+err.Error(), http.StatusInternalServerError)
		return
	}
	go planner.StartPoll(topic.PodID, topic.ID)
}

// handlePodVote handles players votes during a round of voting
func handlePodVote(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")
	choice, err := strconv.Atoi(r.URL.Query().Get("c"))
	if err != nil {
		http.Error(w, "Invalid choice: "+r.FormValue("choice"), http.StatusBadRequest)
		return
	}
	player := planner.CurrentPlayer(r, podID)
	planner.Vote(podID, player.ID, choice)
	props := getPlannerProps(r)
	render(w, "voting-content/answered", props)
}
