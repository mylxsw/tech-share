package service

import (
	"context"
	"database/sql"

	"github.com/jinzhu/copier"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
	"gopkg.in/guregu/null.v3"
)

type UserService interface {
	Users(ctx context.Context) ([]UserBasic, error)
	LoadUser(ctx context.Context, uuid string, userInfo UserInfo) (*UserInfo, error)
}

func NewUserService(cc infra.Resolver, db *sql.DB) UserService {
	return &userService{cc: cc, db: db}
}

type userService struct {
	cc infra.Resolver
	db *sql.DB
}

type UserInfo struct {
	Id      int64
	Uuid    string
	Name    string
	Account string
	Status  int8
}

const (
	UserStatusEnabled  int8 = 1
	UserStatusDisabled int8 = 0
)

type UserBasic struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}

func (s userService) Users(ctx context.Context) ([]UserBasic, error) {
	users, err := model.NewUserModel(s.db).Get(query.Builder().Where(model.UserFieldStatus, UserStatusEnabled))
	if err != nil {
		return nil, err
	}

	results := make([]UserBasic, 0)
	_ = coll.MustNew(users).Map(func(user model.User) UserBasic {
		var res UserBasic
		_ = copier.Copy(&res, user)
		return res
	}).All(&results)

	return results, nil
}

func (s userService) LoadUser(ctx context.Context, uuid string, userInfo UserInfo) (*UserInfo, error) {
	user, err := model.NewUserModel(s.db).First(query.Builder().Where(model.UserFieldUuid, uuid))
	if err != nil && err != query.ErrNoResult {
		return nil, err
	}

	if err == query.ErrNoResult {
		userID, err := model.NewUserModel(s.db).Create(query.KV{
			model.UserFieldUuid:    uuid,
			model.UserFieldName:    userInfo.Name,
			model.UserFieldAccount: userInfo.Account,
			model.UserFieldStatus:  userInfo.Status,
		})
		if err != nil {
			return nil, err
		}

		userInfo.Id = userID
		userInfo.Uuid = uuid

		return &userInfo, nil
	}

	user.Name = null.StringFrom(userInfo.Name)
	user.Account = null.StringFrom(userInfo.Account)
	user.Status = null.IntFrom(int64(userInfo.Status))

	if err := user.Save(model.UserFieldName, model.UserFieldAccount, model.UserFieldStatus); err != nil {
		return nil, err
	}

	res := UserInfo{}
	_ = copier.Copy(&res, user)

	return &res, nil
}
