package repository

import (
	"errors"

	"github.com/Cere6rum/MicroBlog2/internal/models"
	"github.com/Cere6rum/MicroBlog2/internal/syncutils"
)

// PostRepository defines abstraction for post storage.
type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id int) (*models.Post, error)
	List() []*models.Post
	Update(post *models.Post) error
}

// InMemoryPostRepo is an adapter over syncutils.SafePostStorage.
type InMemoryPostRepo struct {
	storage *syncutils.SafePostStorage
}

func NewInMemoryPostRepo() *InMemoryPostRepo {
	return &InMemoryPostRepo{storage: syncutils.NewSafePostStorage()}
}

func (r *InMemoryPostRepo) Create(post *models.Post) error {
	r.storage.Add(post)
	return nil
}

func (r *InMemoryPostRepo) GetByID(id int) (*models.Post, error) {
	v, ok := r.storage.GetByIndex(id - 1)
	if !ok {
		return nil, errors.New("post not found")
	}
	p, ok := v.(*models.Post)
	if !ok {
		return nil, errors.New("stored value has unexpected type")
	}
	return p, nil
}

func (r *InMemoryPostRepo) List() []*models.Post {
	raw := r.storage.GetAll()
	out := make([]*models.Post, 0, len(raw))
	for _, v := range raw {
		if p, ok := v.(*models.Post); ok {
			out = append(out, p)
		}
	}
	return out
}

func (r *InMemoryPostRepo) Update(post *models.Post) error {
	// Naive implementation: replace by index if exists
	if post.ID <= 0 || post.ID > r.storage.Len() {
		return errors.New("post not found")
	}
	r.storage.SetByIndex(post.ID-1, post)
	return nil
}
