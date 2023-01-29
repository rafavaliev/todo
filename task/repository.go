package task

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context, options QueryOptions) ([]*Task, error)
	CountAll(ctx context.Context, options QueryOptions) (int64, error)
	FindByIDs(ctx context.Context, options QueryOptions) ([]*Task, error)
	FindByID(ctx context.Context, userID uint, id string) (*Task, error)
	Create(ctx context.Context, userId uint, task *Task) (*Task, error)
	Update(ctx context.Context, userId uint, task *UpdateTask) error
	Delete(ctx context.Context, userId uint, id string) error
}

type MockRepository struct {
	FindAllFn   func(ctx context.Context, options QueryOptions) ([]*Task, error)
	CountAllFn  func(ctx context.Context, options QueryOptions) (int64, error)
	FindByIDsFn func(ctx context.Context, options QueryOptions) ([]*Task, error)
	FindByIDFn  func(ctx context.Context, userID uint, id string) (*Task, error)
	CreateFn    func(ctx context.Context, userId uint, task *Task) (*Task, error)
	UpdateFn    func(ctx context.Context, userId uint, task *UpdateTask) error
	DeleteFn    func(ctx context.Context, userId uint, id string) error
}

func (m MockRepository) FindAll(ctx context.Context, options QueryOptions) ([]*Task, error) {
	return m.FindAllFn(ctx, options)
}

func (m MockRepository) CountAll(ctx context.Context, options QueryOptions) (int64, error) {
	return m.CountAllFn(ctx, options)
}

func (m MockRepository) FindByIDs(ctx context.Context, options QueryOptions) ([]*Task, error) {
	return m.FindByIDsFn(ctx, options)
}

func (m MockRepository) FindByID(ctx context.Context, userID uint, id string) (*Task, error) {
	return m.FindByIDFn(ctx, userID, id)
}

func (m MockRepository) Create(ctx context.Context, userId uint, task *Task) (*Task, error) {
	return m.CreateFn(ctx, userId, task)
}

func (m MockRepository) Update(ctx context.Context, userId uint, task *UpdateTask) error {
	return m.UpdateFn(ctx, userId, task)
}

func (m MockRepository) Delete(ctx context.Context, userId uint, id string) error {
	return m.DeleteFn(ctx, userId, id)
}
