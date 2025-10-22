package repository

import (
	"errors"

	"github.com/Cere6rum/MicroBlog2/internal/models"
	"github.com/Cere6rum/MicroBlog2/internal/syncutils"
)

// UserRepository defines abstraction for user storage.
type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	Exists(username string) bool
}

// InMemoryUserRepo is an adapter over syncutils.SafeUserStorage.
type InMemoryUserRepo struct {
	storage *syncutils.SafeUserStorage
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{storage: syncutils.NewSafeUserStorage()}
}

func (r *InMemoryUserRepo) Create(user *models.User) error {
	if r.Exists(user.Username) {
		return errors.New("user already exists")
	}
	r.storage.Set(user.Username, user)
	return nil
}

func (r *InMemoryUserRepo) GetByUsername(username string) (*models.User, error) {
	v, ok := r.storage.Get(username)
	if !ok {
		return nil, errors.New("user not found")
	}
	u, ok := v.(*models.User)
	if !ok {
		return nil, errors.New("stored value has unexpected type")
	}
	return u, nil
}

func (r *InMemoryUserRepo) Exists(username string) bool {
	return r.storage.Exists(username)
}
