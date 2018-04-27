package helpers

import (
	"strconv"
	"time"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/models"
)

func UpdateOrAddThreadVote(vote *models.Vote, coeff int) (*models.Thread, error) {
	var thread models.Thread
	tx := database.StartTransaction()
	defer tx.Rollback()
	var createdInTime time.Time
	err := tx.QueryRow(
		`
			UPDATE threads
			SET votes = votes + ( $1::SMALLINT * $2::SMALLINT )
			WHERE id = $3
			RETURNING id, author, created AT TIME ZONE 'UTC', forum, message, slug, title, votes
			`, coeff, vote.Voice, vote.Thread).
		Scan(
			&thread.ID,
			&thread.Author,
			&createdInTime,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
	if err != nil {
		return nil, errors.NotFoundError
	}
	thread.Created = createdInTime.Format(config.Instance.API.TimestampFormat)
	database.CommitTransaction(tx)
	return &thread, nil
}

func VoteThread(slugOrID *string, vote *models.Vote) (*models.Thread, error) {
	var err error
	if vote.Voice != -1 && vote.Voice != 1 {
		return nil, errors.WrongParamsError
	}
	if IsNumber(slugOrID) {
		if vote.Thread, err = strconv.Atoi(*slugOrID); err != nil {
			return nil, err
		}
	} else {
		if vote.Thread, err = GetThreadIDBySlug(slugOrID); err != nil {
			return nil, errors.NotFoundError
		}
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	oldVoice := int16(0)
	rows, err := tx.Query(
		`
		SELECT voice
		FROM votes
		WHERE thread = $1 AND nickname = $2
		`, vote.Thread, vote.Nickname)

	if rows.Next() {
		err = rows.Scan(&oldVoice)
		rows.Close()
	} else {
		if !UserExists(&vote.Nickname) {
			return nil, errors.NotFoundError
		}
	}
	var thread *models.Thread

	switch vote.Voice {
	case oldVoice:
		str := strconv.Itoa(vote.Thread)
		thread, err = GetThreadDetailsBySlugOrID(&str)
		if err != nil {
			return nil, err
		}
	default:
		if oldVoice == 0 {
			thread, err = UpdateOrAddThreadVote(vote, 1)
		} else {
			thread, err = UpdateOrAddThreadVote(vote, 2)
		}
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(
			`
			INSERT INTO 
			votes (thread, nickname, voice) 
			VALUES ($1, $2, $3)
			ON CONFLICT (thread, nickname) DO UPDATE
			SET voice = excluded.voice
			`, vote.Thread, vote.Nickname, vote.Voice)
		if err != nil {
			return nil, err
		}
	}
	database.CommitTransaction(tx)
	return thread, nil
}
