package task

import (
	"encoding/json"
	"errors"
	"time"
)

type Status string

// By default todo lists only have a boolean flag that represents task status(completed/active)
// Enum statuses allows us to expand product, for example as a Kanban board of tasks with several statuses
var (
	CreatedStatus  Status = "created"
	FinishedStatus Status = "finished"
	ArchivedStatus Status = "archived"
)

var (
	ErrEmptyTitle    = errors.New("title is empty")
	ErrInvalidStatus = errors.New("invalid status")
	ErrNotFound      = errors.New("task not found")
)

const TaskContextKey string = "task_ctx"

type Task struct {
	ID          string    `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return ErrEmptyTitle
	}
	switch t.Status {
	case CreatedStatus, FinishedStatus, ArchivedStatus:
		break
	default:
		return ErrInvalidStatus
	}
	return nil
}

func (t Task) String() string {
	jsonTask, _ := json.Marshal(t)
	return string(jsonTask)
}

type UpdateTask struct {
	ID          string
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *Status `json:"status"`
}

func (t *UpdateTask) Validate() error {
	if t.Title != nil && *t.Title == "" {
		return ErrEmptyTitle
	}
	if t.Status == nil {
		return ErrInvalidStatus
	}
	if t.Status != nil {
		switch *t.Status {
		case CreatedStatus, FinishedStatus, ArchivedStatus:
			break
		default:
			return ErrInvalidStatus
		}
	}
	return nil
}
