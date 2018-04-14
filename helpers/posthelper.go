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
		INSERT INTO posts (author, forum, created, message, parent, thread)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created
		`, post.Author, forumSlug, currentTime, post.Message, post.Parent, threadID).
			Scan(&post.ID, &currentTime)
		if err != nil {
			return nil, errors.NotFoundError
		}
		post.Created = currentTime.Format(config.Instance.API.TimestampFormat)
		post.IsEdited = false
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

func GetPostDetails(id *string) (*models.Post, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	post := models.Post{}

	var createdInTime time.Time
	err := tx.QueryRow(
		`
		SELECT 
			id,
			author,
			created,
			forum,
			edited,
			message,
			parent,
			thread
		FROM
			posts
		WHERE
			id = $1
		`, *id).Scan(
		&post.ID,
		&post.Author,
		&createdInTime,
		&post.Forum,
		&post.IsEdited,
		&post.Message,
		&post.Parent,
		&post.Thread)
	if err != nil {
		return nil, errors.NotFoundError
	}

	post.Created = createdInTime.Format(config.Instance.API.TimestampFormat)
	database.CommitTransaction(tx)
	return &post, nil
}

func GetPostFullDetails(id *string, relatedParams []string) (*models.PostFull, error) {

	postFull := models.PostFull{}
	var err error
	for _, value := range relatedParams {
		switch value {
		case "post":
			postFull.Post, err = GetPostDetails(id)
		case "thread":
			threadID := strconv.Itoa(postFull.Post.Thread)
			postFull.Thread, err = GetThreadDetailsBySlugOrID(&threadID)
		case "forum":
			forumSlug := postFull.Post.Forum
			postFull.Forum, err = GetForumDetailsBySlug(&forumSlug)
		case "user":
			userNickname := postFull.Post.Author
			postFull.Author, err = GetUserByNickname(&userNickname)
		}
		if err != nil {
			return nil, err
		}
	}
	return &postFull, nil
}

func UpdatePostDetails(id *string, postUpdate *models.PostUpdate) (*models.Post, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()
	post, err := GetPostDetails(id)
	if err != nil {
		return nil, errors.NotFoundError
	}
	if len(postUpdate.Message) != 0 && post.Message != postUpdate.Message {
		post.IsEdited = true
		post.Message = postUpdate.Message
	} else {
		return post, nil
	}
	_, err = tx.Exec(
		`
		UPDATE 
			posts
		SET
			message = $1, 
			edited = true
		WHERE 
			id = $2
		`, postUpdate.Message, post.ID)
	if err != nil {
		return nil, err
	}
	database.CommitTransaction(tx)
	return post, nil
}
