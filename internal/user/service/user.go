package service

import (
	"context"
	"github.com/nelsonstr/o801/api"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/internal/model"
	"sync"
)

type userService struct {
	mutex        sync.Mutex
	userCache    map[int64]User
	repository   api.Repository[*model.User]
	errorHandler internal.ErrorHandler
}

var (
	_ api.ServiceAPI[*User] = (*userService)(nil)
)

func NewUserService(repo api.Repository[*model.User]) *userService {
	return &userService{
		mutex:        sync.Mutex{},
		userCache:    make(map[int64]User),
		repository:   repo,
		errorHandler: internal.DefaultErrorHandler,
	}
}

func (s userService) Get(ctx context.Context, usr *User) (*User, error) {
	if user := s.GetFromCache(usr.ID); user != NilUser {

		return &user, nil
	}

	mUser, err := s.repository.Fetch(ctx, &model.User{ID: usr.ID})
	if err != nil {
		return &NilUser, err
	}

	if mUser == &model.NilUser {
		return &NilUser, &internal.NotFoundError{}
	}

	user := &User{
		ID:   mUser.ID,
		Name: mUser.Name,
	}

	s.AddToCache(*user)

	return user, nil
}

func (s userService) Create(ctx context.Context, usr *User) (*User, error) {

	mUser, err := s.repository.Create(ctx, &model.User{Name: usr.Name})
	if err != nil {

		return &NilUser, err
	}

	user := &User{
		ID:   mUser.ID,
		Name: mUser.Name,
	}

	s.AddToCache(*user)

	return user, nil
}

func (s *userService) AddToCache(user User) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.userCache[user.ID] = user
}

func (s *userService) GetFromCache(id int64) User {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	if u, exist := s.userCache[id]; exist {
		return u
	}

	return NilUser
}
