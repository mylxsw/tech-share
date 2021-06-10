package ldap

import (
	"errors"
	"fmt"

	lp "github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/auth"
	"github.com/mylxsw/tech-share/internal/service"
)

type LdapAuth struct {
	conf *config.Config
}

func New(conf *config.Config) auth.Auth {
	return &LdapAuth{conf: conf}
}

func (provider *LdapAuth) Login(username, password string) (*auth.AuthedUser, error) {
	conf := provider.conf

	l, err := lp.DialURL(conf.LDAP.URL)
	if err != nil {
		return nil, fmt.Errorf("无法连接 LDAP 服务器: %w", err)
	}

	defer l.Close()

	if err := l.Bind(conf.LDAP.Username, conf.LDAP.Password); err != nil {
		return nil, fmt.Errorf("LDAP 服务器鉴权失败: %w", err)
	}

	searchReq := lp.NewSearchRequest(
		conf.LDAP.BaseDN,
		lp.ScopeWholeSubtree,
		lp.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", conf.LDAP.UID, lp.EscapeFilter(username)),
		[]string{"objectguid", conf.LDAP.UID, conf.LDAP.DisplayName, "userAccountControl"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP 用户查询失败: %w", err)
	}

	if len(sr.Entries) != 1 {
		return nil, service.NewValidateError(fmt.Errorf("用户不存在"))
	}

	// 514-禁用 512-启用
	if sr.Entries[0].GetAttributeValue("userAccountControl") == "514" {
		return nil, service.NewValidateError(errors.New("LDAP 用户账户已禁用"))
	}

	if conf.WeakPasswordMode {
		// TODO 弱密码模式
	} else {
		if err := l.Bind(sr.Entries[0].DN, password); err != nil {
			return nil, service.NewValidateError(fmt.Errorf("用户密码错误: %w", err))
		}
	}

	authedUser := buildAuthedUserFromLDAPEntry(conf.LDAP, sr.Entries[0])
	return &authedUser, nil
}

func buildAuthedUserFromLDAPEntry(conf config.LDAP, entry *lp.Entry) auth.AuthedUser {
	userStatus := service.UserStatusEnabled
	if entry.GetAttributeValue("userAccountControl") == "514" {
		userStatus = service.UserStatusDisabled
	}

	return auth.AuthedUser{
		UUID:    uuid.Must(uuid.FromBytes(entry.GetRawAttributeValue("objectGUID"))).String(),
		Name:    entry.GetAttributeValue(conf.DisplayName),
		Account: entry.GetAttributeValue(conf.UID),
		Status:  userStatus,
	}
}

func (provider *LdapAuth) Users() ([]auth.AuthedUser, error) {
	log.Debugf("start sync ldap users")

	l, err := lp.DialURL(provider.conf.LDAP.URL)
	if err != nil {
		return nil, fmt.Errorf("无法连接 LDAP 服务器: %w", err)
	}

	defer l.Close()

	if err := l.Bind(provider.conf.LDAP.Username, provider.conf.LDAP.Password); err != nil {
		return nil, fmt.Errorf("LDAP 服务器鉴权失败: %w", err)
	}

	searchReq := lp.NewSearchRequest(
		provider.conf.LDAP.BaseDN,
		lp.ScopeWholeSubtree,
		lp.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(memberOf=%s))", lp.EscapeFilter(provider.conf.LDAP.UserFilter)),
		[]string{"objectguid", provider.conf.LDAP.UID, provider.conf.LDAP.DisplayName, "userAccountControl"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP 用户查询失败: %w", err)
	}

	authedUsers := make([]auth.AuthedUser, 0)
	for _, ent := range sr.Entries {
		authedUsers = append(authedUsers, buildAuthedUserFromLDAPEntry(provider.conf.LDAP, ent))
	}

	return authedUsers, nil
}
