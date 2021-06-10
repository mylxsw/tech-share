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
	localMigrate "github.com/mylxsw/tech-share/migrate"
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
		Name:   "auth_provider",
		Usage:  "鉴权服务提供者：默认支持 database/ldap",
		EnvVar: "TECH_SHARE_AUTH_PROVIDER",
		Value:  "database",
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
		Name:   "ldap_user_filter",
		Usage:  "LDAP user filter",
		EnvVar: "TECH_SHARE_LDAP_USER_FILTER",
		Value:  "CN=all-staff,CN=Users,DC=example,DC=com",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "session_key",
		Usage:  "Session key",
		EnvVar: "TECH_SHARE_SESSION_KEY",
		Value:  "49a95f4cdaac9dedbc3298c5f5a7aa83",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "weak_password_mode",
		Usage:  "启用弱密码模式，启用该模式后，只使用 LDAP 的账号体系，密码使用本地密码",
		EnvVar: "TECH_SHARE_WEAK_PASSWORD_MODE",
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
			Version:          Version,
			GitCommit:        GitCommit,
			Listen:           c.String("listen"),
			Debug:            c.Bool("debug"),
			LogPath:          c.String("log_path"),
			DBConnStr:        c.String("db_conn_str"),
			StoragePath:      c.String("storage_path"),
			SessionKey:       c.String("session_key"),
			WeakPasswordMode: c.Bool("weak_password_mode"),
			AuthProvider:     c.String("auth_provider"),
			LDAP: config.LDAP{
				URL:         c.String("ldap_url"),
				BaseDN:      c.String("ldap_base_dn"),
				Username:    c.String("ldap_username"),
				Password:    c.String("ldap_password"),
				DisplayName: "displayName",
				UID:         "sAMAccountName",
				UserFilter:  c.String("ldap_user_filter"),
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
		localMigrate.Migrations(m)

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
