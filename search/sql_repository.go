package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"todo/user"
)

// We have to import this structure because I used gorm with auto migration.
type SQLUserIndex struct {
	UserID    uint `gorm:"primaryKey"`
	Index     string
	Analyzer  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(gorm *gorm.DB) UserIndexRepository {
	return &SQLRepository{db: gorm}
}

func (S *SQLRepository) Find(ctx context.Context, userID uint) (*UserIndex, error) {

	var sqlUserIndex SQLUserIndex
	err := S.db.WithContext(ctx).First(&sqlUserIndex, "user_id = ?", userID).Error
	switch {
	case err == nil:
		break
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("could not find user index: %w", err)
	}

	res := &UserIndex{UserID: userID}

	var idx Index
	err = json.Unmarshal([]byte(sqlUserIndex.Index), &idx)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal user index: %w", err)
	}
	res.Index = idx

	var analyzer Analyzer
	err = json.Unmarshal([]byte(sqlUserIndex.Analyzer), &analyzer)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal user index analyzer: %w", err)
	}
	res.Analyzer = &analyzer

	return res, nil

}

func (s *SQLRepository) Update(ctx context.Context, idx *UserIndex) error {
	usr := ctx.Value(user.UserContextKey).(user.User)

	jsonIdx, err := json.Marshal(idx.Index)
	if err != nil {
		return fmt.Errorf("could not marshal user index: %w", err)
	}

	err = s.db.WithContext(ctx).Model(&SQLUserIndex{}).Where("user_id = ?", usr.ID).Update("index", jsonIdx).Error
	if err != nil {
		return fmt.Errorf("could not update user index: %w", err)
	}

	return nil
}

func (s *SQLRepository) Create(ctx context.Context, idx *UserIndex) error {
	usr := ctx.Value(user.UserContextKey).(user.User)

	sqlIdx := &SQLUserIndex{
		UserID: usr.ID,
	}

	jsonIdx, err := json.Marshal(idx.Index)
	if err != nil {
		return fmt.Errorf("could not marshal user index: %w", err)
	}
	sqlIdx.Index = string(jsonIdx)

	analyzer, err := json.Marshal(idx.Analyzer)
	if err != nil {
		return fmt.Errorf("could not marshal user index: %w", err)
	}
	sqlIdx.Analyzer = string(analyzer)

	err = s.db.WithContext(ctx).Create(sqlIdx).Error
	if err != nil {
		return fmt.Errorf("could not update user index: %w", err)
	}

	return nil
}
