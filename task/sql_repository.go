package task

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(gorm *gorm.DB) *SQLRepository {
	return &SQLRepository{db: gorm}
}

func (s *SQLRepository) FindAll(ctx context.Context, options QueryOptions) ([]*Task, error) {
	var tasks []*Task
	tx := s.db.WithContext(ctx).
		Limit(options.Limit).
		Offset(options.Offset).
		Where("user_id = ? ", options.UserID).
		Find(&tasks)

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failed to find tasks by ids: %w", err)
	}
	return tasks, nil
}

func (s *SQLRepository) CountAll(ctx context.Context, options QueryOptions) (int64, error) {
	count := int64(0)
	tx := s.db.WithContext(ctx).Model(&Task{}).
		Where("user_id = ? ", options.UserID).
		Count(&count)
	if err := tx.Error; err != nil {
		return 0, fmt.Errorf("failed to find tasks by ids: %w", err)
	}
	return count, nil
}

func (s *SQLRepository) FindByIDs(ctx context.Context, options QueryOptions) ([]*Task, error) {
	var tasks []*Task
	if err := s.db.Find(&tasks, options.IDs).Error; err != nil {

		return nil, fmt.Errorf("failed to find tasks by ids: %w", err)

	}
	return tasks, nil
}

func (s *SQLRepository) FindByID(ctx context.Context, userID uint, id string) (*Task, error) {
	var task Task
	tx := s.db.Where("user_id = ? AND id = ?", userID, id).First(&task)
	if err := tx.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("failed to find task by id: %w", err)
		}

	}
	return &task, nil
}

func (s *SQLRepository) Create(ctx context.Context, userID uint, task *Task) (*Task, error) {
	err := s.db.Create(task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func (s *SQLRepository) Update(ctx context.Context, userID uint, task *UpdateTask) error {
	tx := s.db.WithContext(ctx).Model(&Task{}).Where("user_id = ? AND id = ?", userID, task.ID)
	if task.Title != nil {
		tx = tx.Update("title", task.Title)
	}
	if task.Description != nil {
		tx = tx.Update("description", task.Description)
	}
	if task.Status != nil {
		tx = tx.Update("status", task.Status)
	}
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (s *SQLRepository) Delete(ctx context.Context, userID uint, id string) error {
	tx := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).Delete(&Task{})
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
