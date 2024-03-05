package planner

import (
	"log"
	"time"
)

func StartPoll(podID, topicID string) {
	UpdatePodStatus(podID, "voting")
	Publish(podID, "content", "voting-content/voting")

	time.Sleep(time.Second * 10)
	UpdatePodStatus(podID, "waiting")

	result, _ := ClosePoll(topicID)
	CompleteTopicWithResults(topicID, result)

	Publish(podID, "topics", "voting-topics")
	Publish(podID, "content", "voting-content/results")
}

func ClosePoll(topicID string) (int, error) {
	var total, count int
	votes, err := SelectVotesForTopic(topicID)
	if err != nil {
		log.Println("Error closing poll", err)
		return 0, err
	}
	for _, v := range votes {
		total += v.Choice
		count += 1
	}
	if count == 0 {
		return 0, nil
	}
	avg, rem := total/count, total%count
	if rem == 0 {
		return avg, nil
	}
	return avg + 1, nil
}
