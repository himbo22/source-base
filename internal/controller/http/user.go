package http

import (
	"context"

	"github.com/himbo22/source-base/internal/domain/dto"
	"github.com/himbo22/source-base/internal/ports"
)

type userController struct {
	userService ports.UserService
}

func NewUserController(userService ports.UserService) ports.UserController {
	return &userController{
		userService: userService,
	}
}

func (u *userController) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	return u.userService.Create(ctx, req)
}

func (u *userController) GetByPublicID(ctx context.Context, req *dto.GetUserByPublicIDRequest) (*dto.UserResponse, error) {
	return u.userService.GetByPublicID(ctx, req)
}
