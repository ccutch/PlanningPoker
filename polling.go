package planner

import (
	"time"
)

var polls = map[string]map[string]int{}

func StartPoll(podID, topicID string) {
	polls[podID] = map[string]int{}
	UpdatePodStatus(podID, "voting")
	Publish(podID, "content", "voting-content/voting")

	time.Sleep(time.Second * 10)
	UpdatePodStatus(podID, "waiting")

	result := ClosePoll(podID)
	CompleteTopicWithResults(topicID, result)

	Publish(podID, "topics", "voting-topics")
	Publish(podID, "content", "voting-content/results")
}

func Vote(podID, playerID string, choice int) {
	polls[podID][playerID] = choice
}

func Choice(podID, playerID string) int {
	if v, ok := polls[podID]; ok {
		return v[playerID]
	}
	return 0
}

func ClosePoll(podID string) int {
	var total, count int
	for _, n := range polls[podID] {
		total += n
		count += 1
	}
	if count == 0 {
		return 0
	}
	avg, rem := total/count, total%count
	// We are ceiling; if we want to floor we can remove this check and `+ 1`
	if rem == 0 {
		return avg
	}
	return avg + 1
}
