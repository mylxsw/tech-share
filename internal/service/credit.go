package service

import (
	"context"
	"database/sql"
	"sort"
	"time"

	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
)

type CreditService interface {
	CreditRanks(ctx context.Context, startAt time.Time) (CreditRanks, error)
}

type CreditRank struct {
	UserID  int64             `json:"user_id"`
	Name    string            `json:"name"`
	Account string            `json:"account"`
	Status  int8              `json:"status"`
	Credit  int64             `json:"credit"`
	Rank    int64             `json:"rank"`
	Shares  []CreditRankShare `json:"shares"`
}

type CreditRanks []CreditRank

func (t CreditRanks) Len() int { return len(t) }
func (t CreditRanks) Less(i, j int) bool {
	if t[i].Credit == t[j].Credit {
		return t[i].Account < t[j].Account
	}

	return t[i].Credit > t[j].Credit
}
func (t CreditRanks) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

type CreditRankShare struct {
	ShareID int64  `json:"share_id"`
	Subject string `json:"subject"`
}

func NewCreditService(cc infra.Resolver, db *sql.DB) CreditService {
	return &creditService{
		cc: cc,
		db: db,
	}
}

type creditService struct {
	cc infra.Resolver
	db *sql.DB
}

func (srv *creditService) CreditRanks(ctx context.Context, startAt time.Time) (CreditRanks, error) {
	shares, err := model.NewShareModel(srv.db).Get(query.Builder().
		WhereNotNull(model.ShareFieldShareUserId).
		Where(model.ShareFieldStatus, ShareStatusFinished).
		Where(model.ShareFieldShareAt, ">=", startAt))
	if err != nil {
		return nil, err
	}

	credits := make(CreditRanks, 0)
	srv.cc.Must(coll.MustNew(shares).GroupBy(func(s model.Share) int64 { return s.ShareUserId.ValueOrZero() }).Map(func(ss []interface{}, key int64) CreditRank {
		res := CreditRank{}
		res.UserID = key
		res.Shares = make([]CreditRankShare, 0)
		res.Credit = int64(len(ss))

		for _, s := range ss {
			share := s.(model.Share)
			res.Shares = append(res.Shares, CreditRankShare{ShareID: share.Id.ValueOrZero(), Subject: share.Subject.ValueOrZero()})
		}

		return res
	}).AsArray().Filter(func(cr CreditRank) bool { return cr.Credit > 0 }).All(&credits))

	users, err := model.NewUserModel(srv.db).Get(query.Builder().Select(
		model.UserFieldId, model.UserFieldName, model.UserFieldAccount))
	if err != nil {
		return credits, err
	}

	usersMap := make(map[int64]model.User)
	srv.cc.Must(coll.MustNew(users).
		AsMap(func(u model.User) int64 { return u.Id.ValueOrZero() }).
		All(&usersMap))

	var rank int64 = 1
	var lastCredit int64

	sort.Sort(credits)
	for i, cre := range credits {
		credits[i].Name = usersMap[cre.UserID].Name.ValueOrZero()
		credits[i].Account = usersMap[cre.UserID].Account.ValueOrZero()
		credits[i].Status = int8(usersMap[cre.UserID].Status.ValueOrZero())

		if cre.Credit != lastCredit {
			rank = int64(i) + 1
		}

		credits[i].Rank = rank
		lastCredit = cre.Credit
	}

	// 第二次排序，解决第一次排序时没有 account 字段的问题
	sort.Sort(credits)
	return credits, nil
}
