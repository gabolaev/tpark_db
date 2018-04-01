package helpers

import (
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/models"
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
