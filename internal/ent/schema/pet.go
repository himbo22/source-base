package schema

import (
	entpkg "source-base/pkg/database/ent"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.UUID("user_id", uuid.UUID{}),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Field("user_id").
			Required().
			Unique(),
	}
}

func (Pet) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entpkg.UUIDMixin{},
		entpkg.PublicIDMixin{},
		entpkg.TimeMixin{},
	}
}

//func (Pet) Annotations() []schema.Annotation {
//	return []schema.Annotation{
//		entsql.Annotation{},
//	}
//}
