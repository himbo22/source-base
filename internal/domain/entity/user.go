package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	PublicID  string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time

	Pets []Pet
}
