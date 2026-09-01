package ent

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// SoftDeleteMixin adds a nullable deleted_at timestamp field for soft delete.
type SoftDeleteMixin struct {
	mixin.Schema
}

func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time(SoftDeleteAtColumnName).
			Optional().
			Nillable(),
	}
}
