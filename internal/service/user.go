package service

import (
	"context"
	"database/sql"
	"sort"

	"github.com/jinzhu/copier"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
	"gopkg.in/guregu/null.v3"
)

type UserService interface {
	Users(ctx context.Context) (UserBasics, error)
	LoadUser(ctx context.Context, uuid string, userInfo UserInfo) (*UserInfo, error)
	LoadUserByAccount(ctx context.Context, account string) (*UserInfo, error)
}

func NewUserService(cc infra.Resolver, db *sql.DB) UserService {
	return &userService{cc: cc, db: db}
}

type userService struct {
	cc infra.Resolver
	db *sql.DB
}

type UserInfo struct {
	Id       int64  `json:"id"`
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Account  string `json:"account"`
	Status   int8   `json:"status"`
	Password string `json:"-"`
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

type UserBasics []UserBasic

func (t UserBasics) Len() int           { return len(t) }
func (t UserBasics) Less(i, j int) bool { return t[i].Account < t[j].Account }
func (t UserBasics) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (s userService) Users(ctx context.Context) (UserBasics, error) {
	users, err := model.NewUserModel(s.db).Get(query.Builder().Where(model.UserFieldStatus, UserStatusEnabled))
	if err != nil {
		return nil, err
	}

	results := make(UserBasics, 0)
	_ = coll.MustNew(users).Map(func(user model.User) UserBasic {
		var res UserBasic
		_ = copier.Copy(&res, user)
		return res
	}).All(&results)

	sort.Sort(results)
	return results, nil
}

func (s userService) LoadUserByAccount(ctx context.Context, account string) (*UserInfo, error) {
	user, err := model.NewUserModel(s.db).First(query.Builder().Where(model.UserFieldAccount, account))
	if err != nil && err != query.ErrNoResult {
		return nil, err
	}

	res := UserInfo{}
	_ = copier.Copy(&res, user)

	return &res, nil
}

func (s userService) LoadUser(ctx context.Context, uuid string, userInfo UserInfo) (*UserInfo, error) {
	user, err := model.NewUserModel(s.db).First(query.Builder().Where(model.UserFieldUuid, uuid))
	if err != nil && err != query.ErrNoResult {
		return nil, err
	}

	if err == query.ErrNoResult {
		userID, err := model.NewUserModel(s.db).Create(query.KV{
			model.UserFieldUuid:     uuid,
			model.UserFieldName:     userInfo.Name,
			model.UserFieldAccount:  userInfo.Account,
			model.UserFieldStatus:   userInfo.Status,
			model.UserFieldPassword: userInfo.Password,
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
