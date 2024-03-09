package database

import (
	_ "embed"
	"time"

	"github.com/pkg/errors"
)

type Vote struct {
	TopicID   string
	PlayerID  string
	Choice    int
	CreatedAt time.Time
	UpdatedAt time.Time
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
