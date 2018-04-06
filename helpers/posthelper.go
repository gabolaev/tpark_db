package helpers

import (
	"strconv"
	"time"

	"github.com/gabolaev/tpark_db/config"
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

	tx := database.StartTransaction()
	defer tx.Rollback()

	forumSlug, err := GetThreadForum(&threadID)
	if err != nil {
		return nil, errors.NotFoundError
	}
	currentTime := time.Now()
	currentTimeString := currentTime.Format(config.Instance.API.TimestampFormat)
	for _, post := range *posts {
		if post.Parent != 0 {
			rows, err := tx.Query("SELECT 1 FROM posts WHERE id = $1 AND thread = $2", post.Parent, threadID)
			if err != nil {
				return nil, err
			}
			if !rows.Next() {
				return nil, errors.ConflictError
			}
			rows.Close()
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
			return nil, errors.NotFoundError
		}
		post.Created = currentTimeString
		post.Edited = false
		post.Forum = forumSlug
		post.Thread = threadID
		err = IncrementCounters(&forumSlug, "posts")
		if err != nil {
			return nil, err
		}
	}
	database.CommitTransaction(tx)
	return posts, nil
}
