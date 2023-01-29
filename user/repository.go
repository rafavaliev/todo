package user

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}

type SQLRepository struct {
	db *gorm.DB
}

func (s SQLRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var usr User
	tx := s.db.WithContext(ctx).Where("username = ?", username).First(&usr)
	switch {
	case tx.Error == nil:
		return &usr, nil
	case errors.Is(tx.Error, gorm.ErrRecordNotFound):
		return nil, ErrNotFound
	default:
		return nil, tx.Error
	}

}

func (s SQLRepository) Create(ctx context.Context, user *User) (*User, error) {
	tx := s.db.WithContext(ctx).Create(user)
	switch {
	case tx.Error == nil:
		return user, nil
	default:
		return nil, tx.Error
	}
}

func NewSQLRepository(gorm *gorm.DB) Repository {
	return &SQLRepository{db: gorm}
}

type MockRepository struct {
	FindByUsernameFn func(ctx context.Context, username string) (*User, error)
	CreateFn         func(ctx context.Context, user *User) (*User, error)
}

func (m MockRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	return m.FindByUsernameFn(ctx, username)
}

func (m MockRepository) Create(ctx context.Context, user *User) (*User, error) {
	return m.CreateFn(ctx, user)
}
