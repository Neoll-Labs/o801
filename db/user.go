/*
 license x
*/

package db

import (
	"context"
	"database/sql"
	"github.com/nelsonstr/o801/models"
)

type UserStorage struct {
	db *sql.DB
}

var _ CRService[*models.User] = (*UserStorage)(nil)

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

// Create User with name
func (u *UserStorage) Create(_ context.Context, name string) (*models.User, error) {
	// Prepare the SQL statement with placeholders
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

	var insertedID int

	if err = tx.Stmt(stmt).QueryRow(name).Scan(&insertedID); err != nil {
		// Rollback the transaction if an error occurs
		_ = tx.Rollback()
		return &models.NilUser, err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit(); err != nil {
		return &models.NilUser, err
	}

	return &models.User{ID: int64(insertedID), Name: name}, nil
}

// Get User with id
func (u *UserStorage) Get(_ context.Context, id int) (*models.User, error) {
	// Prepare the SQL statement with placeholders
	stmt, err := u.db.Prepare("select id, name from users where id = $1")
	if err != nil {
		return &models.NilUser, err
	}
	defer func() { _ = stmt.Close() }()

	// Start a transaction
	tx, err := u.db.Begin()
	if err != nil {
		return &models.NilUser, err
	}
	defer func() { _ = tx.Rollback() }()

	user := &models.User{}

	err = tx.Stmt(stmt).QueryRow(id).Scan(&user.ID, &user.Name)
	if err != nil {
		return &models.NilUser, err
	}

	return user, nil
}
