package service

import "time"

type ValidateError struct {
	err error
}

func NewValidateError(err error) *ValidateError {
	return &ValidateError{err: err}
}

func (v *ValidateError) Error() string {
	return v.err.Error()
}

const (
	ShareStatusVoting   int8 = 1
	ShareStatusPlaned   int8 = 2
	ShareStatusFinished int8 = 3

	RelTypeLike int8 = 1
	RelTypeJoin int8 = 2
)

type ShareDetail struct {
	Share       Share        `json:"share"`
	Plan        *Plan        `json:"plan"`
	Attachments []Attachment `json:"attachments"`
	LikeUsers   []User       `json:"like_users"`
	JoinUsers   []User       `json:"join_users"`
}

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	relType int8
}

type Attachment struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	AttaType  string    `json:"atta_type"`
	AttaPath  string    `json:"atta_path"`
	CreatedAt time.Time `json:"created_at"`
}

type Plan struct {
	PlanUpdateFields
	Id           int64     `json:"id"`
	RealDuration int64     `json:"real_duration" validate:"gte=0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PlanUpdateFields struct {
	ShareAt      time.Time `json:"share_at" validate:"required"`
	ShareRoom    string    `json:"share_room"`
	PlanDuration int64     `json:"plan_duration" validate:"required,gt=0"`
	Note         string    `json:"note"`
}

type Share struct {
	ShareUpdateFields
	Id           int64     `json:"id"`
	Status       int8      `json:"status" validate:"oneof=1 2 3"`
	Note         string    `json:"note"`
	LikeCount    int64     `json:"like_count" validate:"gte=0"`
	JoinCount    int64     `json:"join_count" validate:"gte=0"`
	Attachments  string    `json:"attachments"`
	CreateUserId int64     `json:"create_user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ShareExt struct {
	Share
	ShareAt      time.Time `json:"share_at"`
	ShareRoom    string    `json:"share_room"`
	PlanDuration int64     `json:"plan_duration"`
}

type ShareUpdateFields struct {
	Subject     string `json:"subject" validate:"required,gte=2"`
	SubjectType string `json:"subject_type" validate:"required"`
	Description string `json:"description"`
	ShareUser   string `json:"share_user"`
}

type ShareFinishFields struct {
	RealDuration int64   `json:"real_duration" validate:"required,gte=0"`
	Attachments  []int64 `json:"attachments"`
	Note         string  `json:"note"`
}

type ShareFilter struct {
	Statuses []int8 `json:"statuses"`
	Creator  int64  `json:"creator"`
	Type     string `json:"type"`
}

type UserLikeOrJoinShare struct {
	ShareID int64 `json:"share_id"`
	Like    bool  `json:"like"`
	Join    bool  `json:"join"`
}
