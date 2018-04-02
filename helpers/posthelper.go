package helpers

import (
	"strconv"
	"time"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/models"
)

func CreatePostsByThreadSlugOrID(posts *models.Posts, slugOrID *string) (*models.Posts, error) {
	var threadID int
	var err error
	if IsNumber(slugOrID) {
		if threadID, err = strconv.Atoi(*slugOrID); err != nil {
			return nil, err
		}
	} else {
		if threadID, err = GetThreadIDBySlug(slugOrID); err != nil {
			return nil, err
		}
	}

	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	forumSlug, err := GetForumSlugByThreadID(&threadID)
	if err != nil {
		return nil, errors.NotFoundError
	}
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02T15:04:05.000Z")
	parentExists := 0
	for _, post := range *posts {
		if post.Parent != 0 {
			err = tx.QueryRow("SELECT COUNT(*) FROM posts WHERE id = $1", post.Parent).Scan(&parentExists)
			if err != nil {
				return nil, err
			}
			if parentExists != 1 {
				return nil, errors.ConflictError
			}
		}
		err = tx.QueryRow(
			`
		INSERT
		INTO posts (author, forum, created, message, parent, thread)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
		`, post.Author, forumSlug, currentTime, post.Message, post.Parent, threadID).
			Scan(&post.ID)
		if err != nil {
			return nil, errors.ConflictError
		}
		post.Created = currentTimeString
		post.Edited = false
		post.Forum = forumSlug
		post.Thread = threadID
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return posts, nil
}
