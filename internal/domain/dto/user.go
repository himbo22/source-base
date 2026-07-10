package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	PublicID  string        `json:"public_id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
	Pets      []PetResponse `json:"pets,omitempty"`
}

type GetUserByPublicIDRequest struct {
	PublicID string `param:"public_id" validate:"required"`
}

type UserSimple struct {
	ID       uuid.UUID `json:"id"`
	PublicID string    `json:"public_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
}
