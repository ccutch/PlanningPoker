package planner

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*
var templates embed.FS

var (
	t, err = template.ParseFS(templates, "templates/*.html")
	views  = template.Must(t, err)
)

func render(w io.Writer, name string, data interface{}) {
	var buf bytes.Buffer
	if err := views.ExecuteTemplate(&buf, name, data); err != nil {
		log.Println("Error rendering:", err)
		return
	}
	fmt.Fprint(w, strings.ReplaceAll(buf.String(), "\n", ""))
}

type PlannerProps struct {
	Pod            *Pod
	CurrentPlayer  *Player
	NextTopic      *Topic
	LastTopic      *Topic
	CurrentChoice  int
	Players        []*Player
	UpcomingTopics []*Topic
	CompleteTopics []*Topic
	Error          error
}

func getPlannerProps(r *http.Request) (props PlannerProps) {
	id := r.PathValue("id")
	props.Pod, props.Error = GetPod(id)
	if props.Error != nil {
		return
	}
	props.NextTopic = GetNextTopic(props.Pod.ID)
	props.LastTopic = GetLastTopic(props.Pod.ID)
	props.Players, props.Error = GetPlayersForPod(props.Pod.ID)
	if props.Error != nil {
		return
	}
	props.UpcomingTopics, props.Error = GetUpcomingTopicsForPod(props.Pod.ID)
	if props.Error != nil {
		return
	}
	props.CompleteTopics, props.Error = GetCompleteTopicsForPod(props.Pod.ID)
	if props.Error != nil {
		return
	}
	props.CurrentPlayer = CurrentPlayer(r, id)
	if props.CurrentPlayer != nil && props.NextTopic != nil {
		vote, err := SelectVoteForPlayer(
			props.NextTopic.ID,
			props.CurrentPlayer.ID,
		)
		if err == nil {
			props.CurrentChoice = vote.Choice
		}
	}
	return
}
