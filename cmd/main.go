package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/asteria/writer"
	"github.com/mylxsw/container"
	"github.com/mylxsw/eloquent/migrate"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/mylxsw/tech-share/internal/service"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"

	"github.com/mylxsw/tech-share/api"
	"github.com/mylxsw/tech-share/config"
	localEvt "github.com/mylxsw/tech-share/internal/event"
	"github.com/mylxsw/tech-share/internal/scheduler"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	app := application.Create(fmt.Sprintf("%s %s", Version, GitCommit))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "listen",
		Usage:  "服务监听地址",
		Value:  "127.0.0.1:19921",
		EnvVar: "TECH_SHARE_LISTEN",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "debug",
		Usage:  "是否使用调试模式，调试模式下，静态资源使用本地文件",
		EnvVar: "TECH_SHARE_DEBUG",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "log_path",
		Usage: "日志文件输出目录（非文件名），默认为空，输出到标准输出",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "db_conn_str",
		Usage:  "数据库连接字符串",
		EnvVar: "TECH_SHARE_DB_CONN",
		Value:  "tech_share:tech-share123@tcp(127.0.0.1:3306)/tech_share?parseTime=true",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "storage_path",
		Usage:  "文件存储目录",
		EnvVar: "TECH_SHARE_STORAGE_PATH",
		Value:  "/home/mylxsw/Downloads",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "ldap_url",
		Usage:  "LDAP 服务器地址",
		EnvVar: "TECH_SHARE_LDAP_URL",
		Value:  "ldap://127.0.0.1:389",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "ldap_username",
		Usage:  "LDAP 账号",
		EnvVar: "TECH_SHARE_LDAP_USERNAME",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "ldap_password",
		Usage:  "LDAP 密码",
		EnvVar: "TECH_SHARE_LDAP_PASSWORD",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "ldap_base_dn",
		Usage:  "LDAP base dn",
		EnvVar: "TECH_SHARE_LDAP_BASE_DN",
		Value:  "dc=example,dc=com",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "session_key",
		Usage:  "Session key",
		EnvVar: "TECH_SHARE_SESSION_KEY",
		Value:  "49a95f4cdaac9dedbc3298c5f5a7aa83",
	}))

	app.BeforeServerStart(func(cc container.Container) error {
		stackWriter := writer.NewStackWriter()
		cc.MustResolve(func(ctx context.Context, c infra.FlagContext) {
			if !c.Bool("debug") {
				log.All().LogLevel(level.Info)
			}

			logPath := c.String("log_path")
			if logPath == "" {
				stackWriter.PushWithLevels(writer.NewStdoutWriter())
				return
			}

			log.All().LogFormatter(formatter.NewJSONWithTimeFormatter())
			stackWriter.PushWithLevels(writer.NewDefaultRotatingFileWriter(ctx, func(le level.Level, module string) string {
				return filepath.Join(logPath, fmt.Sprintf("%s-%s.log", time.Now().Format("20060102"), le.GetLevelName()))
			}))
		})

		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(app.Container()),
			level.Error,
			level.Emergency,
			level.Critical,
		)
		log.All().LogWriter(stackWriter)

		return nil
	})

	app.Singleton(func(c infra.FlagContext) *config.Config {
		return &config.Config{
			Version:     Version,
			GitCommit:   GitCommit,
			Listen:      c.String("listen"),
			Debug:       c.Bool("debug"),
			LogPath:     c.String("log_path"),
			DBConnStr:   c.String("db_conn_str"),
			StoragePath: c.String("storage_path"),
			SessionKey:  c.String("session_key"),
			LDAP: config.LDAP{
				URL:         c.String("ldap_url"),
				BaseDN:      c.String("ldap_base_dn"),
				Username:    c.String("ldap_username"),
				Password:    c.String("ldap_password"),
				DisplayName: "displayName",
				UID:         "sAMAccountName",
			},
		}
	})
	app.Singleton(func(conf *config.Config) (*sql.DB, error) {
		return sql.Open("mysql", conf.DBConnStr)
	})

	app.BeforeServerStop(func(cc infra.Resolver) error {
		return cc.Resolve(func(em event.Publisher) {
			em.Publish(localEvt.SystemUpDownEvent{
				Up:        false,
				CreatedAt: time.Now(),
			})
		})
	})

	app.Main(func(conf *config.Config, em event.Publisher, db *sql.DB) {
		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"config": conf,
			}).Debug("configuration")
		}

		em.Publish(localEvt.SystemUpDownEvent{
			Up:        true,
			CreatedAt: time.Now(),
		})

		m := migrate.NewManager(db).Init()
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
			builder.Integer("share_id", false, true).Default(migrate.RawExpr("0")).Comment("分享id")
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

		if err := m.Run(); err != nil {
			panic(err)
		}
	})

	app.Provider(api.Provider{}, localEvt.Provider{}, scheduler.Provider{})
	app.Provider(service.Provider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

type ErrorCollectorWriter struct {
	cc container.Container
}

func NewErrorCollectorWriter(cc container.Container) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{cc: cc}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	// TODO  Error collector implementation
	return nil
}

func (e *ErrorCollectorWriter) ReOpen() error {
	return nil
}

func (e *ErrorCollectorWriter) Close() error {
	return nil
}
