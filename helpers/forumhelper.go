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

	_, err = database.Instance.Pool.Exec(
		`
		INSERT
		INTO forums (slug, title, creator) 
		VALUES ($1, $2, $3)
		`,
		forum.Slug, forum.Title, forum.Creator)

	if err != nil {
		sError := err.Error()
		// dirty hack with error code
		if sError[len(sError)-2] == '5' {
			err := database.Instance.Pool.QueryRow(
				`
				SELECT posts, threads, title, creator
				FROM forums 
				WHERE slug = $1
				`,
				forum.Slug).Scan(
				&forum.Posts,
				&forum.Threads,
				&forum.Title,
				&forum.Creator)
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
		return nil, false, err
	}
	return forum, true, nil // 201 created
}
