package service

import (
	"context"
	"net/http"
	"source-base/internal/domain/dto"
	"source-base/internal/domain/entity"
	"source-base/internal/ports"
	"source-base/pkg/common/apperror"
	"source-base/pkg/common/http/response"

	"go.uber.org/zap"
)

type userService struct {
	userRepo ports.UserRepository
	logger   *zap.Logger
}

func NewUserService(
	userRepo ports.UserRepository,
	logger *zap.Logger,
) ports.UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *userService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := entity.User{
		Name:  req.Name,
		Email: req.Email,
	}

	err := u.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	r := &dto.UserResponse{
		PublicID:  user.PublicID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return r, nil
}

func (u *userService) GetByPublicID(ctx context.Context, req *dto.GetUserByPublicIDRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetByPublicID(ctx, req.PublicID)
	if err == nil && user == nil {
		return nil, apperror.New(response.CodeNotFound, "user not found", http.StatusNotFound)
	}
	if err != nil {
		return nil, err
	}

	r := &dto.UserResponse{
		PublicID:  user.PublicID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if len(user.Pets) > 0 {
		r.Pets = make([]dto.PetResponse, len(user.Pets))
		for i, p := range user.Pets {
			r.Pets[i] = dto.PetResponse{
				PublicID:  p.PublicID,
				Name:      p.Name,
				UserID:    p.UserID,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
			}
		}
	}

	return r, nil
}
