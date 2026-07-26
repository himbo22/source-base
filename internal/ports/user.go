package ports

import (
	"context"
	"github.com/himbo22/source-base/internal/domain/dto"
	"github.com/himbo22/source-base/internal/domain/entity"
)

type UserController interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByPublicID(ctx context.Context, req *dto.GetUserByPublicIDRequest) (*dto.UserResponse, error)
}

type UserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByPublicID(ctx context.Context, req *dto.GetUserByPublicIDRequest) (*dto.UserResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPublicID(ctx context.Context, publicID string) (*entity.User, error)
}
