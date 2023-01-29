package user

import (
	"fmt"
	"time"
)

const UserContextKey = "user"

var (
	ErrNotFound = fmt.Errorf("user not found")
)

type User struct {
	ID             uint   `gorm:"primarykey"`
	Username       string `gorm:"uniqueIndex"`
	HashedPassword []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
