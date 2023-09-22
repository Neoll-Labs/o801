/*
 license x
*/

package db

import (
	"context"
	"github.com/nelsonstr/o801/models"
)

type UserStorage struct {
	db DBInterface
}

var _ CRService[*models.User] = (*UserStorage)(nil)

func NewUserStorage(db DBInterface) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

// Create User with name
func (u *UserStorage) Create(_ context.Context, name string) (*models.User, error) {
	stmt, err := u.db.Prepare("INSERT INTO users (name) VALUES ($1)  RETURNING id")
	if err != nil {
		return &models.NilUser, err
	}
	defer func() { _ = stmt.Close() }()

	// Start a transaction
	tx, err := u.db.Begin()
	if err != nil {
		return &models.NilUser, err
	}

	var insertedID int64
	if err = tx.Stmt(stmt).QueryRow(name).Scan(&insertedID); err != nil {
		// Rollback the transaction if an error occurs
		_ = tx.Rollback()
		return &models.NilUser, err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit(); err != nil {
		return &models.NilUser, err
	}

	return &models.User{ID: insertedID, Name: name}, nil
}

// Get User with id
func (u *UserStorage) Get(_ context.Context, id int) (*models.User, error) {
	rows := u.db.QueryRow("select id, name from users where id = $1", id)

	user := &models.User{}
	err := rows.Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return &models.NilUser, nil
		}
		return &models.NilUser, err
	}

	return user, nil
}
