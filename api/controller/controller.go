package controller

import (
	"github.com/mylxsw/eloquent/query"
)

// PaginateRes is a struct for paginate response
type PaginateRes struct {
	Page   query.PaginateMeta     `json:"page"`
	Data   interface{}            `json:"data"`
	Search map[string]interface{} `json:"search,omitempty"`
	Extra  interface{}            `json:"extra,omitempty"`
}

type IDRes struct {
	ID int64 `json:"id"`
}
