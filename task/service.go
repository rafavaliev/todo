package task

import (
	"context"
	"fmt"
	"todo/search"
	"todo/user"
)

type Service struct {
	Repo          Repository
	SearchService *search.Service
}

type QueryOptions struct {
	IDs             []string
	UserID          uint
	Limit           int
	Offset          int
	IncludeStatuses []Status
}

func NewService(repo Repository, searchService *search.Service) *Service {
	return &Service{
		Repo:          repo,
		SearchService: searchService,
	}
}

func (s *Service) Search(ctx context.Context, query string, opts QueryOptions) ([]*Task, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)
	opts.UserID = usr.ID

	documentIDs, err := s.SearchService.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	opts.IDs = documentIDs
	return s.Repo.FindByIDs(ctx, opts)
}

func (s *Service) FindAll(ctx context.Context, opts QueryOptions) ([]*Task, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)
	opts.UserID = usr.ID

	return s.Repo.FindAll(ctx, opts)
}

func (s *Service) CountAll(ctx context.Context, opts QueryOptions) (int64, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)
	opts.UserID = usr.ID

	return s.Repo.CountAll(ctx, opts)
}

func (s *Service) FindByID(ctx context.Context, id string) (*Task, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)

	return s.Repo.FindByID(ctx, usr.ID, id)
}

func (s *Service) Create(ctx context.Context, task *Task) (*Task, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)
	task.UserID = usr.ID
	if err := task.Validate(); err != nil {

	}
	t, err := s.Repo.Create(ctx, usr.ID, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	_ = s.SearchService.Insert(ctx, search.Document{
		ID:      task.ID,
		Content: fmt.Sprintf("%s %s", task.Title, task.Description),
	})

	return t, nil
}

func (s *Service) Update(ctx context.Context, task *UpdateTask) (*Task, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)
	oldTask := ctx.Value(TaskContextKey).(*Task)

	// Delete old task from search index
	_ = s.SearchService.Delete(ctx, search.Document{
		ID:      oldTask.ID,
		Content: fmt.Sprintf("%s %s", oldTask.Title, oldTask.Description),
	})
	// Update task in database
	err := s.Repo.Update(ctx, usr.ID, task)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	newTask, err := s.Repo.FindByID(ctx, usr.ID, oldTask.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Update search index
	_ = s.SearchService.Insert(ctx, search.Document{
		ID:      oldTask.ID,
		Content: fmt.Sprintf("%s %s", newTask.Title, newTask.Description),
	})

	return newTask, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	usr := ctx.Value(user.UserContextKey).(user.User)
	task := ctx.Value(TaskContextKey).(*Task)

	// Delete the task from search index
	_ = s.SearchService.Delete(ctx, search.Document{
		ID:      id,
		Content: fmt.Sprintf("%s %s", task.Title, task.Description),
	})
	// Delete the task from database
	err := s.Repo.Delete(ctx, usr.ID, id)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
