package controller

import (
	"context"
	"strconv"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
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
		router.Get("/my/", ctl.MyShares)
		router.Get("/{id:[0-9]+}/", ctl.Share)
		router.Post("/", ctl.CreateShare)
		router.Post("/{id:[0-9]+}/", ctl.UpdateShare)
		router.Delete("/{id:[0-9]+}/", ctl.DeleteShare)

		router.Post("/{id:[0-9]+}/like/", ctl.LikeShare)
		router.Delete("/{id:[0-9]+}/like/", ctl.DislikeShare)

		router.Post("/{id:[0-9]+}/join/", ctl.JoinShare)
		router.Delete("/{id:[0-9]+}/join/", ctl.LeaveShare)

		router.Post("/{id:[0-9]+}/plan/", ctl.CreateOrUpdateSharePlan)
		router.Put("/{id:[0-9]+}/plan/", ctl.CreateOrUpdateSharePlan)
		router.Delete("/{id:[0-9]+}/plan/", ctl.CancelSharePlan)

		router.Post("/{id:[0-9]+}/finish/", ctl.FinishShare)
	})
}

// Shares return all shares
//   - status: 0 所有 1/2/3 基于状态筛选
//   - creator: 基于创建人筛选 0 - 所有人
func (ctl ShareController) Shares(ctx web.Context, req web.Request, shareSrv service.ShareService) (*PaginateRes, error) {
	page := req.Int64Input("page", 1)
	perPage := req.Int64Input("per_page", 20)

	status := req.Int64Input("status", 0)
	creator := req.Int64Input("creator", 0)

	filter := service.ShareFilter{
		Status:  int8(status),
		Creator: creator,
	}
	shares, meta, err := shareSrv.GetShares(context.TODO(), filter, page, perPage)
	if err != nil {
		return nil, err
	}

	var shareIDs []int64
	_ = coll.Map(shares, &shareIDs, func(s service.Share) int64 { return s.Id })

	userLikeOrJoinShares, err := shareSrv.IsUserLikeOrJoinShares(context.TODO(), currentUser(req).Id, shareIDs)
	if err != nil {
		log.Errorf("query user_like_or_join_shares failed: %v", err)
	}

	return &PaginateRes{
		Page: meta,
		Data: shares,
		Extra: SharesExtra{
			UserLikeOrJoin: userLikeOrJoinShares,
		},
		Search: map[string]interface{}{
			"status":   status,
			"creator":  creator,
			"page":     page,
			"per_page": perPage,
		},
	}, nil
}

type SharesExtra struct {
	UserLikeOrJoin map[int64]service.UserLikeOrJoinShare `json:"user_like_or_join"`
}

// MyShares return all shares for current user
//   - status: 0 所有 1/2/3 基于状态筛选
func (ctl ShareController) MyShares(ctx web.Context, req web.Request, shareSrv service.ShareService) (*PaginateRes, error) {
	page := req.Int64Input("page", 1)
	perPage := req.Int64Input("per_page", 20)

	status := req.Int64Input("status", 0)

	filter := service.ShareFilter{
		Status:  int8(status),
		Creator: currentUser(req).Id,
	}
	shares, meta, err := shareSrv.GetShares(context.TODO(), filter, page, perPage)
	if err != nil {
		return nil, err
	}

	return &PaginateRes{
		Page: meta,
		Data: shares,
		Search: map[string]interface{}{
			"status":   status,
			"page":     page,
			"per_page": perPage,
		},
	}, nil
}

// CreateShare create a share
func (ctl ShareController) CreateShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	var share service.ShareUpdateFields
	if err := req.Unmarshal(&share); err != nil {
		return err
	}

	if share.ShareUser == "" {
		share.ShareUser = currentUser(req).Name
	}

	_, err := shareSrv.CreateShare(context.TODO(), service.Share{
		ShareUpdateFields: share,
		CreateUserId:      currentUser(req).Id,
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

	return shareSrv.UpdateShare(context.TODO(), int64(id), share)
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

	_, err = shareSrv.LikeShare(context.TODO(), int64(id), currentUser(req).Id, true)
	return err
}

// DislikeShare cancel a vote to share
func (ctl ShareController) DislikeShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.LikeShare(context.TODO(), int64(id), currentUser(req).Id, false)
	return err
}

// JoinShare join a share as an audience
func (ctl ShareController) JoinShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.JoinShare(context.TODO(), int64(id), currentUser(req).Id, true)
	return err
}

// LeaveShare leave a share
func (ctl ShareController) LeaveShare(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	_, err = shareSrv.JoinShare(context.TODO(), int64(id), currentUser(req).Id, false)
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

	_, err = shareSrv.FinishShare(context.TODO(), int64(id), sf)
	return err
}
