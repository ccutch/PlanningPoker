package planner

import (
	"crypto/rand"
	"database/sql"
	_ "embed"
	"encoding/base64"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	dbURL = os.Getenv("DATABASE_URL")
	db    *sql.DB
)

func init() {
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect to db: "+dbURL))
	}

	m, err := migrate.New("github://ccutch/PlanningPoker/migrations", dbURL)
	if err != nil {
		log.Println("failed ot parse", err)
	} else {
		m.Up()
	}
}

func genID(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("failed to generate random secret")
	}
	return "x" + base64.URLEncoding.EncodeToString(b)[:n-1]
}

//go:embed queries/insert-pod.sql
var insertPodSQL string

func CreatePod(name, strategy string, private bool) (*Pod, *Player, error) {
	if _, ok := VotingStrategies[strategy]; !ok {
		return nil, nil, errors.Errorf("Invalid strategy: %s", strategy)
	}
	pod := Pod{Name: name, Strategy: strategy, Private: private}
	row := db.QueryRow(insertPodSQL, genID(10), pod.Name, pod.Strategy, pod.Private)
	if err = row.Scan(&pod.ID, &pod.CreatedAt, &pod.UpdatedAt); err != nil {
		return nil, nil, errors.Wrap(err, "failed to create pod in db")
	}
	player, err := CreatePlayer(pod.ID, "Pod Leader", true)
	return &pod, player, errors.Wrap(err, "failed to owner for pod")
}

//go:embed queries/select-pod-by-id.sql
var selectPodByIDSQL string

func GetPod(id string) (*Pod, error) {
	pod := Pod{ID: id}
	row := db.QueryRow(selectPodByIDSQL, pod.ID)
	err := row.Scan(&pod.Name, &pod.Strategy, &pod.Private, &pod.Status, &pod.CreatedAt, &pod.UpdatedAt)
	return &pod, errors.Wrap(err, "failed to get pod: "+id)
}

//go:embed queries/insert-player.sql
var insertPlayerSQL string

func CreatePlayer(podID, name string, owner bool) (*Player, error) {
	p := Player{PodID: podID, Name: name, Owner: owner}
	row := db.QueryRow(insertPlayerSQL, genID(10), p.PodID, p.Name, p.Owner)
	err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return &p, errors.Wrap(err, "failed to create player in db")
}

//go:embed queries/select-player-by-id.sql
var selectPlayerByIDSQL string

func GetPlayer(id string) (*Player, error) {
	player := Player{ID: id}
	row := db.QueryRow(selectPlayerByIDSQL, player.ID)
	err := row.Scan(&player.PodID, &player.Name, &player.Owner, &player.CreatedAt, &player.UpdatedAt)
	return &player, errors.Wrap(err, "failed to get player: "+id)
}

//go:embed queries/select-players-for-pod.sql
var selectPlayersForPodSQL string

func GetPlayersForPod(podID string) ([]*Player, error) {
	var players []*Player
	rows, err := db.Query(selectPlayersForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query players for pod: "+podID)
	}
	for rows.Next() {
		p := Player{PodID: podID}
		err := rows.Scan(&p.ID, &p.Name, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse players in: "+podID)
		}
		players = append(players, &p)
	}
	return players, errors.Wrap(err, "failed to get players for: "+podID)
}

//go:embed queries/insert-topic.sql
var insertTopicSQL string

func CreateTopic(podID, prompt string) (*Topic, error) {
	t := Topic{PodID: podID, Prompt: prompt}
	row := db.QueryRow(insertTopicSQL, genID(10), t.PodID, t.Prompt)
	err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	return &t, errors.Wrap(err, "failed to create topic for: "+podID)
}

//go:embed queries/select-upcoming-topics-for-pod.sql
var selectUpcomingTopicsForPodSQL string

func GetNextTopic(podID string) *Topic {
	t := Topic{PodID: podID, Status: "upcoming"}
	row := db.QueryRow(selectUpcomingTopicsForPodSQL, t.PodID)
	err := row.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil
	}
	return &t
}

func GetUpcomingTopicsForPod(podID string) ([]*Topic, error) {
	var topics []*Topic
	rows, err := db.Query(selectUpcomingTopicsForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query topics for pod: "+podID)
	}
	for rows.Next() {
		t := Topic{PodID: podID, Status: "upcoming"}
		err := rows.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse topics in: "+podID)
		}
		topics = append(topics, &t)
	}
	return topics, errors.Wrap(err, "failed to get topics for: "+podID)
}

//go:embed queries/select-complete-topics-for-pod.sql
var selectCompleteTopicsForPodSQL string

func GetLastTopic(podID string) *Topic {
	t := Topic{PodID: podID, Status: "complete"}
	row := db.QueryRow(selectCompleteTopicsForPodSQL, t.PodID)
	err := row.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		log.Println("error getting last topic: ", err)
		return nil
	}
	return &t
}

func GetCompleteTopicsForPod(podID string) ([]*Topic, error) {
	var topics []*Topic
	rows, err := db.Query(selectCompleteTopicsForPodSQL, podID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query topics for pod: "+podID)
	}
	for rows.Next() {
		t := Topic{PodID: podID, Status: "complete"}
		err := rows.Scan(&t.ID, &t.Prompt, &t.Result, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse topics in: "+podID)
		}
		topics = append(topics, &t)
	}
	return topics, errors.Wrap(err, "failed to get topics for: "+podID)
}

//go:embed queries/update-pod-status-by-id.sql
var updatePodStatusByIDSQL string

func UpdatePodStatus(id, status string) error {
	_, err := db.Exec(updatePodStatusByIDSQL, id, status)
	return errors.Wrap(err, "failed to update pod: "+id)
}

//go:embed queries/update-topic-with-result.sql
var updateTopicWithResultsSQL string

func CompleteTopicWithResults(id string, result int) error {
	_, err := db.Exec(updateTopicWithResultsSQL, id, result)
	return errors.Wrap(err, "failed to complete topic: "+id)
}

//go:embed queries/upsert-vote-for-player.sql
var upsertVoteForPlayerSQL string

func UpsertVoteForPlayer(topicID, playerID string, choice int) error {
	_, err := db.Exec(upsertVoteForPlayerSQL, topicID, playerID, choice)
	return errors.Wrap(err, "failed to upsert vote: "+topicID+" - "+playerID)
}

//go:embed queries/select-vote-for-player.sql
var selectVoteForPlayerSQL string

func SelectVoteForPlayer(topicID, playerID string) (*Vote, error) {
	vote := Vote{TopicID: topicID, PlayerID: playerID}
	row := db.QueryRow(selectVoteForPlayerSQL, topicID, playerID)
	err := row.Scan(&vote.Choice, &vote.CreatedAt, &vote.UpdatedAt)
	return &vote, errors.Wrap(err, "failed to select vote: "+topicID+" - "+playerID)
}

//go:embed queries/select-votes-for-topic.sql
var selectVotesForTopicSQL string

func SelectVotesForTopic(topicID string) ([]*Vote, error) {
	var votes []*Vote
	rows, err := db.Query(selectVotesForTopicSQL, topicID)
	if err != nil {
		return nil, errors.Wrap(err, "faield to query votes")
	}
	for rows.Next() {
		v := Vote{TopicID: topicID}
		err := rows.Scan(&v.PlayerID, &v.Choice, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse vote")
		}
		votes = append(votes, &v)
	}
	return votes, nil
}
