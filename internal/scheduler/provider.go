package scheduler

import (
	"context"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/scheduler"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		scheduler.Provider(s.jobCreator),
	}
}

func (s Provider) Register(cc infra.Binder) {}
func (s Provider) Boot(cc infra.Resolver)   {}

func (s Provider) jobCreator(cc infra.Resolver, creator scheduler.JobCreator) {
	creator.MustAddAndRunOnServerReady("ldap-user-sync", "@every 1h", syncLdapUsers)
}

func syncLdapUsers(conf *config.Config, userSrv service.UserService) error {
	log.Debugf("start sync ldap users")

	l, err := ldap.DialURL(conf.LDAP.URL)
	if err != nil {
		return fmt.Errorf("无法连接 LDAP 服务器: %w", err)
	}

	defer l.Close()

	if err := l.Bind(conf.LDAP.Username, conf.LDAP.Password); err != nil {
		return fmt.Errorf("LDAP 服务器鉴权失败: %w", err)
	}

	searchReq := ldap.NewSearchRequest(
		conf.LDAP.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(memberOf=%s))", ldap.EscapeFilter(conf.LDAP.UserFilter)),
		[]string{"objectguid", conf.LDAP.UID, conf.LDAP.DisplayName, "userAccountControl"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return fmt.Errorf("LDAP 用户查询失败: %w", err)
	}

	for _, ent := range sr.Entries {
		userStatus := service.UserStatusEnabled
		if ent.GetAttributeValue("userAccountControl") == "514" {
			userStatus = service.UserStatusDisabled
		}

		if _, err := userSrv.LoadUser(context.TODO(), uuid.Must(uuid.FromBytes(ent.GetRawAttributeValue("objectGUID"))).String(), service.UserInfo{
			Name:    ent.GetAttributeValue(conf.LDAP.DisplayName),
			Account: ent.GetAttributeValue(conf.LDAP.UID),
			Status:  userStatus,
		}); err != nil {
			log.With(ent).Errorf("load user failed: %v", err)
		}
	}

	return nil
}
