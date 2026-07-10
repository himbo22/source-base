package repository

import (
	"context"
	"source-base/internal/domain/entity"
	"source-base/internal/ent"
	"source-base/internal/ent/generate"
	"source-base/internal/ent/generate/user"
	"source-base/internal/ports"
)

type userRepository struct {
	client *generate.Client
}

func NewUserRepository(client *generate.Client) ports.UserRepository {
	return &userRepository{
		client: client,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	client := GetClient(ctx, u.client)

	model, err := client.User.Create().SetName(user.Name).SetEmail(user.Email).Save(ctx)
	if err != nil {
		return err
	}

	created := ent.ToDomainUser(model)
	*user = *created

	return err
}

func (u *userRepository) GetByPublicID(ctx context.Context, publicID string) (*entity.User, error) {
	client := GetClient(ctx, u.client)

	model, err := client.User.Query().Where(user.PublicIDEQ(publicID)).WithPets().Only(ctx)
	if err != nil {
		if generate.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return ent.ToDomainUser(model), nil
}
