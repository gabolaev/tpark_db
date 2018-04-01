package helpers

import (
	"time"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/models"
)

func CreateNewOrGetExistingThread(thread *models.Thread) (*models.Thread, bool, error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Commit()
	createdInTime, err := time.Parse("2006-01-02T15:04:05.000Z07:00", thread.Created)
	if err != nil {
		return nil, false, err
	}
	_, err = database.Instance.Pool.Exec(
		`
		INSERT
			INTO threads (slug, forum, author, created, message, title) 
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		thread.Slug,
		thread.Forum,
		thread.Author,
		createdInTime,
		thread.Message,
		thread.Title)

	if err != nil {
		sError := err.Error()
		// dirty hack with error code
		if sError[len(sError)-2] == '5' {
			err = database.Instance.Pool.QueryRow(
				`
				SELECT id, author, created, forum, message, title, votes
				FROM threads
				WHERE slug = $1
				`,
				thread.Slug).Scan(
				&thread.ID,
				&thread.Author,
				&createdInTime,
				&thread.Forum,
				&thread.Message,
				&thread.Title,
				&thread.Votes,
			)
			if err != nil {
				return nil, false, nil
			}
			thread.Created = createdInTime.Format("2006-01-02T15:04:05.000Z")
			return thread, false, err // existing forum
		}
		return nil, false, nil // 404 user
	}
	if err := tx.Commit(); err != nil {
		return nil, false, err
	}
	return thread, true, nil // 201 created
}
