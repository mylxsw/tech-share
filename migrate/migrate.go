package migrate

import "github.com/mylxsw/eloquent/migrate"

func Migrations(m *migrate.Manager) {
	m.Schema("202104222300").Create("share", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("subject", 255).Comment("分享主题")
		builder.Text("description").Nullable(true).Comment("内容简介")
		builder.String("subject_type", 50).Comment("分享类型")
		builder.TinyInteger("status", false, true).Default(migrate.RawExpr("0")).Comment("状态：0-投票中 1-已排期 2-已完结")
		builder.String("share_user", 100).Nullable(true).Comment("分享人")
		builder.Integer("create_user_id", false, true).Nullable(true).Comment("创建人id")
		builder.Text("note").Nullable(true).Comment("备注")
		builder.Integer("like_count", false, true).Default(migrate.RawExpr("0")).Comment("感兴趣人数")
		builder.Integer("join_count", false, true).Default(migrate.RawExpr("0")).Comment("参加人数")
		builder.String("attachments", 255).Default(migrate.RawExpr("")).Nullable(true).Comment("相关附件id列表，多个使用英文逗号分隔")
		builder.Timestamps(0)
		builder.SoftDeletes("deleted_at", 0)
	})
	m.Schema("202104222309").Create("share_user_rel", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.Integer("user_id", false, true).Comment("用户id")
		builder.Integer("share_id", false, true).Comment("分享id")
		builder.TinyInteger("rel_type", false, true).Comment("关联类型：1-感兴趣 2-参加")
		builder.Timestamp("created_at", 0).Nullable(true)
	})
	m.Schema("202104222322").Create("user", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("uuid", 255).Comment("用户 uuid")
		builder.String("name", 100).Comment("用户名")
		builder.Timestamps(0)
	})
	m.Schema("202104222323").Create("share_plan", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.Integer("share_id", false, true).Comment("分享id")
		builder.Timestamp("share_at", 0).Nullable(true).Comment("分享时间")
		builder.Integer("plan_duration", false, true).Nullable(true).Comment("预计分享时间")
		builder.Integer("real_duration", false, true).Nullable(true).Comment("实际分享时间")
		builder.Text("note").Nullable(true).Comment("备注")
		builder.Timestamps(0)
	})
	m.Schema("202104222325").Create("attachment", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("name", 255).Nullable(true).Comment("附件名称")
		builder.String("atta_type", 20).Nullable(true).Comment("附件类型")
		builder.String("atta_path", 255).Nullable(true).Comment("附件地址")
		builder.Timestamp("created_at", 0).Nullable(true)
	})
	m.Schema("202104222327").Create("journal", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("action", 30).Nullable(true).Comment("动作")
		builder.Text("params").Nullable(true).Comment("动作参数")
		builder.Integer("trigger_user_id", false, true).Comment("用户id")
		builder.Timestamp("created_at", 0).Nullable(true)
	})
	m.Schema("202105071404").Table("share_plan", func(builder *migrate.Builder) {
		builder.String("share_room", 255).Nullable(true).Comment("分享会议室")
	})
	m.Schema("202106071717").Table("share", func(builder *migrate.Builder) {
		builder.Integer("share_user_id", false, true).Nullable(true).Comment("分享用户 id")
		builder.Timestamp("share_at", 0).Nullable(true).Comment("分享时间，冗余字段，便于查询")
	})
	m.Schema("202106100943").Table("user", func(builder *migrate.Builder) {
		builder.String("account", 100).Comment("账号名")
		builder.TinyInteger("status", false, true).Default(migrate.RawExpr("1")).Comment("状态：0-禁用 1-启用")
	})
	m.Schema("202106102309").Table("user", func(builder *migrate.Builder) {
		builder.String("password", 256).Nullable(true).Comment("密码")
	})
}
