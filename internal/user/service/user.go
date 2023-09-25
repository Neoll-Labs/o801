/*
 license x
*/

package service

import (
	"context"

	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/internal/interfaces"
	"github.com/nelsonstr/o801/internal/model"
)

type UserService struct {
	cache        interfaces.Cache[model.UserView]
	repository   interfaces.Repository[*model.User]
	errorHandler internal.ErrorHandler
}

var (
	_ interfaces.ServiceAPI[*model.UserView] = (*UserService)(nil)
)

func NewUserService(repo interfaces.Repository[*model.User], cache interfaces.Cache[model.UserView]) *UserService {
	return &UserService{
		cache:        cache,
		repository:   repo,
		errorHandler: internal.DefaultErrorHandler,
	}
}

func (s UserService) Get(ctx context.Context, usr *model.UserView) (*model.UserView, error) {
	if user, exists := s.cache.Get(usr.ID); exists {
		return &user, nil
	}

	mUser, err := s.repository.Get(ctx, &model.User{ID: usr.ID})
	if err != nil {
		return &model.NilUserView, &internal.StorageError{
			Err: err,
		}
	}

	if mUser == &model.NilUser {
		return &model.NilUserView, &internal.NotFoundError{}
	}

	chanUser := make(chan *model.UserView, 1)

	s.cacheSet(chanUser, mUser)

	return <-chanUser, nil
}

func (s UserService) Create(ctx context.Context, usr *model.UserView) (*model.UserView, error) {
	mUser, err := s.repository.Create(ctx, &model.User{Name: usr.Name})
	if err != nil {
		return &model.NilUserView, &internal.StorageError{
			Err: err,
		}
	}

	chanUser := make(chan *model.UserView, 1)

	s.cacheSet(chanUser, mUser)

	return <-chanUser, nil
}

func (s UserService) cacheSet(userChan chan *model.UserView, mUser *model.User) {
	user := &model.UserView{
		ID:   mUser.ID,
		Name: mUser.Name,
	}

	s.cache.Set(user.ID, user)

	userChan <- user
}
