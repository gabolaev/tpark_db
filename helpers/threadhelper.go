package helpers

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/models"
)

func checkThreadSlugExisting(slug *string) (count int, err error) {
	tx := database.StartTransaction()
	defer tx.Rollback()
	_ = tx.QueryRow("SELECT 1 FROM threads WHERE slug = $1", slug).Scan(&count)
	database.CommitTransaction(tx)
	return
}

func GetThreadDetailsBySlugOrID(slugOrID *string) (*models.Thread, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	var fieldName string
	if IsNumber(slugOrID) {
		fieldName = "id"
	} else {
		fieldName = "slug"
	}

	createdInTime := time.Time{}
	thread := models.Thread{}
	err := tx.QueryRow(
		`
		SELECT id, slug, author, created AT TIME ZONE 'UTC', forum, message, title, votes
		FROM threads
		WHERE `+fieldName+` = $1
		`,
		*slugOrID).Scan(
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
		return nil, errors.NotFoundError
	}
	thread.Created = createdInTime.Format(config.Instance.API.TimestampFormat)
	database.CommitTransaction(tx)
	return &thread, nil
}

func GetThreadForum(tID *int) (slug string, err error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	err = tx.QueryRow("SELECT forum FROM threads WHERE id = $1", *tID).Scan(&slug)
	database.CommitTransaction(tx)
	return
}

func GetThreadIDBySlug(slug *string) (result int, err error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	err = tx.QueryRow("SELECT id FROM threads WHERE slug = $1", slug).Scan(&result)
	if err != nil {
		return -1, errors.NotFoundError
	}
	database.CommitTransaction(tx)
	return
}

func CreateNewOrGetExistingThread(thread *models.Thread) (*models.Thread, error) {
	nowTime := time.Now()
	tx := database.StartTransaction()
	defer tx.Rollback()

	if thread.Slug != "" {
		slugCounts, err := checkThreadSlugExisting(&thread.Slug)
		if err != nil {
			return nil, err
		}
		if slugCounts > 0 {
			existThread, err := GetThreadDetailsBySlugOrID(&thread.Slug)
			if err != nil {
				return nil, err
			}
			database.CommitTransaction(tx)
			return existThread, errors.ConflictError
		}
	}

	var err error
	var createdInTime time.Time
	if thread.Created == "" {
		createdInTime = nowTime
	} else {
		createdInTime, err = time.Parse(config.Instance.Database.TimestampFormat, thread.Created)
		if err != nil {
			return nil, err
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
			thread, err := GetThreadDetailsBySlugOrID(&thread.Slug)
			if err != nil {
				return nil, err
			}
			database.CommitTransaction(tx)
			return thread, errors.ConflictError
		}
		return nil, errors.NotFoundError
	}
	database.CommitTransaction(tx)
	err = IncrementCounters(&thread.Forum, "threads")
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func UpdateThreadDetails(slugOrID *string, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	var thread models.Thread

	tx := database.StartTransaction()
	defer tx.Rollback()
	var createdIntime time.Time
	var fieldName string
	if IsNumber(slugOrID) {
		fieldName = "id"
	} else {
		fieldName = "slug"
	}
	err := tx.QueryRow(
		`
		UPDATE threads
		SET
			message = coalesce(coalesce(nullif($2, ''), message)), 
			title = coalesce(coalesce(nullif($3, ''), title))
		WHERE
			`+fieldName+` = $1
		RETURNING
			author,
			created AT TIME ZONE 'UTC',
			forum,
			id,
			message,
			slug,
			title
		`, *slugOrID, threadUpdate.Message, threadUpdate.Title).Scan(
		&thread.Author,
		&createdIntime,
		&thread.Forum,
		&thread.ID,
		&thread.Message,
		&thread.Slug,
		&thread.Title)
	if err != nil {
		return nil, errors.NotFoundError
	}
	thread.Created = createdIntime.Format(config.Instance.API.TimestampFormat)
	database.CommitTransaction(tx)
	return &thread, nil
}

func GetThreadPostsFlat(slugOrID *string, limit, since, desc []byte) (*models.Posts, error) {
	var ID int
	if IsNumber(slugOrID) {
		ID, _ = strconv.Atoi(*slugOrID)
	} else {
		ID, _ = GetThreadIDBySlug(slugOrID)
	}
	posts := models.Posts{}
	queryStringBuffer := bytes.Buffer{}
	queryStringBuffer.WriteString(
		`SELECT 
			p.author,
			p.created,
			p.forum,
			p.id,
			p.Edited,
			p.message,
			p.parent,
			p.thread
		FROM posts p
		WHERE p.thread = $1`)
	sinceExists := lsdBuilder(&queryStringBuffer, limit, since, desc, "p.id", "p.id", false)
	fmt.Println(queryStringBuffer.String())
	tx := database.StartTransaction()
	defer tx.Rollback()

	var rows *pgx.Rows
	var err error
	fmt.Println(queryStringBuffer.String())
	if sinceExists {
		rows, err = tx.Query(queryStringBuffer.String(), ID, string(since))
	} else {
		rows, err = tx.Query(queryStringBuffer.String(), ID)
	}
	if err != nil {
		return nil, err // TODO this
	}
	var createdInTime time.Time
	for rows.Next() {
		currentPost := models.Post{}
		if err = rows.Scan(
			&currentPost.Author,
			&createdInTime,
			&currentPost.Forum,
			&currentPost.ID,
			&currentPost.IsEdited,
			&currentPost.Message,
			&currentPost.Parent,
			&currentPost.Thread,
		); err != nil {
			fmt.Println(err)
			return nil, err
		}
		currentPost.Created = createdInTime.Format(config.Instance.API.TimestampFormat)
		posts = append(posts, &currentPost)
	}
	rows.Close()

	database.CommitTransaction(tx)
	return &posts, nil
}
