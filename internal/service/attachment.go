package service

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/go-utils/str"
	"github.com/mylxsw/tech-share/internal/service/model"
)

type AttachmentService interface {
	GetByID(ctx context.Context, id int64) (*Attachment, error)
	GetByShareID(ctx context.Context, shareID int64) ([]Attachment, error)
	CreateAttachment(ctx context.Context, atta Attachment) (int64, error)
}

func NewAttachmentService(cc infra.Resolver, db *sql.DB) AttachmentService {
	return &attachmentService{db: db, cc: cc}
}

type attachmentService struct {
	db *sql.DB
	cc infra.Resolver
}

func getAttachmentsByShare(ctx context.Context, db *sql.DB, share model.Share) ([]Attachment, error) {
	if share.Attachments.ValueOrZero() == "" {
		return []Attachment{}, nil
	}

	var attaIDs []int64
	_ = coll.MustNew(strings.Split(share.Attachments.ValueOrZero(), ",")).
		Map(func(s string) string { return strings.TrimSpace(s) }).
		Filter(func(s string) bool { return s != "" }).
		Map(func(s string) int64 { res, _ := strconv.Atoi(s); return int64(res) }).All(&attaIDs)

	attas, err := model.NewAttachmentModel(db).Get(query.Builder().WhereIn(model.AttachmentFieldId, sliceToInterface(attaIDs)...))
	if err != nil {
		return nil, err
	}

	var res []Attachment
	_ = coll.Map(attas, &res, func(atta model.Attachment) Attachment {
		a := Attachment{}
		_ = copier.Copy(&a, atta.ToAttachmentPlain())

		return a
	})

	return res, nil
}

func (p attachmentService) GetByShareID(ctx context.Context, shareID int64) ([]Attachment, error) {
	share, err := model.NewShareModel(p.db).First(query.Builder().Where(model.ShareFieldId, shareID))
	if err != nil {
		return nil, err
	}

	return getAttachmentsByShare(ctx, p.db, share)
}

func (p attachmentService) GetByID(ctx context.Context, id int64) (*Attachment, error) {
	atta, err := model.NewAttachmentModel(p.db).First(query.Builder().Where(model.AttachmentFieldId, id))
	if err != nil {
		return nil, err
	}

	res := Attachment{}
	_ = copier.Copy(&res, atta.ToAttachmentPlain())

	return &res, nil
}

func (p attachmentService) CreateAttachment(ctx context.Context, atta Attachment) (int64, error) {
	attaP := model.AttachmentPlain{}
	_ = copier.Copy(&attaP, atta)

	return model.NewAttachmentModel(p.db).Save(
		attaP.ToAttachment(str.Exclude(
			model.AttachmentFields(),
			model.AttachmentFieldId,
			model.AttachmentFieldCreatedAt,
		)...))
}
