package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePetRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type PetResponse struct {
	PublicID  string      `json:"public_id"`
	Name      string      `json:"name"`
	UserID    uuid.UUID   `json:"user_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Owner     *UserSimple `json:"owner,omitempty"`
}

type PetSimple struct {
	ID       uuid.UUID `json:"id"`
	PublicID string    `json:"public_id"`
	UserID   uuid.UUID `json:"user_id"`
}
