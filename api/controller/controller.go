package controller

import (
	"strconv"

	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/glacier/web"
)

// PaginateRes is a struct for paginate response
type PaginateRes struct {
	Page query.PaginateMeta `json:"page"`
	Data interface{}        `json:"data"`
}

// currentUserID extract current user id from request
func currentUserID(req web.Request) int64 {
	userID, err := strconv.Atoi(req.Header("mock-user-id"))
	if err != nil {
		return 101
	}

	return int64(userID)
}
