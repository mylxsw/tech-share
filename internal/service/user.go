package service

import (
	"context"
	"database/sql"

	"github.com/jinzhu/copier"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
	"gopkg.in/guregu/null.v3"
)

type UserService interface {
	LoadUser(ctx context.Context, uuid string, name string) (*UserInfo, error)
}

func NewUserService(cc infra.Resolver, db *sql.DB) UserService {
	return &userService{cc: cc, db: db}
}

type userService struct {
	cc infra.Resolver
	db *sql.DB
}

type UserInfo struct {
	Id   int64
	Uuid string
	Name string
}

func (s userService) LoadUser(ctx context.Context, uuid string, name string) (*UserInfo, error) {
	user, err := model.NewUserModel(s.db).First(query.Builder().Where(model.UserFieldUuid, uuid))
	if err != nil && err != query.ErrNoResult {
		return nil, err
	}

	if err == query.ErrNoResult {
		userID, err := model.NewUserModel(s.db).Create(query.KV{
			model.UserFieldUuid: uuid,
			model.UserFieldName: name,
		})
		if err != nil {
			return nil, err
		}

		return &UserInfo{Id: userID, Uuid: uuid, Name: name}, nil
	}

	user.Name = null.StringFrom(name)
	if err := user.Save(model.UserFieldName); err != nil {
		return nil, err
	}

	res := UserInfo{}
	_ = copier.Copy(&res, user)

	return &res, nil
}
