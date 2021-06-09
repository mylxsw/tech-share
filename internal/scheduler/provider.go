package scheduler

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/scheduler"
	"github.com/mylxsw/tech-share/config"
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
	creator.Add("ldap-user-sync", "@every 1h", syncLdapUsers)
	cc.Call(syncLdapUsers)
}

func syncLdapUsers(conf *config.Config) {
	l, err := ldap.DialURL(conf.LDAP.URL)
	if err != nil {
		panic(fmt.Errorf("无法连接 LDAP 服务器: %w", err))
	}

	defer l.Close()

	if err := l.Bind(conf.LDAP.Username, conf.LDAP.Password); err != nil {
		panic(fmt.Errorf("LDAP 服务器鉴权失败: %w", err))
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
		panic(fmt.Errorf("LDAP 用户查询失败: %w", err))
	}

	for _, ent := range sr.Entries {
		logFields := log.Fields{
			"dn":                 ent.DN,
			"account":            ent.GetAttributeValue(conf.LDAP.UID),
			"displayName":        ent.GetAttributeValue(conf.LDAP.DisplayName),
			"uuid":               uuid.Must(uuid.FromBytes(ent.GetRawAttributeValue("objectGUID"))).String(),
			"userAccountControl": ent.GetAttributeValue("userAccountControl"),
		}
		if ent.GetAttributeValue("userAccountControl") == "514" {
			log.WithFields(logFields).Error("ldap entries")
		} else {
			log.WithFields(logFields).Debug("ldap entries")
		}
	}
}
