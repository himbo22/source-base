package ports

import (
	"context"
	"source-base/internal/domain/dto"
)

type PetController interface {
	Create(ctx context.Context, req *dto.CreatePetRequest) (dto.PetResponse, error)
}
