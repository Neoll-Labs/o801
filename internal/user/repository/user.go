/*
 license x
*/

package repository

import (
	"context"
	"fmt"
	"github.com/nelsonstr/o801/api"
	"github.com/nelsonstr/o801/internal/model"
	"github.com/nelsonstr/o801/internal/repository"
)

type UserRepository struct {
	db repository.DBInterface
}

var _ api.Repository[*model.User] = (*UserRepository)(nil)

func NewUserRepository(db repository.DBInterface) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create User with name.
func (u *UserRepository) Create(_ context.Context, user *model.User) (*model.User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return &model.NilUser, fmt.Errorf("begin tx error: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	prepare, err := tx.Prepare("INSERT INTO users (name) VALUES ($1)  RETURNING id,name")
	if err != nil {
		return &model.NilUser, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() { _ = prepare.Close() }()

	newUser := &model.User{}
	if err := prepare.QueryRow(user.Name).Scan(&newUser.ID, &newUser.Name); err != nil {
		return &model.NilUser, fmt.Errorf("db insert row: %w", err)
	}

	// Commit the transaction if everything is successful.
	if err := tx.Commit(); err != nil {
		return &model.NilUser, fmt.Errorf("db on commit tx: %w", err)
	}

	return newUser, nil
}

// Fetch User with id.
func (u *UserRepository) Fetch(_ context.Context, usr *model.User) (*model.User, error) {
	rows := u.db.QueryRow("select id, name from users where id = $1", usr.ID)

	user := &model.User{}

	err := rows.Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return &model.NilUser, nil
		}

		return &model.NilUser, fmt.Errorf("db error feting user data: %w", err)
	}

	return user, nil
}
