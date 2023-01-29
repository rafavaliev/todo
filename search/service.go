package search

import (
	"context"
	"fmt"
	"todo/user"
)

var (
	ErrNotFound = fmt.Errorf("user index not found")
)

// Service is a service for searching documents.
type Service struct {
	Repo UserIndexRepository
}

func NewService(repo UserIndexRepository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) Search(ctx context.Context, query string) ([]string, error) {
	usr := ctx.Value(user.UserContextKey).(user.User)

	userIndex, err := s.Repo.Find(ctx, usr.ID)
	if err == ErrNotFound {
		userIndex = &UserIndex{UserID: usr.ID, Index: Index{}}
		err = s.Repo.Create(ctx, userIndex)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user index: %w", err)
	}
	return userIndex.Search(query), nil
}

func (s *Service) Insert(ctx context.Context, document Document) error {
	usr := ctx.Value(user.UserContextKey).(user.User)

	userIndex, err := s.Repo.Find(ctx, usr.ID)
	if err == ErrNotFound {
		userIndex = &UserIndex{UserID: usr.ID, Index: Index{}}
		err = s.Repo.Create(ctx, userIndex)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return fmt.Errorf("failed to find user index: %w", err)
	}
	userIndex.Insert(document)
	err = s.Repo.Update(ctx, userIndex)
	if err != nil {
		return fmt.Errorf("failed to update user index: %w", err)
	}
	return nil

}

func (s *Service) Delete(ctx context.Context, document Document) error {
	usr := ctx.Value(user.UserContextKey).(user.User)

	userIndex, err := s.Repo.Find(ctx, usr.ID)
	if err != nil {
		return fmt.Errorf("failed to find user index: %w", err)
	}
	userIndex.Delete(document)
	err = s.Repo.Update(ctx, userIndex)
	if err != nil {
		return fmt.Errorf("failed to update user index: %w", err)
	}
	return nil
}
