package helpers

import (
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/models"
)

func CreateNewOrGetExistingUsers(user *models.User) (*models.Users, bool, error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	users := models.Users{}
	execResult, err := database.Instance.Pool.Exec(
		`
		INSERT
		INTO users (nickname, fullname, email, about) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT DO NOTHING
		`,
		user.Nickname, user.Fullname, user.Email, user.About)

	if err != nil {
		return nil, false, err
	}

	if execResult.RowsAffected() != 0 {
		users = append(users, user)
		return &users, true, nil
	}

	rows, err := database.Instance.Pool.Query(
		`
		SELECT nickname, fullname, email, about
		FROM users 
		WHERE nickname = $1 or email = $2
		`,
		user.Nickname, user.Email)
	defer rows.Close()

	for rows.Next() {
		existingUser := &models.User{}
		if err := rows.Scan(
			&existingUser.Nickname,
			&existingUser.Fullname,
			&existingUser.Email,
			&existingUser.About); err != nil {
			return nil, false, err
		} else {
			users = append(users, existingUser)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, false, err
	}
	return &users, false, nil
}

func GetUserByNickname(nickname string) (*models.User, error) {
	findedUser := models.User{}
	err := database.Instance.Pool.QueryRow(
		`
		SELECT nickname, fullname, email, about 
		FROM users 
		WHERE nickname = $1
		`,
		nickname).Scan(
		&findedUser.Nickname,
		&findedUser.Fullname,
		&findedUser.Email,
		&findedUser.About)
	if err != nil {
		return nil, err
	}
	return &findedUser, nil
}

func UpdateUserInfo(user *models.User) (bool, error) {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Commit()

	execResult, err := database.Instance.Pool.Exec(
		`
		UPDATE users
		SET fullname = $1, email = $2, about = $3
		WHERE nickname = $4
		`,
		user.Fullname, user.Email, user.About, user.Nickname)

	return execResult.RowsAffected() != 0, err
}
