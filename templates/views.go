package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"planner"
	db "planner/queries"
	"strings"
)

var Routes = http.NewServeMux()

//go:embed *.html
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
	Pod            *planner.Pod
	CurrentPlayer  *planner.Player
	NextTopic      *planner.Topic
	LastTopic      *planner.Topic
	CurrentChoice  int
	Players        []*planner.Player
	UpcomingTopics []*planner.Topic
	CompleteTopics []*planner.Topic
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
	props.CurrentPlayer = planner.CurrentPlayer(r, id)
	if props.CurrentPlayer != nil {
		props.CurrentChoice = planner.Choice(id, props.CurrentPlayer.ID)
	}
	return
}
