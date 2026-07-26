package ports

import (
	"context"
	"github.com/himbo22/source-base/internal/domain/dto"
)

type PetController interface {
	Create(ctx context.Context, req *dto.CreatePetRequest) (dto.PetResponse, error)
}
