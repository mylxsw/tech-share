package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
	"gopkg.in/guregu/null.v3"
)

const (
	ShareStatusVoting   int64 = 0
	ShareStatusPlaned   int64 = 1
	ShareStatusFinished int64 = 2

	RelTypeLike int64 = 1
	RelTypeJoin int64 = 2
)

type ShareService interface {
	// GetByID get a share by id
	GetByID(ctx context.Context, id int64) (model.Share, error)
	// CreateShare create a share
	CreateShare(ctx context.Context, share Share) (int64, error)
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
	// LikeShare user join a share or cancel
	JoinShare(ctx context.Context, id int64, userID int64, positive bool) (bool, error)
}

func NewShareService(cc infra.Resolver, db *sql.DB) ShareService {
	return &shareService{db: db, cc: cc}
}

type shareService struct {
	db *sql.DB
	cc infra.Resolver
}

func (p shareService) GetByID(ctx context.Context, id int64) (model.Share, error) {
	return model.NewShareModel(p.db).First(query.Builder().Where("id", id))
}

type Share struct {
	ShareUpdateFields
	ID           int64     `json:"id"`
	Status       int       `json:"status"`
	Note         string    `json:"note"`
	LikeCount    int64     `json:"like_count"`
	JoinCount    int64     `json:"join_count"`
	CreateUserID int64     `json:"create_user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ShareUpdateFields struct {
	Subject       string `json:"subject"`
	SubjectType   string `json:"subject_type"`
	Description   string `json:"description"`
	ShareUserName string `json:"share_user"`
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
		return Share{
			ShareUpdateFields: ShareUpdateFields{
				Subject:       s.Subject.ValueOrZero(),
				SubjectType:   s.SubjectType.ValueOrZero(),
				Description:   s.Description.ValueOrZero(),
				ShareUserName: s.ShareUser.ValueOrZero(),
			},
			ID:           s.Id.ValueOrZero(),
			CreateUserID: s.CreateUserId.ValueOrZero(),
			Note:         s.Note.ValueOrZero(),
			LikeCount:    s.LikeCount.ValueOrZero(),
			JoinCount:    s.JoinCount.ValueOrZero(),
			CreatedAt:    s.CreatedAt.ValueOrZero(),
			UpdatedAt:    s.UpdatedAt.ValueOrZero(),
		}
	}).All(&results)

	return results, meta, nil
}

func (p shareService) CreateShare(ctx context.Context, share Share) (int64, error) {
	return model.NewShareModel(p.db).Create(query.KV{
		"subject":        share.Subject,
		"subject_type":   share.SubjectType,
		"description":    share.Description,
		"share_user":     share.ShareUserName,
		"create_user_id": share.CreateUserID,
		"like_count":     share.LikeCount,
		"join_count":     share.JoinCount,
		"status":         share.Status,
	})
}

func (p shareService) UpdateShare(ctx context.Context, id int64, share Share) error {
	old, err := model.NewShareModel(p.db).First(query.Builder().Where("id", id))
	if err != nil {
		return err
	}

	old.Subject = null.StringFrom(share.Subject)
	old.SubjectType = null.StringFrom(share.SubjectType)
	old.Description = null.StringFrom(share.Description)
	old.ShareUser = null.StringFrom(share.ShareUserName)

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

func (p shareService) shareOpt(ctx context.Context, db query.Database, id int64, userID int64, relType int64, positive bool) (bool, error) {
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
			RelType: null.IntFrom(relType),
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
func (p shareService) updateShareRelTypeCount(ctx context.Context, db query.Database, id int64, relType int64) error {
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

func relTypeToField(relType int64) string {
	if relType == RelTypeLike {
		return "like_count"
	}

	return "join_count"
}
