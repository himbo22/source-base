package ent

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

func utcNow() time.Time { return time.Now().UTC() }

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time(CreatedAtColumnName).
			Default(utcNow).
			Immutable().
			Annotations(entsql.Default("CURRENT_TIMESTAMP")),
		field.Time(UpdatedAtColumnName).
			Default(utcNow).
			UpdateDefault(utcNow).
			Annotations(entsql.Default("CURRENT_TIMESTAMP")),
	}
}
