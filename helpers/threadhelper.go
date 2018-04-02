package helpers

import (
	"fmt"
	"time"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/models"
)

func checkThreadSlugExisting(slug *string) (count int, err error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	_ = tx.QueryRow("SELECT COUNT(*) FROM threads WHERE slug = $1", slug).Scan(&count)
	tx.Commit()
	return
}

func GetThreadBySlug(slug *string) (*models.Thread, error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	createdInTime := time.Time{}
	thread := models.Thread{}
	err = tx.QueryRow(
		`
		SELECT id, slug, author, created AT TIME ZONE 'UTC', forum, message, title, votes
		FROM threads
		WHERE slug = $1
		`,
		*slug).Scan(
		&thread.ID,
		&thread.Slug,
		&thread.Author,
		&createdInTime,
		&thread.Forum,
		&thread.Message,
		&thread.Title,
		&thread.Votes,
	)
	if err != nil {
		return nil, err
	}
	thread.Created = createdInTime.Format("2006-01-02T15:04:05.000Z")
	return &thread, nil
}

func GetForumSlugByThreadID(tID *int) (slug string, err error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT forum FROM threads WHERE id = $1", *tID).Scan(&slug)
	fmt.Println(err)
	tx.Commit()
	return
}

func GetThreadIDBySlug(slug *string) (result int, err error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT id FROM threads WHERE slug = $1", slug).Scan(&result)
	tx.Commit()
	return
}

func CreateNewOrGetExistingThread(thread *models.Thread) (*models.Thread, int16, error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, 500, err
	}
	defer tx.Commit()

	if thread.Slug != "" {
		slugCounts, err := checkThreadSlugExisting(&thread.Slug)
		if err != nil {
			return nil, 500, err
		}
		if slugCounts > 0 {
			existThread, err := GetThreadBySlug(&thread.Slug)
			if err != nil {
				return nil, 500, err
			}
			return existThread, 409, nil
		}
	}

	var createdInTime time.Time
	if thread.Created == "" {
		createdInTime = time.Now()
	} else {
		createdInTime, err = time.Parse("2006-01-02T15:04:05.000Z07:00", thread.Created)
		if err != nil {
			return nil, 500, err
		}
	}

	err = tx.QueryRow(
		`
		INSERT
			INTO threads (slug, forum, author, created, message, title) 
		VALUES ($1, (SELECT slug FROM forums WHERE slug = $2), $3, $4, $5, $6)
		RETURNING id, forum
		`,
		thread.Slug,
		thread.Forum,
		thread.Author,
		createdInTime,
		thread.Message,
		thread.Title).Scan(&thread.ID, &thread.Forum)

	if err != nil {
		sError := err.Error()
		if sError[len(sError)-2] == '5' {
			thread, err := GetThreadBySlug(&thread.Slug)
			if err != nil {
				return nil, 501, err
			}
			return thread, 409, nil
		}
		return nil, 404, nil
	}
	if err := tx.Commit(); err != nil {
		return nil, 500, err
	}
	return thread, 201, nil // 201 created
}
