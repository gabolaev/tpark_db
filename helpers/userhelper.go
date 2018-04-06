package helpers

import (
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/models"
)

func UserExists(nickname *string) bool {
	tx := database.StartTransaction()
	rows, err := tx.Query("SELECT 1 FROM users WHERE nickname = $1", nickname)
	if err != nil {
		return false
	}
	if rows.Next() {
		rows.Close()
		database.CommitTransaction(tx)
		return true
	}
	return false
}

func CreateNewOrGetExistingUsers(user *models.User) (*models.Users, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

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
		return nil, err
	}

	if execResult.RowsAffected() != 0 {
		users = append(users, user)
		database.CommitTransaction(tx)
		return &users, nil
	}

	rows, err := tx.Query(
		`
		SELECT nickname, fullname, email, about
		FROM users 
		WHERE nickname = $1 or email = $2
		`,
		user.Nickname, user.Email)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		existingUser := &models.User{}
		if err := rows.Scan(
			&existingUser.Nickname,
			&existingUser.Fullname,
			&existingUser.Email,
			&existingUser.About); err != nil {
			return nil, err
		}
		users = append(users, existingUser)
	}
	database.CommitTransaction(tx)
	return &users, errors.ConflictError
}

func GetUserByNickname(nickname string) (*models.User, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	findedUser := models.User{}
	err := tx.QueryRow(
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
		return nil, errors.NotFoundError
	}
	database.CommitTransaction(tx)
	return &findedUser, nil
}

func UpdateUserInfo(user *models.User) error {
	tx := database.StartTransaction()
	defer tx.Rollback()

	err := tx.QueryRow(
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
		sError := err.Error()
		if sError[len(sError)-2] == '5' {
			return errors.ConflictError
		}
		return errors.NotFoundError
	}
	database.CommitTransaction(tx)
	return nil
}
