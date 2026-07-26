package ent

import (
	"github.com/himbo22/source-base/internal/domain/entity"
	"github.com/himbo22/source-base/internal/ent/generate"
)

func ToDomainPet(p *generate.Pet) *entity.Pet {
	pet := &entity.Pet{
		ID:        p.ID,
		PublicID:  p.PublicID,
		Name:      p.Name,
		UserID:    p.UserID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if p.Edges.Owner != nil {
		pet.Owner = ToDomainUser(p.Edges.Owner)
	}
	return pet
}

func ToDomainUser(u *generate.User) *entity.User {
	user := &entity.User{
		ID:        u.ID,
		PublicID:  u.PublicID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.Edges.Pets != nil {
		user.Pets = make([]entity.Pet, len(u.Edges.Pets))
		for i, p := range u.Edges.Pets {
			user.Pets[i] = *ToDomainPet(p)
		}
	}
	return user
}
