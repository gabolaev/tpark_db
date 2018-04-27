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

func emptyForumSearchOrNF(tx *pgx.Tx, slug *string) error {
	var exists int
	if err := tx.QueryRow("SELECT 1 FROM forums WHERE slug = $1", slug).Scan(&exists); err != nil {
		return errors.NotFoundError
	}
	return errors.EmptySearchError
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
			forum, err = GetForumDetailsBySlug(&forum.Slug)
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

	var findedForum models.Forum
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
	sinceExists := lsdBuilder(&queryStringBuffer, limit, since, desc, "created", "created", true)
	tx := database.StartTransaction()
	defer tx.Rollback()
	var rows *pgx.Rows
	var err error
	fmt.Println(queryStringBuffer.String())
	if sinceExists {
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
		var currRowThread models.Thread
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
		return nil, emptyForumSearchOrNF(tx, slug)
	}
	database.CommitTransaction(tx)
	return &threads, nil
}

func GetForumUsersBySlug(slug *string, limit, desc, since []byte) (*models.Users, error) {
	var queryStringBuffer bytes.Buffer
	queryStringBuffer.WriteString(
		`
		SELECT DISTINCT 
			u.nickname, 
			u.email, 
			u.fullname,
			u.about
		FROM users u
			LEFT JOIN threads t on u.nickname = t.author
			LEFT JOIN posts p on u.nickname = p.author
		WHERE (p.forum = $1 OR t.forum = $1) 
		`)

	sinceExists := lsdBuilder(&queryStringBuffer, limit, since, desc, "u.nickname", "u.nickname", false)
	var users models.Users
	tx := database.StartTransaction()
	defer tx.Rollback()
	var rows *pgx.Rows
	var err error
	println(queryStringBuffer.String())
	if sinceExists {
		rows, err = tx.Query(queryStringBuffer.String(), *slug, string(since))
	} else {
		rows, err = tx.Query(queryStringBuffer.String(), *slug)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var currUser models.User
		if err := rows.Scan(
			&currUser.Nickname,
			&currUser.Email,
			&currUser.Fullname,
			&currUser.About,
		); err != nil {
			return nil, err
		}
		users = append(users, &currUser)
	}
	rows.Close()
	if len(users) == 0 {
		return nil, emptyForumSearchOrNF(tx, slug)
	}
	database.CommitTransaction(tx)
	return &users, err
}
