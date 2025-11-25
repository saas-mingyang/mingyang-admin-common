package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/saas-mingyang/mingyang-admin-common/orm/ent/mixins"
)

type Apk struct {
	ent.Schema
}

// Fields of the APK.
func (Apk) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("APK名称"),
		field.String("version").
			Comment("版本号"),
		field.String("version_code").Unique().
			Comment("版本代码(内部版本号)"),
		field.Uint64("file_size").
			Default(0).
			Comment("文件id"),
		field.Uint64("file_id").
			Comment("下载地址"),
		field.String("file_path").
			Optional().
			Comment("文件存储路径"),
		field.String("md5").
			Optional().
			Comment("文件MD5值"),
		field.String("sha1").
			Optional().
			Comment("文件SHA1值"),
		field.String("sha256").
			Optional().
			Comment("文件SHA256值"),
		field.String("package_name").
			Optional().
			Comment("应用包名"),
		field.Text("description").
			Optional().
			Comment("版本描述"),
		field.String("update_log").
			Optional().
			Comment("更新日志"),
		field.Bool("is_force_update").
			Default(false).
			Comment("是否强制更新"),
		field.Int64("download_count").
			Default(0).
			Comment("下载次数"),
	}
}

func (Apk) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.TenantMixin{},
	}
}

func (Apk) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("APK版本记录"),
		entsql.Annotation{Table: "apk_file"},
	}
}

func (Apk) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("version_code").Unique(),
	}
}
