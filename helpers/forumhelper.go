package helpers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/models"
	"github.com/jackc/pgx"
)

func CreateNewOrGetExistingForum(forum *models.Forum) (*models.Forum, bool, error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(
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
		forum.Slug, forum.Title, forum.User).
		Scan(&forum.User)

	if err != nil {
		sError := err.Error()
		// dirty hack with error code
		if sError[len(sError)-2] == '5' {
			var result string
			err = database.Instance.Pool.QueryRow(
				`
				SELECT slug, posts, threads, title, "user"
				FROM forums 
				WHERE slug = $1
				`,
				forum.Slug).Scan(
				&forum.Slug,
				&forum.Posts,
				&forum.Threads,
				&forum.Title,
				&result)
			if err != nil {
				return nil, false, err
			}
			return forum, false, nil // existing forum
		}
		return nil, false, nil // 404 user
	}
	forum.Posts = 0
	forum.Threads = 0
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, false, err
	}
	return forum, true, nil // 201 created
}

func GetForumInfoBySlug(slug *string) (*models.Forum, error) {
	findedForum := models.Forum{}
	err := database.Instance.Pool.QueryRow(
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
	return &findedForum, nil
}

func GetThreadsByForumSlug(slug *string, limit, desc, since []byte) (*models.Threads, bool, error) {
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
	fmt.Println(queryStringBuffer.String())
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Commit()

	var rows *pgx.Rows
	if len(since) != 0 {
		sinceTime, err := time.Parse("2006-01-02T15:04:05.000Z07:00", string(since))
		if err != nil {
			return nil, false, err
		}
		rows, err = tx.Query(queryStringBuffer.String(), slug, sinceTime)
	} else {
		rows, err = tx.Query(queryStringBuffer.String(), slug)
	}
	if err != nil {
		return nil, false, err
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
			return nil, false, err
		}
		currRowThread.Created = createdInTime.Format("2006-01-02T15:04:05.000Z")
		threads = append(threads, &currRowThread)
	}
	if len(threads) == 0 {
		cnt := 0
		if err = tx.QueryRow("SELECT COUNT(*) FROM forums WHERE slug = $1", slug).Scan(&cnt); err != nil {
			return nil, false, err
		}
		if cnt != 0 {
			return nil, true, err // search is empty, but forum exists
		}
		return nil, false, nil
	}
	rows.Close()
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, false, err
	}
	return &threads, false, nil
}
