/*
 license x
*/

package repository

import (
	"context"
	"fmt"

	"github.com/nelsonstr/o801/api"
	"github.com/nelsonstr/o801/models"
)

type UserRepository struct {
	db DBInterface
}

var _ api.Repository[*models.User] = (*UserRepository)(nil)

func NewUserRepo(db DBInterface) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create User with name.
func (u *UserRepository) Create(_ context.Context, name string) (*models.User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return &models.NilUser, fmt.Errorf("begin tx error: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	prepare, err := tx.Prepare("INSERT INTO users (name) VALUES ($1)  RETURNING id")
	if err != nil {
		return &models.NilUser, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() { _ = prepare.Close() }()

	var insertedID int64
	if err := prepare.QueryRow(name).Scan(&insertedID); err != nil {
		return &models.NilUser, fmt.Errorf("db insert row: %w", err)
	}

	// Commit the transaction if everything is successful.
	if err := tx.Commit(); err != nil {
		return &models.NilUser, fmt.Errorf("db on commit tx: %w", err)
	}

	return &models.User{ID: insertedID, Name: name}, nil
}

// Get User with id.
func (u *UserRepository) Get(_ context.Context, id int) (*models.User, error) {
	rows := u.db.QueryRow("select id, name from users where id = $1", id)

	user := &models.User{}

	err := rows.Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return &models.NilUser, nil
		}

		return &models.NilUser, fmt.Errorf("db error feting user data: %w", err)
	}

	return user, nil
}
