package controller

import (
	"github.com/mylxsw/eloquent/query"
)

// PaginateRes is a struct for paginate response
type PaginateRes struct {
	Page query.PaginateMeta `json:"page"`
	Data interface{}        `json:"data"`
}

type IDRes struct {
	ID int64 `json:"id"`
}
