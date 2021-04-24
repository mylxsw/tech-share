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

}

// ShareUserRel is a ShareUserRel object
type ShareUserRel struct {
	original          *shareUserRelOriginal
	shareUserRelModel *ShareUserRelModel

	Id        null.Int
	ShareId   null.Int
	UserId    null.Int
	RelType   null.Int
	CreatedAt null.Time
}

// As convert object to other type
// dst must be a pointer to struct
func (inst *ShareUserRel) As(dst interface{}) error {
	return coll.CopyProperties(inst, dst)
}

// SetModel set model for ShareUserRel
func (inst *ShareUserRel) SetModel(shareUserRelModel *ShareUserRelModel) {
	inst.shareUserRelModel = shareUserRelModel
}

// shareUserRelOriginal is an object which stores original ShareUserRel from database
type shareUserRelOriginal struct {
	Id        null.Int
	ShareId   null.Int
	UserId    null.Int
	RelType   null.Int
	CreatedAt null.Time
}

// Staled identify whether the object has been modified
func (inst *ShareUserRel) Staled(onlyFields ...string) bool {
	if inst.original == nil {
		inst.original = &shareUserRelOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			return true
		}
		if inst.ShareId != inst.original.ShareId {
			return true
		}
		if inst.UserId != inst.original.UserId {
			return true
		}
		if inst.RelType != inst.original.RelType {
			return true
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			return true
		}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {

			case "id":
				if inst.Id != inst.original.Id {
					return true
				}
			case "share_id":
				if inst.ShareId != inst.original.ShareId {
					return true
				}
			case "user_id":
				if inst.UserId != inst.original.UserId {
					return true
				}
			case "rel_type":
				if inst.RelType != inst.original.RelType {
					return true
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					return true
				}
			default:
			}
		}
	}

	return false
}

// StaledKV return all fields has been modified
func (inst *ShareUserRel) StaledKV(onlyFields ...string) query.KV {
	kv := make(query.KV, 0)

	if inst.original == nil {
		inst.original = &shareUserRelOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			kv["id"] = inst.Id
		}
		if inst.ShareId != inst.original.ShareId {
			kv["share_id"] = inst.ShareId
		}
		if inst.UserId != inst.original.UserId {
			kv["user_id"] = inst.UserId
		}
		if inst.RelType != inst.original.RelType {
			kv["rel_type"] = inst.RelType
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			kv["created_at"] = inst.CreatedAt
		}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {

			case "id":
				if inst.Id != inst.original.Id {
					kv["id"] = inst.Id
				}
			case "share_id":
				if inst.ShareId != inst.original.ShareId {
					kv["share_id"] = inst.ShareId
				}
			case "user_id":
				if inst.UserId != inst.original.UserId {
					kv["user_id"] = inst.UserId
				}
			case "rel_type":
				if inst.RelType != inst.original.RelType {
					kv["rel_type"] = inst.RelType
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					kv["created_at"] = inst.CreatedAt
				}
			default:
			}
		}
	}

	return kv
}

// Save create a new model or update it
func (inst *ShareUserRel) Save() error {
	if inst.shareUserRelModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := inst.shareUserRelModel.SaveOrUpdate(*inst)
	if err != nil {
		return err
	}

	inst.Id = null.IntFrom(id)
	return nil
}

// Delete remove a share_user_rel
func (inst *ShareUserRel) Delete() error {
	if inst.shareUserRelModel == nil {
		return query.ErrModelNotSet
	}

	_, err := inst.shareUserRelModel.DeleteById(inst.Id.Int64)
	if err != nil {
		return err
	}

	return nil
}

// String convert instance to json string
func (inst *ShareUserRel) String() string {
	rs, _ := json.Marshal(inst)
	return string(rs)
}

func (inst *ShareUserRel) User() *ShareUserRelBelongsToUserRel {
	return &ShareUserRelBelongsToUserRel{
		source:   inst,
		relModel: NewUserModel(inst.shareUserRelModel.GetDB()),
	}
}

type ShareUserRelBelongsToUserRel struct {
	source   *ShareUserRel
	relModel *UserModel
}

func (rel *ShareUserRelBelongsToUserRel) Create(target User) (int64, error) {
	targetId, err := rel.relModel.Save(target)
	if err != nil {
		return 0, err
	}

	target.Id = null.IntFrom(targetId)

	rel.source.UserId = target.Id
	if err := rel.source.Save(); err != nil {
		return targetId, err
	}

	return targetId, nil
}

func (rel *ShareUserRelBelongsToUserRel) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("id", rel.source.UserId).Merge(builders...)

	return rel.relModel.Exists(builder)
}

func (rel *ShareUserRelBelongsToUserRel) First(builders ...query.SQLBuilder) (User, error) {
	builder := query.Builder().Where("id", rel.source.UserId).Limit(1).Merge(builders...)

	return rel.relModel.First(builder)
}

func (rel *ShareUserRelBelongsToUserRel) Associate(target User) error {
	rel.source.UserId = target.Id
	return rel.source.Save()
}

func (rel *ShareUserRelBelongsToUserRel) Dissociate() error {
	rel.source.UserId = null.IntFrom(0)
	return rel.source.Save()
}

type shareUserRelScope struct {
	name  string
	apply func(builder query.Condition)
}

var shareUserRelGlobalScopes = make([]shareUserRelScope, 0)
var shareUserRelLocalScopes = make([]shareUserRelScope, 0)

// AddGlobalScopeForShareUserRel assign a global scope to a model
func AddGlobalScopeForShareUserRel(name string, apply func(builder query.Condition)) {
	shareUserRelGlobalScopes = append(shareUserRelGlobalScopes, shareUserRelScope{name: name, apply: apply})
}

// AddLocalScopeForShareUserRel assign a local scope to a model
func AddLocalScopeForShareUserRel(name string, apply func(builder query.Condition)) {
	shareUserRelLocalScopes = append(shareUserRelLocalScopes, shareUserRelScope{name: name, apply: apply})
}

func (m *ShareUserRelModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range shareUserRelGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range shareUserRelLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *ShareUserRelModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *ShareUserRelModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type ShareUserRelPlain struct {
	Id        int64
	ShareId   int64
	UserId    int64
	RelType   int8
	CreatedAt time.Time
}

func (w ShareUserRelPlain) ToShareUserRel(allows ...string) ShareUserRel {
	if len(allows) == 0 {
		return ShareUserRel{

			Id:        null.IntFrom(int64(w.Id)),
			ShareId:   null.IntFrom(int64(w.ShareId)),
			UserId:    null.IntFrom(int64(w.UserId)),
			RelType:   null.IntFrom(int64(w.RelType)),
			CreatedAt: null.TimeFrom(w.CreatedAt),
		}
	}

	res := ShareUserRel{}
	for _, al := range allows {
		switch strcase.ToSnake(al) {

		case "id":
			res.Id = null.IntFrom(int64(w.Id))
		case "share_id":
			res.ShareId = null.IntFrom(int64(w.ShareId))
		case "user_id":
			res.UserId = null.IntFrom(int64(w.UserId))
		case "rel_type":
			res.RelType = null.IntFrom(int64(w.RelType))
		case "created_at":
			res.CreatedAt = null.TimeFrom(w.CreatedAt)
		default:
		}
	}

	return res
}

// As convert object to other type
// dst must be a pointer to struct
func (w ShareUserRelPlain) As(dst interface{}) error {
	return coll.CopyProperties(w, dst)
}

func (w *ShareUserRel) ToShareUserRelPlain() ShareUserRelPlain {
	return ShareUserRelPlain{

		Id:        w.Id.Int64,
		ShareId:   w.ShareId.Int64,
		UserId:    w.UserId.Int64,
		RelType:   int8(w.RelType.Int64),
		CreatedAt: w.CreatedAt.Time,
	}
}

// ShareUserRelModel is a model which encapsulates the operations of the object
type ShareUserRelModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var shareUserRelTableName = "share_user_rel"

const (
	ShareUserRelFieldId        = "id"
	ShareUserRelFieldShareId   = "share_id"
	ShareUserRelFieldUserId    = "user_id"
	ShareUserRelFieldRelType   = "rel_type"
	ShareUserRelFieldCreatedAt = "created_at"
)

// ShareUserRelFields return all fields in ShareUserRel model
func ShareUserRelFields() []string {
	return []string{
		"id",
		"share_id",
		"user_id",
		"rel_type",
		"created_at",
	}
}

func SetShareUserRelTable(tableName string) {
	shareUserRelTableName = tableName
}

// NewShareUserRelModel create a ShareUserRelModel
func NewShareUserRelModel(db query.Database) *ShareUserRelModel {
	return &ShareUserRelModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           shareUserRelTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *ShareUserRelModel) GetDB() query.Database {
	return m.db.GetDB()
}

func (m *ShareUserRelModel) clone() *ShareUserRelModel {
	return &ShareUserRelModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *ShareUserRelModel) WithoutGlobalScopes(names ...string) *ShareUserRelModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *ShareUserRelModel) WithLocalScopes(names ...string) *ShareUserRelModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Condition add query builder to model
func (m *ShareUserRelModel) Condition(builder query.SQLBuilder) *ShareUserRelModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *ShareUserRelModel) Find(id int64) (ShareUserRel, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *ShareUserRelModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *ShareUserRelModel) Count(builders ...query.SQLBuilder) (int64, error) {
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

func (m *ShareUserRelModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]ShareUserRel, query.PaginateMeta, error) {
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
func (m *ShareUserRelModel) Get(builders ...query.SQLBuilder) ([]ShareUserRel, error) {
	b := m.query.Merge(builders...).Table(m.tableName).AppendCondition(m.applyScope())
	if len(b.GetFields()) == 0 {
		b = b.Select(
			"id",
			"share_id",
			"user_id",
			"rel_type",
			"created_at",
		)
	}

	fields := b.GetFields()
	selectFields := make([]query.Expr, 0)

	for _, f := range fields {
		switch strcase.ToSnake(f.Value) {

		case "id":
			selectFields = append(selectFields, f)
		case "share_id":
			selectFields = append(selectFields, f)
		case "user_id":
			selectFields = append(selectFields, f)
		case "rel_type":
			selectFields = append(selectFields, f)
		case "created_at":
			selectFields = append(selectFields, f)
		}
	}

	var createScanVar = func(fields []query.Expr) (*ShareUserRel, []interface{}) {
		var shareUserRelVar ShareUserRel
		scanFields := make([]interface{}, 0)

		for _, f := range fields {
			switch strcase.ToSnake(f.Value) {

			case "id":
				scanFields = append(scanFields, &shareUserRelVar.Id)
			case "share_id":
				scanFields = append(scanFields, &shareUserRelVar.ShareId)
			case "user_id":
				scanFields = append(scanFields, &shareUserRelVar.UserId)
			case "rel_type":
				scanFields = append(scanFields, &shareUserRelVar.RelType)
			case "created_at":
				scanFields = append(scanFields, &shareUserRelVar.CreatedAt)
			}
		}

		return &shareUserRelVar, scanFields
	}

	sqlStr, params := b.Fields(selectFields...).ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shareUserRels := make([]ShareUserRel, 0)
	for rows.Next() {
		shareUserRelReal, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		shareUserRelReal.SetModel(m)
		shareUserRels = append(shareUserRels, *shareUserRelReal)
	}

	return shareUserRels, nil
}

// First return first result for given query
func (m *ShareUserRelModel) First(builders ...query.SQLBuilder) (ShareUserRel, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return ShareUserRel{}, err
	}

	if len(res) == 0 {
		return ShareUserRel{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new share_user_rel to database
func (m *ShareUserRelModel) Create(kv query.KV) (int64, error) {

	if _, ok := kv["created_at"]; !ok {
		kv["created_at"] = time.Now()
	}

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all share_user_rels to database
func (m *ShareUserRelModel) SaveAll(shareUserRels []ShareUserRel) ([]int64, error) {
	ids := make([]int64, 0)
	for _, shareUserRel := range shareUserRels {
		id, err := m.Save(shareUserRel)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a share_user_rel to database
func (m *ShareUserRelModel) Save(shareUserRel ShareUserRel, onlyFields ...string) (int64, error) {
	return m.Create(shareUserRel.StaledKV(onlyFields...))
}

// SaveOrUpdate save a new share_user_rel or update it when it has a id > 0
func (m *ShareUserRelModel) SaveOrUpdate(shareUserRel ShareUserRel, onlyFields ...string) (id int64, updated bool, err error) {
	if shareUserRel.Id.Int64 > 0 {
		_, _err := m.UpdateById(shareUserRel.Id.Int64, shareUserRel, onlyFields...)
		return shareUserRel.Id.Int64, true, _err
	}

	_id, _err := m.Save(shareUserRel, onlyFields...)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *ShareUserRelModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

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
func (m *ShareUserRelModel) Update(shareUserRel ShareUserRel, builders ...query.SQLBuilder) (int64, error) {
	return m.UpdateFields(shareUserRel.StaledKV(), builders...)
}

// UpdatePart update a model for given query
func (m *ShareUserRelModel) UpdatePart(shareUserRel ShareUserRel, onlyFields ...string) (int64, error) {
	return m.UpdateFields(shareUserRel.StaledKV(onlyFields...))
}

// UpdateById update a model by id
func (m *ShareUserRelModel) UpdateById(id int64, shareUserRel ShareUserRel, onlyFields ...string) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).UpdateFields(shareUserRel.StaledKV(onlyFields...))
}

// Delete remove a model
func (m *ShareUserRelModel) Delete(builders ...query.SQLBuilder) (int64, error) {

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// DeleteById remove a model by id
func (m *ShareUserRelModel) DeleteById(id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Delete()
}