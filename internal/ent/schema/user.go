package schema

import (
	entpkg "source-base/pkg/database/ent"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("email"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entpkg.UUIDMixin{},
		entpkg.PublicIDMixin{},
		entpkg.TimeMixin{},
	}
}
