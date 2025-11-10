package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/saas-mingyang/mingyang-admin-common/utils/sonyflake"
	"time"
)

// IdSonyFlakeMixin id生成器
type IdSonyFlakeMixin struct {
	mixin.Schema
}

func (IdSonyFlakeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Default(uint64(sonyflake.NextID())),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Comment("Create Time | 创建日期"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Update Time | 修改日期"),
	}
}
