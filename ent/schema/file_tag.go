package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/suyuan32/simple-admin-common/orm/ent/mixins"
	"github.com/suyuan32/simple-admin-common/orm/ent/tenantctx"
	ent2 "github.com/suyuan32/simple-admin-file/ent"
	"github.com/suyuan32/simple-admin-file/ent/hook"
)

// FileTag holds the schema definition for the FileTag entity.
type FileTag struct {
	ent.Schema
}

// Fields of the FileTag.
func (FileTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("FileTag's name | 标签名称").
			Annotations(entsql.WithComments(true)),
		field.String("remark").Comment("The remark of tag | 标签的备注").
			Optional().
			Annotations(entsql.WithComments(true)),
	}
}

func (FileTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.TenantMixin{},
	}
}

// Edges of the FileTag.
func (FileTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
	}
}

func (FileTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}

// Hooks of the FileTag.
func (FileTag) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.FileTagFunc(func(ctx context.Context, m *ent2.FileTagMutation) (ent.Value, error) {
					if !tenantctx.GetTenantAdminCtx(ctx) {
						m.SetTenantID(tenantctx.GetTenantIDFromCtx(ctx))
					}
					return next.Mutate(ctx, m)
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),
	}
}

func (FileTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fms_file_tags"}, // fms means FileTag management service
	}
}
