package service

import (
	"context"
	"database/sql"

	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/tech-share/internal/service/model"
)

type AttachmentService interface {
	GetByID(ctx context.Context, id int64) (model.Attachment, error)
	GetByShareID(ctx context.Context, shareID int64) ([]model.Attachment, error)
}

func NewAttachmentService(cc infra.Resolver, db *sql.DB) AttachmentService {
	return &attachmentService{db: db, cc: cc}
}

type attachmentService struct {
	db *sql.DB
	cc infra.Resolver
}

func (p attachmentService) GetByShareID(ctx context.Context, shareID int64) ([]model.Attachment, error) {
	return model.NewAttachmentModel(p.db).Get(query.Builder().Where("share_id", shareID))
}

func (p attachmentService) GetByID(ctx context.Context, id int64) (model.Attachment, error) {
	return model.NewAttachmentModel(p.db).First(query.Builder().Where("id", id))
}
