package controller

import (
	"context"
	"strconv"

	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

type ShareController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewShareController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &ShareController{cc: cc, conf: conf}
}

func (ctl ShareController) Register(router web.Router) {
	router.Group("shares/", func(router web.Router) {
		router.Get("/", ctl.Shares)
		router.Get("/{id}/", ctl.Share)
		router.Post("/", ctl.CreateShare)
		router.Post("/{id}/", ctl.UpdateShare)
		router.Delete("/{id}/", ctl.DeleteShare)

		router.Post("/{id}/like/", ctl.LikeShare)
		router.Delete("/{id}/like/", ctl.DislikeShare)

		router.Post("/{id}/join/", ctl.JoinShare)
		router.Delete("/{id}/join/", ctl.LeaveShare)

		router.Post("/{id}/plan/", ctl.CreateOrUpdateSharePlan)
		router.Put("/{id}/plan/", ctl.CreateOrUpdateSharePlan)
		router.Delete("/{id}/plan/", ctl.CancelSharePlan)

		router.Post("/{id}/finish/", ctl.FinishShare)
	})
}

// Shares return all shares
func (ctl ShareController) Shares(ctx web.Context, req web.Request, shareSrv service.ShareService) (*PaginateRes, error) {
	page := req.Int64Input("page", 1)
	perPage := req.Int64Input("per_page", 20)

	shares, meta, err := shareSrv.GetShares(context.TODO(), page, perPage)
	if err != nil {
		return nil, err
	}

	return &PaginateRes{
		Page: meta,
		Data: shares,
	}, nil
}

// CreateShare create a share
func (ctl ShareController) CreateShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	var share service.ShareUpdateFields
	if err := req.Unmarshal(&share); err != nil {
		return err
	}

	_, err := shareSrv.CreateShare(context.TODO(), service.Share{
		ShareUpdateFields: share,
		CreateUserID:      currentUserID(req).Id,
	})
	return err
}

// UpdateShare update an existed share
func (ctl ShareController) UpdateShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	var share service.ShareUpdateFields
	if err := req.Unmarshal(&share); err != nil {
		return err
	}

	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	return shareSrv.UpdateShare(context.TODO(), int64(id), service.Share{
		ShareUpdateFields: share,
		CreateUserID:      currentUserID(req).Id,
	})
}

// DeleteShare remove a share
func (ctl ShareController) DeleteShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.RemoveShare(context.TODO(), int64(id))
	return err
}

// Share get a share detail
func (ctl ShareController) Share(ctx web.Context, req web.Request, shareSrv service.ShareService) (*service.ShareDetail, error) {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return nil, err
	}

	return shareSrv.GetShareByID(context.TODO(), int64(id))
}

// LikeShare give a vote to a share
func (ctl ShareController) LikeShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.LikeShare(context.TODO(), int64(id), currentUserID(req).Id, true)
	return err
}

// DislikeShare cancel a vote to share
func (ctl ShareController) DislikeShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.LikeShare(context.TODO(), int64(id), currentUserID(req).Id, false)
	return err
}

// JoinShare join a share as an audience
func (ctl ShareController) JoinShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.JoinShare(context.TODO(), int64(id), currentUserID(req).Id, true)
	return err
}

// LeaveShare leave a share
func (ctl ShareController) LeaveShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.JoinShare(context.TODO(), int64(id), currentUserID(req).Id, false)
	return err
}

// CreateOrUpdateSharePlan create or update a share plan
func (ctl ShareController) CreateOrUpdateSharePlan(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	var plan service.PlanUpdateFields
	if err := req.Unmarshal(&plan); err != nil {
		return err
	}

	_, err = shareSrv.CreateOrUpdatePlan(context.TODO(), int64(id), plan)
	return err
}

// CancelSharePlan cancel a share plan
func (ctl ShareController) CancelSharePlan(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.RemovePlan(context.TODO(), int64(id))
	return err
}

// FinishShare set a share as finished
func (ctl ShareController) FinishShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	var sf service.ShareFinishFields
	if err := req.Unmarshal(&sf); err != nil {
		return err
	}

	_, err = shareSrv.ShareFinish(context.TODO(), int64(id), sf)
	return err
}
