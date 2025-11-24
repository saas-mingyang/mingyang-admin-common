package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/entctx/tenantctx"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/mixins"
	ent2 "mingyang-admin-simple-admin-file/ent"
	"mingyang-admin-simple-admin-file/ent/hook"
	"mingyang-admin-simple-admin-file/ent/privacy"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("File's name | 文件名称"),
		field.Uint8("file_type").
			Comment("File's type | 文件类型"),
		field.Uint64("size").
			Comment("File's size | 文件大小"),
		field.String("path").
			Comment("File's path | 文件路径"),
		field.String("user_id").
			Comment("User's UUID | 用户的 UUID"),
		field.String("md5").
			Comment("The md5 of the file | 文件的 md5"),
	}
}

func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.StatusMixin{},
		mixins.TenantMixin{},
	}
}

func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("file_type"),
		index.Fields("path"),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", FileTag.Type).Ref("files"),
	}
}

// Hooks of the File.
func (File) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.FileFunc(func(ctx context.Context, m *ent2.FileMutation) (ent.Value, error) {
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

func (File) Policy() ent.Policy {
	return privacy.FilterTenantRule()
}

func (File) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("File Table | 文件表"),
		entsql.Annotation{Table: "fms_files"}, // fms means file management service
	}
}
