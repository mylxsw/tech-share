package database

import (
	"context"
	"errors"

	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/auth"
	"github.com/mylxsw/tech-share/internal/service"
	"github.com/mylxsw/tech-share/pkg/bcrypt"
)

type DatabaseAuth struct {
	srv  service.UserService
	conf *config.Config
}

func (au DatabaseAuth) CanRegister() bool {
	return true
}

func (au DatabaseAuth) Register(username, password string) (*auth.AuthedUser, error) {
	user, err := au.srv.Register(context.TODO(), username, password)
	if err != nil {
		return nil, err
	}

	return &auth.AuthedUser{
		UUID:    user.Uuid,
		Name:    user.Name,
		Account: user.Account,
		Status:  user.Status,
	}, nil
}

func New(conf *config.Config, srv service.UserService) auth.Auth {
	return &DatabaseAuth{srv: srv, conf: conf}
}

func (au DatabaseAuth) Login(username, password string) (*auth.AuthedUser, error) {
	user, err := au.srv.LoadUserByAccount(context.TODO(), username)
	if err != nil {
		return nil, err
	}

	// 弱密码模式
	if !au.conf.WeakPasswordMode {
		if user.Password == "" || !bcrypt.Match(password, user.Password) {
			return nil, errors.New("密码错误")
		}
	}

	return &auth.AuthedUser{
		UUID:    user.Uuid,
		Name:    user.Name,
		Account: user.Account,
		Status:  user.Status,
	}, nil
}

func (au DatabaseAuth) Users() ([]auth.AuthedUser, error) {
	return []auth.AuthedUser{}, nil
}
