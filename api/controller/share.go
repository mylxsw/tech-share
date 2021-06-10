package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	event2 "github.com/mylxsw/tech-share/internal/event"
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
		router.Get("/recently/", ctl.RecentlyShares)
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

// RecentlyShares return all shares recently
//   - status: 0 所有 1/2/3 基于状态筛选
//   - type: 分享类型
func (ctl ShareController) RecentlyShares(ctx web.Context, req web.Request, shareSrv service.ShareService) (*PaginateRes, error) {
	filter := service.ShareFilter{
		Statuses: []int8{service.ShareStatusPlaned, service.ShareStatusVoting},
		Type:     req.Input("type"),
	}

	return ctl.getShares(context.TODO(), shareSrv, filter, true, req.Int64Input("page", 1), req.Int64Input("per_page", 20), currentUser(req).Id)
}

// Shares return all shares
//   - status: 0 所有 1/2/3 基于状态筛选
//   - creator: 基于创建人筛选 0 - 所有人
//   - type: 分享类型
func (ctl ShareController) Shares(ctx web.Context, req web.Request, shareSrv service.ShareService) (*PaginateRes, error) {
	filter := service.ShareFilter{
		Statuses: statusFilter(req),
		Creator:  req.Int64Input("creator", 0),
		Type:     req.Input("type"),
	}

	return ctl.getShares(context.TODO(), shareSrv, filter, false, req.Int64Input("page", 1), req.Int64Input("per_page", 20), currentUser(req).Id)
}

func statusFilter(req web.Request) []int8 {
	status := int8(req.Int64Input("status", 0))
	if status <= 0 {
		return []int8{}
	}

	return []int8{status}
}

type SharesExtra struct {
	UserLikeOrJoin map[int64]service.UserLikeOrJoinShare `json:"user_like_or_join"`
}

// MyShares return all shares for current user
//   - status: 0 所有 1/2/3 基于状态筛选
//   - type: 分享类型
func (ctl ShareController) MyShares(ctx web.Context, req web.Request, shareSrv service.ShareService) (*PaginateRes, error) {
	currentUserID := currentUser(req).Id
	filter := service.ShareFilter{
		Statuses: statusFilter(req),
		Creator:  currentUserID,
		Type:     req.Input("type"),
	}

	return ctl.getShares(context.TODO(), shareSrv, filter, false, req.Int64Input("page", 1), req.Int64Input("per_page", 20), currentUserID)
}

func (ctl ShareController) getShares(ctx context.Context, shareSrv service.ShareService, filter service.ShareFilter, sortDesc bool, page, perPage int64, currentUserID int64) (*PaginateRes, error) {
	shares, meta, err := shareSrv.GetShares(context.TODO(), filter, sortDesc, page, perPage)
	if err != nil {
		return nil, err
	}

	var shareIDs []int64
	_ = coll.Map(shares, &shareIDs, func(s service.ShareExt) int64 { return s.Id })

	userLikeOrJoinShares, err := shareSrv.IsUserLikeOrJoinShares(context.TODO(), currentUserID, shareIDs)
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
			"status":   returnStatusFilter(filter),
			"type":     filter.Type,
			"creator":  filter.Creator,
			"page":     page,
			"per_page": perPage,
		},
	}, nil
}

func returnStatusFilter(filter service.ShareFilter) int8 {
	if len(filter.Statuses) == 0 {
		return 0
	}

	return filter.Statuses[0]
}

// CreateShare create a share
func (ctl ShareController) CreateShare(ctx web.Context, req web.Request, shareSrv service.ShareService, pub event.Publisher) error {
	var share service.ShareUpdateFields
	if err := req.Unmarshal(&share); err != nil {
		return err
	}

	if share.ShareUser == "" {
		share.ShareUser = currentUser(req).Name
	}

	shareID, err := shareSrv.CreateShare(context.TODO(), service.Share{
		ShareUpdateFields: share,
		CreateUserId:      currentUser(req).Id,
	})
	if err == nil {
		pub.Publish(event2.ShareCreatedEvent{
			ShareID:   shareID,
			Share:     share,
			Creator:   currentUser(req),
			CreatedAt: time.Now(),
		})
	}

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

type ShareDetail struct {
	service.ShareDetail
	UserLike bool `json:"user_like"`
	UserJoin bool `json:"user_join"`
}

// Share get a share detail
func (ctl ShareController) Share(ctx web.Context, req web.Request, shareSrv service.ShareService) (*ShareDetail, error) {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return nil, err
	}

	res, err := shareSrv.GetShareByID(context.TODO(), int64(id))
	if err != nil || res == nil {
		return nil, err
	}

	detail := ShareDetail{ShareDetail: *res}

	userLikeOrJoinShares, err := shareSrv.IsUserLikeOrJoinShares(context.TODO(), currentUser(req).Id, []int64{res.Share.Id})
	if err != nil {
		log.Errorf("query user_like_or_join_shares failed: %v", err)
	}

	if s, ok := userLikeOrJoinShares[res.Share.Id]; ok {
		detail.UserLike = s.Like
		detail.UserJoin = s.Join
	}

	return &detail, nil
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

type PlanUpdateFields struct {
	ShareDate    string `json:"share_date,omitempty"`
	ShareTime    string `json:"share_time,omitempty"`
	ShareRoom    string `json:"share_room"`
	PlanDuration int64  `json:"plan_duration"`
	Note         string `json:"note"`
}

// CreateOrUpdateSharePlan create or update a share plan
func (ctl ShareController) CreateOrUpdateSharePlan(ctx web.Context, req web.Request, shareSrv service.ShareService) error {
	id, err := strconv.Atoi(req.PathVar("id"))
	if err != nil {
		return err
	}

	var plan PlanUpdateFields
	if err := req.Unmarshal(&plan); err != nil {
		return err
	}

	if plan.ShareDate == "" {
		return service.NewValidateError(errors.New("计划分享日期不能为空"))
	}

	shareAt, err := time.Parse("2006-01-02", plan.ShareDate)
	if err != nil {
		return service.NewValidateError(fmt.Errorf("计划分享日期格式不合法: %v", err))
	}

	if plan.ShareTime != "" {
		shareAt, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", shareAt.Format("2006-01-02"), plan.ShareTime))
		if err != nil {
			return service.NewValidateError(fmt.Errorf("计划分享时间格式不合法: %v", err))
		}
	}

	_, err = shareSrv.CreateOrUpdatePlan(context.TODO(), int64(id), service.PlanUpdateFields{
		ShareAt:      shareAt,
		ShareRoom:    plan.ShareRoom,
		PlanDuration: plan.PlanDuration,
		Note:         plan.Note,
	})
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
