package planner

import (
	"time"
)

type Pod struct {
	ID             string
	Name           string
	Private        bool
	Strategy       string
	Status         string
	UpcomingTopics []*Topic
	ResolvedTopics []*Topic
	Votes          map[string]string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Player struct {
	ID        string
	PodID     string
	Name      string
	Owner     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Topic struct {
	ID        string
	PodID     string
	Prompt    string
	Status    string
	Result    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Vote struct {
	TopicID   string
	PlayerID  string
	Choice    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Pod) VotingChoices() []string {
	return VotingStrategies[p.Strategy]
}

func (t Topic) ResultString() string {
	if t.Result == 0 {
		return ""
	}
	p, _ := GetPod(t.PodID)
	return p.VotingChoices()[t.Result]
}

var VotingStrategies = map[string][]string{
	"Fibinocci Numbers": {"1", "2", "3", "5", "8", "13"},
	"T-Shirt Sizes":     {"xs", "s", "m", "lg", "xl", "2xl"},
	"Sentiment":         {"Strongly Disagree", "Disagree", "Mild Disagree", "Mild Agree", "Agree", "Strongly Agree"},
	"Face Cards":        {"Jack", "Queen", "King", "Ace"},
}
