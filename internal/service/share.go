package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/jinzhu/copier"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
	"gopkg.in/guregu/null.v3"
)

var validate = validator.New()

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
	GetShares(ctx context.Context, filter ShareFilter, sortDesc bool, page int64, perPage int64) ([]ShareExt, query.PaginateMeta, error)
	// FinishShare set a share as finished status
	FinishShare(ctx context.Context, shareID int64, sf ShareFinishFields) (bool, error)

	// UpdateShare update a share
	UpdateShare(ctx context.Context, id int64, share ShareUpdateFields) error
	// RemoveShare delete a share
	RemoveShare(ctx context.Context, id int64) (bool, error)
	// LikeShare user like a share or cancel
	LikeShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error)
	// JoinShare user join a share or cancel
	JoinShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error)

	// IsUserLikeOrJoinShares return whether the user join or like the shares
	IsUserLikeOrJoinShares(ctx context.Context, userID int64, shareIDs []int64) (map[int64]UserLikeOrJoinShare, error)
}

func NewShareService(cc infra.Resolver, db *sql.DB) ShareService {
	return &shareService{db: db, cc: cc}
}

type shareService struct {
	db *sql.DB
	cc infra.Resolver
}

func (p shareService) GetShareByID(ctx context.Context, id int64) (*ShareDetail, error) {
	share, err := model.NewShareModel(p.db).First(query.Builder().Where(model.ShareFieldId, id))
	if err != nil {
		return nil, err
	}

	res := ShareDetail{}
	_ = copier.Copy(&res.Share, share.ToSharePlain())

	if plan, err := share.SharePlan().First(); err == nil {
		res.Plan = &Plan{}
		_ = copier.Copy(res.Plan, plan.ToSharePlanPlain())
	}

	if share.Attachments.ValueOrZero() != "" {
		attachments, err := getAttachmentsByShare(ctx, p.db, share)
		if err != nil {
			log.With(share).Errorf("query attachments for share failed: %v", err)
		}

		res.Attachments = attachments
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

func (p shareService) FinishShare(ctx context.Context, shareID int64, sf ShareFinishFields) (bool, error) {
	if err := validate.Struct(sf); err != nil {
		return false, NewValidateError(err)
	}

	ok := false
	err := eloquent.Transaction(p.db, func(tx query.Database) error {
		share, err := model.NewShareModel(tx).First(query.Builder().Where(model.ShareFieldId, shareID))
		if err != nil {
			return err
		}

		if int8(share.Status.ValueOrZero()) == ShareStatusVoting {
			return NewValidateError(fmt.Errorf("the share can not be finished, you need create a plan first"))
		}

		// if int8(share.Status.ValueOrZero()) == ShareStatusFinished {
		// 	return nil
		// }

		share.Status = null.IntFrom(int64(ShareStatusFinished))

		var attas []string
		_ = coll.MustNew(sf.Attachments).Map(func(val int64) string { return strconv.Itoa(int(val)) }).All(&attas)

		share.Attachments = null.StringFrom(strings.Join(attas, ","))
		if err := share.Save(model.ShareFieldStatus, model.ShareFieldAttachments); err != nil {
			return err
		}

		plan, err := share.SharePlan().First()
		if err != nil {
			return err
		}

		plan.RealDuration = null.IntFrom(sf.RealDuration)
		plan.Note = null.StringFrom(sf.Note)
		if err := plan.Save(model.SharePlanFieldRealDuration, model.SharePlanFieldNote); err != nil {
			return err
		}

		ok = true
		return nil
	})

	return ok, err
}

func (p shareService) GetShares(ctx context.Context, filter ShareFilter, sortDesc bool, page int64, perPage int64) ([]ShareExt, query.PaginateMeta, error) {
	statusSort := "asc"
	if sortDesc {
		statusSort = "desc"
	}

	condition := query.Builder().OrderBy(model.ShareFieldStatus, statusSort).OrderBy(model.ShareFieldId, "desc")
	if filter.Creator > 0 {
		condition = condition.Where(model.ShareFieldCreateUserId, filter.Creator)
	}
	if len(filter.Statuses) > 0 {
		condition = condition.WhereIn(model.ShareFieldStatus, sliceToInterface(filter.Statuses)...)
	}
	if filter.Type != "" {
		condition = condition.Where(model.ShareFieldSubjectType, filter.Type)
	}

	return p.getShares(ctx, condition, page, perPage)
}

func (p shareService) getShares(ctx context.Context, qb query.SQLBuilder, page, perPage int64) ([]ShareExt, query.PaginateMeta, error) {
	if page <= 0 {
		return nil, query.PaginateMeta{}, fmt.Errorf("invalid page")
	}
	if perPage <= 0 || perPage > 100 {
		return nil, query.PaginateMeta{}, fmt.Errorf("invalid per_page, must be 1-100")
	}

	shares, meta, err := model.NewShareModel(p.db).Paginate(page, perPage, qb)
	if err != nil {
		return nil, meta, err
	}

	sharePlanRefs := make(map[int64]model.SharePlan)
	var shareIds []int64
	_ = coll.Map(shares, &shareIds, func(s model.Share) int64 { return s.Id.ValueOrZero() })
	if len(shareIds) > 0 {
		plans, err := model.NewSharePlanModel(p.db).Get(query.Builder().WhereIn(model.SharePlanFieldShareId, sliceToInterface(shareIds)...))
		if err == nil {
			_ = coll.MustNew(plans).AsMap(func(plan model.SharePlan) int64 { return plan.ShareId.ValueOrZero() }).All(&sharePlanRefs)
		}
	}

	var results []ShareExt
	_ = coll.MustNew(shares).Map(func(s model.Share) ShareExt {
		var share ShareExt
		_ = copier.Copy(&share, s.ToSharePlain())

		if sp, ok := sharePlanRefs[share.Id]; ok {
			share.ShareAt = sp.ShareAt.ValueOrZero()
			share.ShareRoom = sp.ShareRoom.ValueOrZero()
			share.PlanDuration = sp.PlanDuration.ValueOrZero()
		}

		return share
	}).All(&results)

	return results, meta, nil
}

func (p shareService) CreateShare(ctx context.Context, share Share) (int64, error) {
	var shareID int64
	share.Status = ShareStatusVoting
	if err := validate.Struct(share); err != nil {
		return shareID, NewValidateError(err)
	}

	err := eloquent.Transaction(p.db, func(tx query.Database) error {
		exist, err := model.NewShareModel(tx).Exists(query.Builder().Where(model.ShareFieldSubject, share.Subject))
		if err != nil && err != query.ErrNoResult {
			return err
		}

		if exist {
			return NewValidateError(fmt.Errorf("the subject already existed"))
		}

		var shareP model.SharePlain
		_ = copier.Copy(&shareP, share)

		shareID, err = model.NewShareModel(tx).Create(query.KV{
			model.ShareFieldSubject:      share.Subject,
			model.ShareFieldSubjectType:  share.SubjectType,
			model.ShareFieldDescription:  share.Description,
			model.ShareFieldShareUser:    share.ShareUser,
			model.ShareFieldCreateUserId: share.CreateUserId,
			model.ShareFieldLikeCount:    share.LikeCount,
			model.ShareFieldJoinCount:    share.JoinCount,
			model.ShareFieldStatus:       share.Status,
		})

		return err
	})

	return shareID, err
}

func (p shareService) CreateOrUpdatePlan(ctx context.Context, shareID int64, plan PlanUpdateFields) (int64, error) {
	var planID int64

	if err := validate.Struct(plan); err != nil {
		return planID, NewValidateError(err)
	}

	err := eloquent.Transaction(p.db, func(tx query.Database) error {
		// query share
		share, err := model.NewShareModel(tx).First(query.Builder().Where(model.ShareFieldId, shareID))
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
			model.SharePlanFieldShareRoom,
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
		planExi.ShareRoom = sp.ShareRoom
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
		share, err := model.NewShareModel(tx).First(query.Builder().Where(model.ShareFieldId, shareID))
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

		affected, err = model.NewSharePlanModel(tx).Delete(query.Builder().Where(model.SharePlanFieldShareId, share.Id.ValueOrZero()))
		return err
	})

	return affected > 0, err
}

func (p shareService) UpdateShare(ctx context.Context, id int64, share ShareUpdateFields) error {
	if err := validate.Struct(share); err != nil {
		return NewValidateError(err)
	}

	exist, err := model.NewShareModel(p.db).Exists(query.Builder().
		Where("id", "!=", id).
		Where(model.ShareFieldSubject, share.Subject))
	if err != nil && err != query.ErrNoResult {
		return err
	}

	if exist {
		return NewValidateError(fmt.Errorf("the subject already existed"))
	}

	old, err := model.NewShareModel(p.db).First(query.Builder().Where(model.ShareFieldId, id))
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
	share, err := model.NewShareModel(db).First(query.Builder().Where(model.ShareFieldId, id))
	if err != nil {
		return false, err
	}

	rel, err := share.ShareUserRels().First(query.Builder().
		Where(model.ShareUserRelFieldUserId, userID).
		Where(model.ShareUserRelFieldRelType, relType))
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
	relCount, err := model.NewShareUserRelModel(db).Count(query.Builder().
		Where(model.ShareUserRelFieldShareId, id).
		Where(model.ShareUserRelFieldRelType, relType))
	if err != nil {
		return err
	}

	if _, err := model.NewShareModel(db).UpdateFields(
		query.KV{relTypeToField(relType): relCount},
		query.Builder().Where(model.ShareFieldId, id),
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

func (p shareService) IsUserLikeOrJoinShares(ctx context.Context, userID int64, shareIDs []int64) (map[int64]UserLikeOrJoinShare, error) {
	if len(shareIDs) == 0 {
		return make(map[int64]UserLikeOrJoinShare), nil
	}

	rels, err := model.NewShareUserRelModel(p.db).Get(query.Builder().
		Where(model.ShareUserRelFieldUserId, userID).
		WhereIn(model.ShareUserRelFieldShareId, sliceToInterface(shareIDs)...))
	if err != nil {
		return nil, err
	}

	results := make(map[int64]UserLikeOrJoinShare)
	_ = coll.MustNew(rels).GroupBy(func(rel model.ShareUserRel) int64 {
		return rel.ShareId.ValueOrZero()
	}).Map(func(ss []interface{}, shareID int64) UserLikeOrJoinShare {
		res := UserLikeOrJoinShare{ShareID: shareID}
		for _, s := range ss {
			if s.(model.ShareUserRel).RelType.ValueOrZero() == int64(RelTypeJoin) {
				res.Join = true
			} else {
				res.Like = true
			}
		}

		return res
	}).All(&results)
	return results, nil
}

func relTypeToField(relType int8) string {
	if relType == RelTypeLike {
		return model.ShareFieldLikeCount
	}

	return model.ShareFieldJoinCount
}

func sliceToInterface(items interface{}) []interface{} {
	res := make([]interface{}, 0)
	_ = coll.Map(items, &res)
	return res
}
