package event

import (
	"time"

	"github.com/mylxsw/tech-share/internal/service"
)

type ShareCreatedEvent struct {
	ShareID   int64
	Share     service.ShareUpdateFields
	Creator   service.UserInfo
	CreatedAt time.Time
}
