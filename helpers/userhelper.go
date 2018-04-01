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
	defer tx.Commit()

	users := models.Users{}
	execResult, err := tx.Exec(
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

	rows, err := tx.Query(
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
		tx.Rollback()
		return nil, false, nil
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

func UpdateUserInfo(user *models.User) error {
	tx, err := database.Instance.Pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = tx.QueryRow(
		`
		UPDATE users
		SET
			fullname = coalesce(coalesce(nullif($1, ''), fullname)), 
			email = coalesce(coalesce(nullif($2, ''), email)), 
			about = coalesce(coalesce(nullif($3, ''), about))
		WHERE
			nickname = $4
		RETURNING
			fullname,
			email,
			about
		`,
		user.Fullname, user.Email, user.About, user.Nickname).
		Scan(&user.Fullname, &user.Email, &user.About)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return err
}
