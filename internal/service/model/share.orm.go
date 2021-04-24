package model

// !!! DO NOT EDIT THIS FILE

import (
	"context"
	"encoding/json"
	"github.com/iancoleman/strcase"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
	"time"
)

func init() {

	// AddShareGlobalScope assign a global scope to a model for soft delete
	AddGlobalScopeForShare("soft_delete", func(builder query.Condition) {
		builder.WhereNull("deleted_at")
	})

}

// Share is a Share object
type Share struct {
	original   *shareOriginal
	shareModel *ShareModel

	Id           null.Int
	Subject      null.String
	Description  null.String
	SubjectType  null.String
	Status       null.Int
	ShareUser    null.String
	CreateUserId null.Int
	Note         null.String
	LikeCount    null.Int
	JoinCount    null.Int
	CreatedAt    null.Time
	UpdatedAt    null.Time
	DeletedAt    null.Time
}

// As convert object to other type
// dst must be a pointer to struct
func (inst *Share) As(dst interface{}) error {
	return coll.CopyProperties(inst, dst)
}

// SetModel set model for Share
func (inst *Share) SetModel(shareModel *ShareModel) {
	inst.shareModel = shareModel
}

// shareOriginal is an object which stores original Share from database
type shareOriginal struct {
	Id           null.Int
	Subject      null.String
	Description  null.String
	SubjectType  null.String
	Status       null.Int
	ShareUser    null.String
	CreateUserId null.Int
	Note         null.String
	LikeCount    null.Int
	JoinCount    null.Int
	CreatedAt    null.Time
	UpdatedAt    null.Time
	DeletedAt    null.Time
}

// Staled identify whether the object has been modified
func (inst *Share) Staled(onlyFields ...string) bool {
	if inst.original == nil {
		inst.original = &shareOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			return true
		}
		if inst.Subject != inst.original.Subject {
			return true
		}
		if inst.Description != inst.original.Description {
			return true
		}
		if inst.SubjectType != inst.original.SubjectType {
			return true
		}
		if inst.Status != inst.original.Status {
			return true
		}
		if inst.ShareUser != inst.original.ShareUser {
			return true
		}
		if inst.CreateUserId != inst.original.CreateUserId {
			return true
		}
		if inst.Note != inst.original.Note {
			return true
		}
		if inst.LikeCount != inst.original.LikeCount {
			return true
		}
		if inst.JoinCount != inst.original.JoinCount {
			return true
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			return true
		}
		if inst.UpdatedAt != inst.original.UpdatedAt {
			return true
		}
		if inst.DeletedAt != inst.original.DeletedAt {
			return true
		}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {

			case "id":
				if inst.Id != inst.original.Id {
					return true
				}
			case "subject":
				if inst.Subject != inst.original.Subject {
					return true
				}
			case "description":
				if inst.Description != inst.original.Description {
					return true
				}
			case "subject_type":
				if inst.SubjectType != inst.original.SubjectType {
					return true
				}
			case "status":
				if inst.Status != inst.original.Status {
					return true
				}
			case "share_user":
				if inst.ShareUser != inst.original.ShareUser {
					return true
				}
			case "create_user_id":
				if inst.CreateUserId != inst.original.CreateUserId {
					return true
				}
			case "note":
				if inst.Note != inst.original.Note {
					return true
				}
			case "like_count":
				if inst.LikeCount != inst.original.LikeCount {
					return true
				}
			case "join_count":
				if inst.JoinCount != inst.original.JoinCount {
					return true
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					return true
				}
			case "updated_at":
				if inst.UpdatedAt != inst.original.UpdatedAt {
					return true
				}
			case "deleted_at":
				if inst.DeletedAt != inst.original.DeletedAt {
					return true
				}
			default:
			}
		}
	}

	return false
}

// StaledKV return all fields has been modified
func (inst *Share) StaledKV(onlyFields ...string) query.KV {
	kv := make(query.KV, 0)

	if inst.original == nil {
		inst.original = &shareOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			kv["id"] = inst.Id
		}
		if inst.Subject != inst.original.Subject {
			kv["subject"] = inst.Subject
		}
		if inst.Description != inst.original.Description {
			kv["description"] = inst.Description
		}
		if inst.SubjectType != inst.original.SubjectType {
			kv["subject_type"] = inst.SubjectType
		}
		if inst.Status != inst.original.Status {
			kv["status"] = inst.Status
		}
		if inst.ShareUser != inst.original.ShareUser {
			kv["share_user"] = inst.ShareUser
		}
		if inst.CreateUserId != inst.original.CreateUserId {
			kv["create_user_id"] = inst.CreateUserId
		}
		if inst.Note != inst.original.Note {
			kv["note"] = inst.Note
		}
		if inst.LikeCount != inst.original.LikeCount {
			kv["like_count"] = inst.LikeCount
		}
		if inst.JoinCount != inst.original.JoinCount {
			kv["join_count"] = inst.JoinCount
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			kv["created_at"] = inst.CreatedAt
		}
		if inst.UpdatedAt != inst.original.UpdatedAt {
			kv["updated_at"] = inst.UpdatedAt
		}
		if inst.DeletedAt != inst.original.DeletedAt {
			kv["deleted_at"] = inst.DeletedAt
		}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {

			case "id":
				if inst.Id != inst.original.Id {
					kv["id"] = inst.Id
				}
			case "subject":
				if inst.Subject != inst.original.Subject {
					kv["subject"] = inst.Subject
				}
			case "description":
				if inst.Description != inst.original.Description {
					kv["description"] = inst.Description
				}
			case "subject_type":
				if inst.SubjectType != inst.original.SubjectType {
					kv["subject_type"] = inst.SubjectType
				}
			case "status":
				if inst.Status != inst.original.Status {
					kv["status"] = inst.Status
				}
			case "share_user":
				if inst.ShareUser != inst.original.ShareUser {
					kv["share_user"] = inst.ShareUser
				}
			case "create_user_id":
				if inst.CreateUserId != inst.original.CreateUserId {
					kv["create_user_id"] = inst.CreateUserId
				}
			case "note":
				if inst.Note != inst.original.Note {
					kv["note"] = inst.Note
				}
			case "like_count":
				if inst.LikeCount != inst.original.LikeCount {
					kv["like_count"] = inst.LikeCount
				}
			case "join_count":
				if inst.JoinCount != inst.original.JoinCount {
					kv["join_count"] = inst.JoinCount
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					kv["created_at"] = inst.CreatedAt
				}
			case "updated_at":
				if inst.UpdatedAt != inst.original.UpdatedAt {
					kv["updated_at"] = inst.UpdatedAt
				}
			case "deleted_at":
				if inst.DeletedAt != inst.original.DeletedAt {
					kv["deleted_at"] = inst.DeletedAt
				}
			default:
			}
		}
	}

	return kv
}

// Save create a new model or update it
func (inst *Share) Save() error {
	if inst.shareModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := inst.shareModel.SaveOrUpdate(*inst)
	if err != nil {
		return err
	}

	inst.Id = null.IntFrom(id)
	return nil
}

// Delete remove a share
func (inst *Share) Delete() error {
	if inst.shareModel == nil {
		return query.ErrModelNotSet
	}

	_, err := inst.shareModel.DeleteById(inst.Id.Int64)
	if err != nil {
		return err
	}

	return nil
}

// String convert instance to json string
func (inst *Share) String() string {
	rs, _ := json.Marshal(inst)
	return string(rs)
}

func (inst *Share) SharePlan() *ShareHasOneSharePlanRel {
	return &ShareHasOneSharePlanRel{
		source:   inst,
		relModel: NewSharePlanModel(inst.shareModel.GetDB()),
	}
}

type ShareHasOneSharePlanRel struct {
	source   *Share
	relModel *SharePlanModel
}

func (rel *ShareHasOneSharePlanRel) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Exists(builder)
}

func (rel *ShareHasOneSharePlanRel) First(builders ...query.SQLBuilder) (SharePlan, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Limit(1).Merge(builders...)
	return rel.relModel.First(builder)
}

func (rel *ShareHasOneSharePlanRel) Create(target SharePlan) (int64, error) {
	target.ShareId = rel.source.Id
	return rel.relModel.Save(target)
}

func (rel *ShareHasOneSharePlanRel) Associate(target SharePlan) error {
	_, err := rel.relModel.UpdateFields(
		query.KV{"share_id": rel.source.Id},
		query.Builder().Where("id", target.Id),
	)
	return err
}

func (rel *ShareHasOneSharePlanRel) Dissociate() error {
	_, err := rel.relModel.UpdateFields(
		query.KV{"share_id": nil},
		query.Builder().Where("share_id", rel.source.Id),
	)

	return err
}

func (inst *Share) Attachments() *ShareHasManyAttachmentRel {
	return &ShareHasManyAttachmentRel{
		source:   inst,
		relModel: NewAttachmentModel(inst.shareModel.GetDB()),
	}
}

type ShareHasManyAttachmentRel struct {
	source   *Share
	relModel *AttachmentModel
}

func (rel *ShareHasManyAttachmentRel) Get(builders ...query.SQLBuilder) ([]Attachment, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Get(builder)
}

func (rel *ShareHasManyAttachmentRel) Count(builders ...query.SQLBuilder) (int64, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Count(builder)
}

func (rel *ShareHasManyAttachmentRel) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Exists(builder)
}

func (rel *ShareHasManyAttachmentRel) First(builders ...query.SQLBuilder) (Attachment, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Limit(1).Merge(builders...)
	return rel.relModel.First(builder)
}

func (rel *ShareHasManyAttachmentRel) Create(target Attachment) (int64, error) {
	target.ShareId = rel.source.Id
	return rel.relModel.Save(target)
}

func (inst *Share) ShareUserRels() *ShareHasManyShareUserRelRel {
	return &ShareHasManyShareUserRelRel{
		source:   inst,
		relModel: NewShareUserRelModel(inst.shareModel.GetDB()),
	}
}

type ShareHasManyShareUserRelRel struct {
	source   *Share
	relModel *ShareUserRelModel
}

func (rel *ShareHasManyShareUserRelRel) Get(builders ...query.SQLBuilder) ([]ShareUserRel, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Get(builder)
}

func (rel *ShareHasManyShareUserRelRel) Count(builders ...query.SQLBuilder) (int64, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Count(builder)
}

func (rel *ShareHasManyShareUserRelRel) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Exists(builder)
}

func (rel *ShareHasManyShareUserRelRel) First(builders ...query.SQLBuilder) (ShareUserRel, error) {
	builder := query.Builder().Where("share_id", rel.source.Id).Limit(1).Merge(builders...)
	return rel.relModel.First(builder)
}

func (rel *ShareHasManyShareUserRelRel) Create(target ShareUserRel) (int64, error) {
	target.ShareId = rel.source.Id
	return rel.relModel.Save(target)
}

type shareScope struct {
	name  string
	apply func(builder query.Condition)
}

var shareGlobalScopes = make([]shareScope, 0)
var shareLocalScopes = make([]shareScope, 0)

// AddGlobalScopeForShare assign a global scope to a model
func AddGlobalScopeForShare(name string, apply func(builder query.Condition)) {
	shareGlobalScopes = append(shareGlobalScopes, shareScope{name: name, apply: apply})
}

// AddLocalScopeForShare assign a local scope to a model
func AddLocalScopeForShare(name string, apply func(builder query.Condition)) {
	shareLocalScopes = append(shareLocalScopes, shareScope{name: name, apply: apply})
}

func (m *ShareModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range shareGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range shareLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *ShareModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *ShareModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type SharePlain struct {
	Id           int64
	Subject      string
	Description  string
	SubjectType  string
	Status       int8
	ShareUser    string
	CreateUserId int64
	Note         string
	LikeCount    int64
	JoinCount    int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

func (w SharePlain) ToShare(allows ...string) Share {
	if len(allows) == 0 {
		return Share{

			Id:           null.IntFrom(int64(w.Id)),
			Subject:      null.StringFrom(w.Subject),
			Description:  null.StringFrom(w.Description),
			SubjectType:  null.StringFrom(w.SubjectType),
			Status:       null.IntFrom(int64(w.Status)),
			ShareUser:    null.StringFrom(w.ShareUser),
			CreateUserId: null.IntFrom(int64(w.CreateUserId)),
			Note:         null.StringFrom(w.Note),
			LikeCount:    null.IntFrom(int64(w.LikeCount)),
			JoinCount:    null.IntFrom(int64(w.JoinCount)),
			CreatedAt:    null.TimeFrom(w.CreatedAt),
			UpdatedAt:    null.TimeFrom(w.UpdatedAt),
			DeletedAt:    null.TimeFrom(w.DeletedAt),
		}
	}

	res := Share{}
	for _, al := range allows {
		switch strcase.ToSnake(al) {

		case "id":
			res.Id = null.IntFrom(int64(w.Id))
		case "subject":
			res.Subject = null.StringFrom(w.Subject)
		case "description":
			res.Description = null.StringFrom(w.Description)
		case "subject_type":
			res.SubjectType = null.StringFrom(w.SubjectType)
		case "status":
			res.Status = null.IntFrom(int64(w.Status))
		case "share_user":
			res.ShareUser = null.StringFrom(w.ShareUser)
		case "create_user_id":
			res.CreateUserId = null.IntFrom(int64(w.CreateUserId))
		case "note":
			res.Note = null.StringFrom(w.Note)
		case "like_count":
			res.LikeCount = null.IntFrom(int64(w.LikeCount))
		case "join_count":
			res.JoinCount = null.IntFrom(int64(w.JoinCount))
		case "created_at":
			res.CreatedAt = null.TimeFrom(w.CreatedAt)
		case "updated_at":
			res.UpdatedAt = null.TimeFrom(w.UpdatedAt)
		case "deleted_at":
			res.DeletedAt = null.TimeFrom(w.DeletedAt)
		default:
		}
	}

	return res
}

// As convert object to other type
// dst must be a pointer to struct
func (w SharePlain) As(dst interface{}) error {
	return coll.CopyProperties(w, dst)
}

func (w *Share) ToSharePlain() SharePlain {
	return SharePlain{

		Id:           w.Id.Int64,
		Subject:      w.Subject.String,
		Description:  w.Description.String,
		SubjectType:  w.SubjectType.String,
		Status:       int8(w.Status.Int64),
		ShareUser:    w.ShareUser.String,
		CreateUserId: w.CreateUserId.Int64,
		Note:         w.Note.String,
		LikeCount:    w.LikeCount.Int64,
		JoinCount:    w.JoinCount.Int64,
		CreatedAt:    w.CreatedAt.Time,
		UpdatedAt:    w.UpdatedAt.Time,
		DeletedAt:    w.DeletedAt.Time,
	}
}

// ShareModel is a model which encapsulates the operations of the object
type ShareModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var shareTableName = "share"

const (
	ShareFieldId           = "id"
	ShareFieldSubject      = "subject"
	ShareFieldDescription  = "description"
	ShareFieldSubjectType  = "subject_type"
	ShareFieldStatus       = "status"
	ShareFieldShareUser    = "share_user"
	ShareFieldCreateUserId = "create_user_id"
	ShareFieldNote         = "note"
	ShareFieldLikeCount    = "like_count"
	ShareFieldJoinCount    = "join_count"
	ShareFieldCreatedAt    = "created_at"
	ShareFieldUpdatedAt    = "updated_at"
	ShareFieldDeletedAt    = "deleted_at"
)

// ShareFields return all fields in Share model
func ShareFields() []string {
	return []string{
		"id",
		"subject",
		"description",
		"subject_type",
		"status",
		"share_user",
		"create_user_id",
		"note",
		"like_count",
		"join_count",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func SetShareTable(tableName string) {
	shareTableName = tableName
}

// NewShareModel create a ShareModel
func NewShareModel(db query.Database) *ShareModel {
	return &ShareModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           shareTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *ShareModel) GetDB() query.Database {
	return m.db.GetDB()
}

// WithTrashed force soft deleted models to appear in a result set
func (m *ShareModel) WithTrashed() *ShareModel {
	return m.WithoutGlobalScopes("soft_delete")
}

func (m *ShareModel) clone() *ShareModel {
	return &ShareModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *ShareModel) WithoutGlobalScopes(names ...string) *ShareModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *ShareModel) WithLocalScopes(names ...string) *ShareModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Condition add query builder to model
func (m *ShareModel) Condition(builder query.SQLBuilder) *ShareModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *ShareModel) Find(id int64) (Share, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *ShareModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *ShareModel) Count(builders ...query.SQLBuilder) (int64, error) {
	sqlStr, params := m.query.
		Merge(builders...).
		Table(m.tableName).
		AppendCondition(m.applyScope()).
		ResolveCount()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	rows.Next()
	var res int64
	if err := rows.Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}

func (m *ShareModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]Share, query.PaginateMeta, error) {
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 15
	}

	meta := query.PaginateMeta{
		PerPage: perPage,
		Page:    page,
	}

	count, err := m.Count(builders...)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.LastPage = count / perPage
	if count%perPage != 0 {
		meta.LastPage += 1
	}

	res, err := m.Get(append([]query.SQLBuilder{query.Builder().Limit(perPage).Offset((page - 1) * perPage)}, builders...)...)
	if err != nil {
		return res, meta, err
	}

	return res, meta, nil
}

// Get retrieve all results for given query
func (m *ShareModel) Get(builders ...query.SQLBuilder) ([]Share, error) {
	b := m.query.Merge(builders...).Table(m.tableName).AppendCondition(m.applyScope())
	if len(b.GetFields()) == 0 {
		b = b.Select(
			"id",
			"subject",
			"description",
			"subject_type",
			"status",
			"share_user",
			"create_user_id",
			"note",
			"like_count",
			"join_count",
			"created_at",
			"updated_at",
			"deleted_at",
		)
	}

	fields := b.GetFields()
	selectFields := make([]query.Expr, 0)

	for _, f := range fields {
		switch strcase.ToSnake(f.Value) {

		case "id":
			selectFields = append(selectFields, f)
		case "subject":
			selectFields = append(selectFields, f)
		case "description":
			selectFields = append(selectFields, f)
		case "subject_type":
			selectFields = append(selectFields, f)
		case "status":
			selectFields = append(selectFields, f)
		case "share_user":
			selectFields = append(selectFields, f)
		case "create_user_id":
			selectFields = append(selectFields, f)
		case "note":
			selectFields = append(selectFields, f)
		case "like_count":
			selectFields = append(selectFields, f)
		case "join_count":
			selectFields = append(selectFields, f)
		case "created_at":
			selectFields = append(selectFields, f)
		case "updated_at":
			selectFields = append(selectFields, f)
		case "deleted_at":
			selectFields = append(selectFields, f)
		}
	}

	var createScanVar = func(fields []query.Expr) (*Share, []interface{}) {
		var shareVar Share
		scanFields := make([]interface{}, 0)

		for _, f := range fields {
			switch strcase.ToSnake(f.Value) {

			case "id":
				scanFields = append(scanFields, &shareVar.Id)
			case "subject":
				scanFields = append(scanFields, &shareVar.Subject)
			case "description":
				scanFields = append(scanFields, &shareVar.Description)
			case "subject_type":
				scanFields = append(scanFields, &shareVar.SubjectType)
			case "status":
				scanFields = append(scanFields, &shareVar.Status)
			case "share_user":
				scanFields = append(scanFields, &shareVar.ShareUser)
			case "create_user_id":
				scanFields = append(scanFields, &shareVar.CreateUserId)
			case "note":
				scanFields = append(scanFields, &shareVar.Note)
			case "like_count":
				scanFields = append(scanFields, &shareVar.LikeCount)
			case "join_count":
				scanFields = append(scanFields, &shareVar.JoinCount)
			case "created_at":
				scanFields = append(scanFields, &shareVar.CreatedAt)
			case "updated_at":
				scanFields = append(scanFields, &shareVar.UpdatedAt)
			case "deleted_at":
				scanFields = append(scanFields, &shareVar.DeletedAt)
			}
		}

		return &shareVar, scanFields
	}

	sqlStr, params := b.Fields(selectFields...).ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shares := make([]Share, 0)
	for rows.Next() {
		shareReal, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		shareReal.SetModel(m)
		shares = append(shares, *shareReal)
	}

	return shares, nil
}

// First return first result for given query
func (m *ShareModel) First(builders ...query.SQLBuilder) (Share, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return Share{}, err
	}

	if len(res) == 0 {
		return Share{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new share to database
func (m *ShareModel) Create(kv query.KV) (int64, error) {

	if _, ok := kv["created_at"]; !ok {
		kv["created_at"] = time.Now()
	}

	if _, ok := kv["updated_at"]; !ok {
		kv["updated_at"] = time.Now()
	}

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all shares to database
func (m *ShareModel) SaveAll(shares []Share) ([]int64, error) {
	ids := make([]int64, 0)
	for _, share := range shares {
		id, err := m.Save(share)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a share to database
func (m *ShareModel) Save(share Share, onlyFields ...string) (int64, error) {
	return m.Create(share.StaledKV(onlyFields...))
}

// SaveOrUpdate save a new share or update it when it has a id > 0
func (m *ShareModel) SaveOrUpdate(share Share, onlyFields ...string) (id int64, updated bool, err error) {
	if share.Id.Int64 > 0 {
		_, _err := m.UpdateById(share.Id.Int64, share, onlyFields...)
		return share.Id.Int64, true, _err
	}

	_id, _err := m.Save(share, onlyFields...)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *ShareModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	kv["updated_at"] = time.Now()

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).
		Table(m.tableName).
		ResolveUpdate(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update update a model for given query
func (m *ShareModel) Update(share Share, builders ...query.SQLBuilder) (int64, error) {
	return m.UpdateFields(share.StaledKV(), builders...)
}

// UpdatePart update a model for given query
func (m *ShareModel) UpdatePart(share Share, onlyFields ...string) (int64, error) {
	return m.UpdateFields(share.StaledKV(onlyFields...))
}

// UpdateById update a model by id
func (m *ShareModel) UpdateById(id int64, share Share, onlyFields ...string) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).UpdateFields(share.StaledKV(onlyFields...))
}

// ForceDelete permanently remove a soft deleted model from the database
func (m *ShareModel) ForceDelete(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := m2.query.Merge(builders...).AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *ShareModel) ForceDeleteById(id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).ForceDelete()
}

// Restore restore a soft deleted model into an active state
func (m *ShareModel) Restore(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(query.KV{
		"deleted_at": nil,
	}, builders...)
}

// RestoreById restore a soft deleted model into an active state by id
func (m *ShareModel) RestoreById(id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Restore()
}

// Delete remove a model
func (m *ShareModel) Delete(builders ...query.SQLBuilder) (int64, error) {

	return m.UpdateFields(query.KV{
		"deleted_at": time.Now(),
	}, builders...)

}

// DeleteById remove a model by id
func (m *ShareModel) DeleteById(id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Delete()
}
