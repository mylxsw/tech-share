package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
	"gopkg.in/guregu/null.v3"
)

const (
	ShareStatusVoting   int8 = 0
	ShareStatusPlaned   int8 = 1
	ShareStatusFinished int8 = 2

	RelTypeLike int8 = 1
	RelTypeJoin int8 = 2
)

type ShareService interface {
	// GetShareByID get a share by id
	GetShareByID(ctx context.Context, id int64) (*ShareDetail, error)
	// CreateShare create a share
	CreateShare(ctx context.Context, share Share) (int64, error)
	// CreateOrUpdatePlan create or update a share plan
	CreateOrUpdatePlan(ctx context.Context, shareID int64, plan PlanUpdateFields) (int64, error)
	// RemovePlan delete a plan for a share
	RemovePlan(ctx context.Context, shareID int64) (bool, error)
	// GetShares return all shares
	GetShares(ctx context.Context, page int64, perPage int64) ([]Share, query.PaginateMeta, error)
	// GetSharesForUser return all shares for a user
	GetSharesForUser(ctx context.Context, userID int64, page int64, perPage int64) ([]Share, query.PaginateMeta, error)
	// GetSharesPlaned return all shares has been planed
	GetSharesPlaned(ctx context.Context, page int64, perPage int64) ([]Share, query.PaginateMeta, error)

	// UpdateShare update a share
	UpdateShare(ctx context.Context, id int64, share Share) error
	// RemoveShare delete a share
	RemoveShare(ctx context.Context, id int64) (bool, error)
	// LikeShare user like a share or cancel
	LikeShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error)
	// JoinShare user join a share or cancel
	JoinShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error)
}

func NewShareService(cc infra.Resolver, db *sql.DB) ShareService {
	return &shareService{db: db, cc: cc}
}

type shareService struct {
	db *sql.DB
	cc infra.Resolver
}

func (p shareService) GetShareByID(ctx context.Context, id int64) (*ShareDetail, error) {
	share, err := model.NewShareModel(p.db).First(query.Builder().Where("id", id))
	if err != nil {
		return nil, err
	}

	res := ShareDetail{}
	_ = copier.Copy(&res.Share, share.ToSharePlain())

	if plan, err := share.SharePlan().First(); err == nil {
		res.Plan = &Plan{}
		_ = copier.Copy(res.Plan, plan.ToSharePlanPlain())
	}

	if attas, err := share.Attachments().Get(); err == nil {
		_ = coll.MustNew(attas).Map(func(atta model.Attachment) Attachment {
			var local Attachment
			_ = copier.Copy(&local, atta.ToAttachmentPlain())

			return local
		}).All(&res.Attachments)
	}

	rawSQL := `SELECT rel.user_id, rel.rel_type, user.name FROM share_user_rel rel LEFT JOIN user user ON user.id = rel.user_id WHERE rel.share_id=?`
	if relUsers, err := eloquent.DB(p.db).Query(eloquent.Raw(rawSQL, id), func(row eloquent.Scanner) (interface{}, error) {
		var user User
		if err := row.Scan(&user.ID, &user.relType, &user.Name); err != nil {
			return nil, err
		}

		return user, nil
	}); err == nil {
		groupedUsers := make(map[interface{}][]interface{})
		_ = relUsers.GroupBy(func(user User) int8 { return user.relType }).All(&groupedUsers)

		if likeUsers, ok := groupedUsers[RelTypeLike]; ok {
			_ = coll.Map(likeUsers, &res.LikeUsers)
		}
		if joinUsers, ok := groupedUsers[RelTypeJoin]; ok {
			_ = coll.Map(joinUsers, &res.JoinUsers)
		}
	}

	return &res, nil
}

type ShareDetail struct {
	Share       Share        `json:"share"`
	Plan        *Plan        `json:"plan"`
	Attachments []Attachment `json:"attachments"`
	LikeUsers   []User       `json:"like_users"`
	JoinUsers   []User       `json:"join_users"`
}

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	relType int8
}

type Attachment struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	AttaType  string    `json:"atta_type"`
	AttaPath  string    `json:"atta_path"`
	CreatedAt time.Time `json:"created_at"`
}

type Plan struct {
	PlanUpdateFields
	Id           int64     `json:"id"`
	RealDuration int64     `json:"real_duration"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PlanUpdateFields struct {
	ShareAt      time.Time `json:"share_at"`
	PlanDuration int64     `json:"plan_duration"`
	Note         string    `json:"note"`
}

type Share struct {
	ShareUpdateFields
	Id           int64     `json:"id"`
	Status       int8      `json:"status"`
	Note         string    `json:"note"`
	LikeCount    int64     `json:"like_count"`
	JoinCount    int64     `json:"join_count"`
	CreateUserID int64     `json:"create_user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ShareUpdateFields struct {
	Subject     string `json:"subject"`
	SubjectType string `json:"subject_type"`
	Description string `json:"description"`
	ShareUser   string `json:"share_user"`
}

func (p shareService) GetShares(ctx context.Context, page int64, perPage int64) ([]Share, query.PaginateMeta, error) {
	return p.getShares(ctx, query.Builder().OrderBy("id", "desc"), page, perPage)
}

func (p shareService) GetSharesForUser(ctx context.Context, userID int64, page int64, perPage int64) ([]Share, query.PaginateMeta, error) {
	return p.getShares(ctx, query.Builder().Where("create_user_id", userID).OrderBy("id", "desc"), page, perPage)
}

func (p shareService) GetSharesPlaned(ctx context.Context, page int64, perPage int64) ([]Share, query.PaginateMeta, error) {
	return p.getShares(ctx, query.Builder().Where("status", ShareStatusPlaned).OrderBy("id", "desc"), page, perPage)
}

func (p shareService) getShares(ctx context.Context, qb query.SQLBuilder, page, perPage int64) ([]Share, query.PaginateMeta, error) {
	shares, meta, err := model.NewShareModel(p.db).Paginate(page, perPage, qb)
	if err != nil {
		return nil, meta, err
	}

	var results []Share
	_ = coll.MustNew(shares).Map(func(s model.Share) Share {
		var share Share
		_ = copier.Copy(&share, s.ToSharePlain())

		return share
	}).All(&results)

	return results, meta, nil
}

func (p shareService) CreateShare(ctx context.Context, share Share) (int64, error) {
	return model.NewShareModel(p.db).Create(query.KV{
		"subject":        share.Subject,
		"subject_type":   share.SubjectType,
		"description":    share.Description,
		"share_user":     share.ShareUser,
		"create_user_id": share.CreateUserID,
		"like_count":     share.LikeCount,
		"join_count":     share.JoinCount,
		"status":         share.Status,
	})
}

func (p shareService) CreateOrUpdatePlan(ctx context.Context, shareID int64, plan PlanUpdateFields) (int64, error) {
	var planID int64
	err := eloquent.Transaction(p.db, func(tx query.Database) error {
		// query share
		share, err := model.NewShareModel(tx).First(query.Builder().Where("id", shareID))
		if err != nil {
			return err
		}

		// constraint check
		if share.Status.ValueOrZero() == int64(ShareStatusFinished) {
			return fmt.Errorf("plan can not be updated because share has been finished")
		}

		// update share status
		share.Status = null.IntFrom(int64(ShareStatusPlaned))
		if err := share.Save(); err != nil {
			return err
		}

		// onlyFields only update this fields
		onlyFields := []string{
			model.SharePlanFieldShareAt,
			model.SharePlanFieldPlanDuration,
			model.SharePlanFieldNote,
		}

		// create or update plan
		var planPlain model.SharePlanPlain
		_ = copier.Copy(&planPlain, plan)
		sp := planPlain.ToSharePlan(onlyFields...)

		planExi, err := share.SharePlan().First()
		if err != nil {
			if err != query.ErrNoResult {
				return err
			}

			id, err := share.SharePlan().Create(sp)
			if err != nil {
				return err
			}

			planID = id
			return nil
		}

		planExi.ShareAt = sp.ShareAt
		planExi.PlanDuration = sp.PlanDuration
		planExi.Note = sp.Note

		fmt.Println(planExi.StaledKV(onlyFields...))

		planID = planExi.Id.ValueOrZero()
		return planExi.Save(onlyFields...)
	})

	return planID, err
}

func (p shareService) RemovePlan(ctx context.Context, shareID int64) (bool, error) {
	var affected int64
	err := eloquent.Transaction(p.db, func(tx query.Database) error {
		share, err := model.NewShareModel(tx).First(query.Builder().Where("id", shareID))
		if err != nil {
			return err
		}

		if share.Status.ValueOrZero() == int64(ShareStatusFinished) {
			return fmt.Errorf("remove plan failed because the share has been finished")
		}

		share.Status = null.IntFrom(int64(ShareStatusVoting))
		if err := share.Save(); err != nil {
			return err
		}

		affected, err = model.NewSharePlanModel(tx).Delete(query.Builder().Where("share_id", share.Id.ValueOrZero()))
		return err
	})

	return affected > 0, err
}

func (p shareService) UpdateShare(ctx context.Context, id int64, share Share) error {
	old, err := model.NewShareModel(p.db).First(query.Builder().Where("id", id))
	if err != nil {
		return err
	}

	old.Subject = null.StringFrom(share.Subject)
	old.SubjectType = null.StringFrom(share.SubjectType)
	old.Description = null.StringFrom(share.Description)
	old.ShareUser = null.StringFrom(share.ShareUser)

	return old.Save()
}

func (p shareService) LikeShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return false, err
	}

	ok, err := p.shareOpt(ctx, tx, id, userID, RelTypeLike, positive)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return ok, nil
}

func (p shareService) JoinShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return false, err
	}

	ok, err := p.shareOpt(ctx, tx, id, userID, RelTypeJoin, positive)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return ok, nil
}

func (p shareService) shareOpt(ctx context.Context, db query.Database, id int64, userID int64, relType int8, positive bool) (bool, error) {
	share, err := model.NewShareModel(db).First(query.Builder().Where("id", id))
	if err != nil {
		return false, err
	}

	rel, err := share.ShareUserRels().First(query.Builder().Where("user_id", userID).Where("rel_type", relType))
	if err != nil && err != query.ErrNoResult {
		return false, err
	}

	if positive {
		if err != query.ErrNoResult {
			return false, nil
		}

		if _, err := share.ShareUserRels().Create(model.ShareUserRel{
			UserId:  null.IntFrom(userID),
			RelType: null.IntFrom(int64(relType)),
		}); err != nil {
			return false, err
		}

		return true, p.updateShareRelTypeCount(ctx, db, id, relType)
	}

	if err != query.ErrNoResult {
		if err := rel.Delete(); err != nil {
			return false, err
		}

		return true, p.updateShareRelTypeCount(ctx, db, id, relType)
	}

	return false, nil
}
func (p shareService) updateShareRelTypeCount(ctx context.Context, db query.Database, id int64, relType int8) error {
	relCount, err := model.NewShareUserRelModel(db).Count(query.Builder().Where("share_id", id).Where("rel_type", relType))
	if err != nil {
		return err
	}

	if _, err := model.NewShareModel(db).UpdateFields(
		query.KV{relTypeToField(relType): relCount},
		query.Builder().Where("id", id),
	); err != nil {
		return err
	}

	return nil
}

func (p shareService) RemoveShare(ctx context.Context, id int64) (bool, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return false, err
	}

	affected, err := model.NewShareModel(tx).DeleteById(id)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	// TODO whether need to delete all related resources

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return affected > 0, nil
}

func relTypeToField(relType int8) string {
	if relType == RelTypeLike {
		return "like_count"
	}

	return "join_count"
}
