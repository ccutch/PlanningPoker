package planner

import (
	"log"
	"net/http"
	db "planner/database"
	"planner/templates"
	"strconv"
	"time"
)

func init() {
	Routes.HandleFunc("POST /pod/{id}/start", handlePodStart)
	Routes.HandleFunc("POST /pod/{id}/vote", handlePodVote)
}

// handlePodStart handles the owner starting a new round of voting
func handlePodStart(w http.ResponseWriter, r *http.Request) {
	topic := db.GetNextTopic(r.PathValue("id"))
	if topic == nil {
		http.Error(w, "No next topic", http.StatusInternalServerError)
		return
	}
	go StartPoll(topic.PodID, topic.ID)
}

func StartPoll(podID, topicID string) {
	db.UpdatePodStatus(podID, "voting")
	db.Publish(podID, "content", "voting-content")

	time.Sleep(time.Second * 10)
	ClosePoll(podID, topicID)
}

func ClosePoll(podID, topicID string) {
	log.Println("Closing poll...")
	if t, err := db.GetTopic(topicID); err != nil || t.Status == "complete" {
		log.Println("topic => ", t, err)
		return
	}
	db.UpdatePodStatus(podID, "waiting")
	votes, err := db.SelectVotesForTopic(topicID)
	log.Println("getting votes")
	if err != nil || len(votes) == 0 {
		db.Publish(podID, "content", "voting-content/waiting")
		return
	}
	var total, count int
	for _, v := range votes {
		total += v.Choice - 1
		count += 1
	}
	log.Println("total and count", total, count)
	avg, rem := total/count, total%count
	if rem != 0 {
		avg += 1
	}
	log.Println("total and count", total, count)
	db.CompleteTopicWithResults(topicID, avg)
	db.Publish(podID, "topics", "voting-topics")
	db.Publish(podID, "content", "voting-content")
}

// handlePodVote handles players votes during a round of voting
func handlePodVote(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")
	choice, err := strconv.Atoi(r.URL.Query().Get("c"))
	if err != nil {
		http.Error(w, "Invalid choice: "+r.FormValue("choice"), http.StatusBadRequest)
		return
	}
	player := db.CurrentPlayer(r, podID)
	topic := db.GetNextTopic(podID)
	db.UpsertVoteForPlayer(topic.ID, player.ID, choice)
	if votes, err := db.SelectVotesForTopic(topic.ID); err == nil {
		players, _ := db.GetPlayersForPod(podID)
		log.Println("Votes", len(votes))
		log.Println("Players", len(players))
		if len(votes) == len(players) {
			ClosePoll(podID, topic.ID)
		}
	}
	props := getPlannerProps(r)
	templates.Render(w, "voting-content", props)
}
