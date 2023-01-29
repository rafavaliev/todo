package search

import "context"

// UserIndexRepository is a repository for each user's search index.
type UserIndexRepository interface {
	Find(ctx context.Context, userID uint) (*UserIndex, error)
	Update(ctx context.Context, userIndex *UserIndex) error
	Create(ctx context.Context, userIndex *UserIndex) error
}

type MockUserIndexRepository struct {
	FindFn   func(ctx context.Context, userID uint) (*UserIndex, error)
	UpdateFn func(ctx context.Context, userIndex *UserIndex) error
	CreateFn func(ctx context.Context, userIndex *UserIndex) error
}

func (m MockUserIndexRepository) Find(ctx context.Context, userID uint) (*UserIndex, error) {
	return m.FindFn(ctx, userID)
}

func (m MockUserIndexRepository) Update(ctx context.Context, userIndex *UserIndex) error {
	return m.UpdateFn(ctx, userIndex)
}

func (m MockUserIndexRepository) Create(ctx context.Context, userIndex *UserIndex) error {
	return m.CreateFn(ctx, userIndex)
}
