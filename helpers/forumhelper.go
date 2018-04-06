package helpers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/models"
	"github.com/jackc/pgx"
)

func IncrementCounters(slug *string, fieldName string) error {
	tx := database.StartTransaction()
	defer tx.Rollback()

	_, err := tx.Exec(
		`
		UPDATE forums
		SET `+fieldName+` = `+fieldName+` + 1
		WHERE slug = $1
		`, slug)
	if err != nil {
		return err
	}
	database.CommitTransaction(tx)
	return nil
}

func GetForumBySlug(slug *string) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	forum := models.Forum{}
	err := tx.QueryRow(
		`
		SELECT slug, posts, threads, title, "user"
		FROM forums 
		WHERE slug = $1
		`,
		*slug).Scan(
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
		&forum.Title,
		&forum.User)
	if err != nil {
		return nil, err
	}
	database.CommitTransaction(tx)
	return &forum, nil
}

func CreateNewOrGetExistingForum(forum *models.Forum) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(
		`
		INSERT
		INTO forums (slug, title, "user") 
		VALUES ($1, $2, (SELECT nickname
						 FROM users 
						 WHERE nickname = $3
						 )
				)
		RETURNING "user"
		`,
		forum.Slug, forum.Title, forum.User)
	err := rows.Scan(&forum.User)
	if err != nil {
		sError := err.Error()
		if sError[len(sError)-2] == '5' {
			forum, err = GetForumBySlug(&forum.Slug)
			return forum, errors.ConflictError
		}
		return nil, errors.NotFoundError
	}
	forum.Posts = 0
	forum.Threads = 0
	database.CommitTransaction(tx)
	return forum, nil
}

func GetForumDetailsBySlug(slug *string) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	findedForum := models.Forum{}
	err := tx.QueryRow(
		`
		SELECT posts, slug, threads, title, "user" 
		FROM forums
		WHERE slug = $1
		`,
		slug).Scan(
		&findedForum.Posts,
		&findedForum.Slug,
		&findedForum.Threads,
		&findedForum.Title,
		&findedForum.User)
	if err != nil {
		return nil, err
	}
	database.CommitTransaction(tx)
	return &findedForum, nil
}

func GetThreadsByForumSlug(slug *string, limit, desc, since []byte) (*models.Threads, error) {
	var queryStringBuffer bytes.Buffer
	queryStringBuffer.WriteString(
		`
		SELECT author, created AT TIME ZONE 'UTC', forum, id, message, slug, title, votes
		FROM threads 
		WHERE forum = $1`)

	faseDescChecker := false
	if len(since) != 0 {
		sign := ">"
		if desc != nil && bytes.Equal([]byte("true"), desc) {
			faseDescChecker = true
			sign = "<"
		}
		queryStringBuffer.WriteString(fmt.Sprintf(" AND created %s= $2", sign))
	}

	queryStringBuffer.WriteString("\nORDER BY created")
	if faseDescChecker || desc != nil && bytes.Equal([]byte("true"), desc) {
		queryStringBuffer.WriteString(" DESC")
	}

	if limit != nil {
		queryStringBuffer.WriteString(fmt.Sprintf("\nLIMIT %s", limit))
	}
	tx := database.StartTransaction()
	defer tx.Rollback()

	var rows *pgx.Rows
	var err error
	if len(since) != 0 {
		sinceTime, err := time.Parse(config.Instance.Database.TimestampFormat, string(since))
		if err != nil {
			return nil, err
		}
		rows, err = tx.Query(queryStringBuffer.String(), slug, sinceTime)
	} else {
		rows, err = tx.Query(queryStringBuffer.String(), slug)
	}
	if err != nil {
		return nil, err
	}
	var threads models.Threads
	for rows.Next() {
		currRowThread := models.Thread{}
		var createdInTime time.Time
		if err = rows.Scan(
			&currRowThread.Author,
			&createdInTime,
			&currRowThread.Forum,
			&currRowThread.ID,
			&currRowThread.Message,
			&currRowThread.Slug,
			&currRowThread.Title,
			&currRowThread.Votes); err != nil {
			return nil, err
		}
		currRowThread.Created = createdInTime.Format(config.Instance.API.TimestampFormat)
		threads = append(threads, &currRowThread)
	}
	if len(threads) == 0 {
		var threadsCount int
		if err = tx.QueryRow("SELECT 1 FROM forums WHERE slug = $1", slug).Scan(&threadsCount); err != nil {
			return nil, errors.NotFoundError
		}
		return nil, errors.EmptySearchError
	}
	database.CommitTransaction(tx)
	return &threads, nil
}
