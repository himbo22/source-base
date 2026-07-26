package ent

import (
	"github.com/himbo22/source-base/pkg/unique"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type UUIDMixin struct {
	mixin.Schema
}

func (u UUIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(NewUUID).
			Immutable(),
	}
}

func NewUUID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

type PublicIDMixin struct {
	mixin.Schema
	Prefix  string
	RandLen int
}

func (u PublicIDMixin) Fields() []ent.Field {
	n := u.RandLen
	if n <= 0 {
		n = 6
	}
	return []ent.Field{
		field.String("public_id").
			Unique().
			Immutable().
			DefaultFunc(func() string {
				return unique.PublicID(u.Prefix, time.Now(), n)
			}),
	}
}
