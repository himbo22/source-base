package entity

import (
	"time"

	"github.com/google/uuid"
)

type Pet struct {
	ID        uuid.UUID
	PublicID  string
	Name      string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Owner *User
}
