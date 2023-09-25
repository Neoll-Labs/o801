package service

import (
	"context"
	"github.com/nelsonstr/o801/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockRepository mocks to run unit test.
type MockRepository struct {
	user  *User
	error error
}

func (f *MockRepository) Create(_ context.Context, _ *User) (*User, error) {
	return f.user, f.error
}

func (f *MockRepository) Get(_ context.Context, _ *User) (*User, error) {
	return f.user, f.error
}

func TestNewUserService(t *testing.T) {
	s := NewUserService(repoMock{})
	assert.NotNil(t, s.userCache)
	assert.NotNil(t, s.repository)
	assert.NotNil(t, s.errorHandler)
}

// MOCKS
type repoMock struct {
	user *model.User
	err  error
}

func (c repoMock) Create(ctx context.Context, t *model.User) (*model.User, error) {
	return c.user, c.err
}

func (c repoMock) Fetch(ctx context.Context, t *model.User) (*model.User, error) {
	return c.user, c.err
}

func Test_userService_GetFromCache(t *testing.T) {
	s := userService{
		userCache: map[int64]User{
			1: {
				ID:   1,
				Name: "nelson",
			},
		},
	}
	got, err := s.Get(context.Background(), &User{ID: 1})

	assert.NoError(t, err)
	assert.Equal(t, "nelson", got.Name)
}
