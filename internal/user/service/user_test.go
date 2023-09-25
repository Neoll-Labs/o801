/*
 license x
*/

package service

import (
	"context"
	"errors"
	"github.com/nelsonstr/o801/internal/interfaces"
	user "github.com/nelsonstr/o801/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUserService(t *testing.T) {
	cache := &mockCache[user.UserView]{}
	s := NewUserService(repoMock{}, cache)

	assert.NotNil(t, s.cache)
	assert.NotNil(t, s.repository)
	assert.NotNil(t, s.errorHandler)
}

func TestUserService_Get_SuccessfulNotCached(t *testing.T) {
	t.Parallel()

	repo := repoMock{
		user: &user.User{
			ID:   1,
			Name: "nelson",
		},
	}

	cache := &mockCache[user.UserView]{}
	s := NewUserService(repo, cache)

	u, err := s.Get(context.Background(), &user.UserView{ID: 1})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserService_Get_SuccessfulCached(t *testing.T) {
	t.Parallel()

	repo := repoMock{}

	cache := &mockCache[user.UserView]{
		user: user.UserView{
			ID:   1,
			Name: "nelson",
		},
		exists: true,
	}
	s := NewUserService(repo, cache)

	u, err := s.Get(context.Background(), &user.UserView{ID: 1})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserService_Get_RepoError(t *testing.T) {
	t.Parallel()

	repo := repoMock{
		err: errors.New("repo error"),
	}

	cache := &mockCache[user.UserView]{}
	s := NewUserService(repo, cache)

	u, err := s.Get(context.Background(), &user.UserView{ID: 1})

	assert.Error(t, err)
	assert.Equal(t, &user.NilUserView, u)
}

func TestUserService_Get_RepoUserNotFound(t *testing.T) {
	t.Parallel()

	repo := repoMock{
		user: &user.NilUser,
	}

	cache := &mockCache[user.UserView]{}
	s := NewUserService(repo, cache)

	u, err := s.Get(context.Background(), &user.UserView{ID: 1})

	assert.Error(t, err)
	assert.Equal(t, &user.NilUserView, u)
}

func TestUserService_Create_Successful(t *testing.T) {
	t.Parallel()

	repo := repoMock{
		user: &user.User{
			ID:   1,
			Name: "nelson",
		},
	}

	cache := &mockCache[user.UserView]{}
	s := NewUserService(repo, cache)

	u, err := s.Create(context.Background(), &user.UserView{ID: 1})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserService_Create_RepoError(t *testing.T) {
	t.Parallel()

	repo := repoMock{
		err: errors.New("create error"),
	}

	cache := &mockCache[user.UserView]{}
	s := NewUserService(repo, cache)

	u, err := s.Create(context.Background(), &user.UserView{ID: 1})

	assert.Error(t, err)
	assert.Equal(t, &user.NilUserView, u)
}

// mockCache for testing purposes
type mockCache[T interfaces.Cacheable] struct {
	user   T
	err    error
	exists bool
	length int
}

func (m *mockCache[T]) Get(_ int64) (T, bool) {
	return m.user, m.exists
}

func (m *mockCache[T]) Set(_ int64, _ *T) {
	return
}

func (m *mockCache[T]) Len() int {
	return m.length
}

// MOCKS
type repoMock struct {
	user *user.User
	err  error
}

func (c repoMock) Create(_ context.Context, _ *user.User) (*user.User, error) {
	return c.user, c.err
}

func (c repoMock) Get(_ context.Context, _ *user.User) (*user.User, error) {
	return c.user, c.err
}
