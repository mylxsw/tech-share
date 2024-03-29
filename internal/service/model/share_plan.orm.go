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

// SharePlan is a SharePlan object
type SharePlan struct {
	original       *sharePlanOriginal
	sharePlanModel *SharePlanModel

	Id           null.Int
	ShareId      null.Int
	ShareAt      null.Time
	ShareRoom    null.String
	PlanDuration null.Int
	RealDuration null.Int
	Note         null.String
	CreatedAt    null.Time
	UpdatedAt    null.Time
}

// As convert object to other type
// dst must be a pointer to struct
func (inst *SharePlan) As(dst interface{}) error {
	return coll.CopyProperties(inst, dst)
}

// SetModel set model for SharePlan
func (inst *SharePlan) SetModel(sharePlanModel *SharePlanModel) {
	inst.sharePlanModel = sharePlanModel
}

// sharePlanOriginal is an object which stores original SharePlan from database
type sharePlanOriginal struct {
	Id           null.Int
	ShareId      null.Int
	ShareAt      null.Time
	ShareRoom    null.String
	PlanDuration null.Int
	RealDuration null.Int
	Note         null.String
	CreatedAt    null.Time
	UpdatedAt    null.Time
}

// Staled identify whether the object has been modified
func (inst *SharePlan) Staled(onlyFields ...string) bool {
	if inst.original == nil {
		inst.original = &sharePlanOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			return true
		}
		if inst.ShareId != inst.original.ShareId {
			return true
		}
		if inst.ShareAt != inst.original.ShareAt {
			return true
		}
		if inst.ShareRoom != inst.original.ShareRoom {
			return true
		}
		if inst.PlanDuration != inst.original.PlanDuration {
			return true
		}
		if inst.RealDuration != inst.original.RealDuration {
			return true
		}
		if inst.Note != inst.original.Note {
			return true
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			return true
		}
		if inst.UpdatedAt != inst.original.UpdatedAt {
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
			case "share_at":
				if inst.ShareAt != inst.original.ShareAt {
					return true
				}
			case "share_room":
				if inst.ShareRoom != inst.original.ShareRoom {
					return true
				}
			case "plan_duration":
				if inst.PlanDuration != inst.original.PlanDuration {
					return true
				}
			case "real_duration":
				if inst.RealDuration != inst.original.RealDuration {
					return true
				}
			case "note":
				if inst.Note != inst.original.Note {
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
			default:
			}
		}
	}

	return false
}

// StaledKV return all fields has been modified
func (inst *SharePlan) StaledKV(onlyFields ...string) query.KV {
	kv := make(query.KV, 0)

	if inst.original == nil {
		inst.original = &sharePlanOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			kv["id"] = inst.Id
		}
		if inst.ShareId != inst.original.ShareId {
			kv["share_id"] = inst.ShareId
		}
		if inst.ShareAt != inst.original.ShareAt {
			kv["share_at"] = inst.ShareAt
		}
		if inst.ShareRoom != inst.original.ShareRoom {
			kv["share_room"] = inst.ShareRoom
		}
		if inst.PlanDuration != inst.original.PlanDuration {
			kv["plan_duration"] = inst.PlanDuration
		}
		if inst.RealDuration != inst.original.RealDuration {
			kv["real_duration"] = inst.RealDuration
		}
		if inst.Note != inst.original.Note {
			kv["note"] = inst.Note
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			kv["created_at"] = inst.CreatedAt
		}
		if inst.UpdatedAt != inst.original.UpdatedAt {
			kv["updated_at"] = inst.UpdatedAt
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
			case "share_at":
				if inst.ShareAt != inst.original.ShareAt {
					kv["share_at"] = inst.ShareAt
				}
			case "share_room":
				if inst.ShareRoom != inst.original.ShareRoom {
					kv["share_room"] = inst.ShareRoom
				}
			case "plan_duration":
				if inst.PlanDuration != inst.original.PlanDuration {
					kv["plan_duration"] = inst.PlanDuration
				}
			case "real_duration":
				if inst.RealDuration != inst.original.RealDuration {
					kv["real_duration"] = inst.RealDuration
				}
			case "note":
				if inst.Note != inst.original.Note {
					kv["note"] = inst.Note
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					kv["created_at"] = inst.CreatedAt
				}
			case "updated_at":
				if inst.UpdatedAt != inst.original.UpdatedAt {
					kv["updated_at"] = inst.UpdatedAt
				}
			default:
			}
		}
	}

	return kv
}

// Save create a new model or update it
func (inst *SharePlan) Save(onlyFields ...string) error {
	if inst.sharePlanModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := inst.sharePlanModel.SaveOrUpdate(*inst, onlyFields...)
	if err != nil {
		return err
	}

	inst.Id = null.IntFrom(id)
	return nil
}

// Delete remove a share_plan
func (inst *SharePlan) Delete() error {
	if inst.sharePlanModel == nil {
		return query.ErrModelNotSet
	}

	_, err := inst.sharePlanModel.DeleteById(inst.Id.Int64)
	if err != nil {
		return err
	}

	return nil
}

// String convert instance to json string
func (inst *SharePlan) String() string {
	rs, _ := json.Marshal(inst)
	return string(rs)
}

func (inst *SharePlan) Share() *SharePlanBelongsToShareRel {
	return &SharePlanBelongsToShareRel{
		source:   inst,
		relModel: NewShareModel(inst.sharePlanModel.GetDB()),
	}
}

type SharePlanBelongsToShareRel struct {
	source   *SharePlan
	relModel *ShareModel
}

func (rel *SharePlanBelongsToShareRel) Create(target Share) (int64, error) {
	targetId, err := rel.relModel.Save(target)
	if err != nil {
		return 0, err
	}

	target.Id = null.IntFrom(targetId)

	rel.source.ShareId = target.Id
	if err := rel.source.Save(); err != nil {
		return targetId, err
	}

	return targetId, nil
}

func (rel *SharePlanBelongsToShareRel) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("id", rel.source.ShareId).Merge(builders...)

	return rel.relModel.Exists(builder)
}

func (rel *SharePlanBelongsToShareRel) First(builders ...query.SQLBuilder) (Share, error) {
	builder := query.Builder().Where("id", rel.source.ShareId).Limit(1).Merge(builders...)

	return rel.relModel.First(builder)
}

func (rel *SharePlanBelongsToShareRel) Associate(target Share) error {
	rel.source.ShareId = target.Id
	return rel.source.Save()
}

func (rel *SharePlanBelongsToShareRel) Dissociate() error {
	rel.source.ShareId = null.IntFrom(0)
	return rel.source.Save()
}

type sharePlanScope struct {
	name  string
	apply func(builder query.Condition)
}

var sharePlanGlobalScopes = make([]sharePlanScope, 0)
var sharePlanLocalScopes = make([]sharePlanScope, 0)

// AddGlobalScopeForSharePlan assign a global scope to a model
func AddGlobalScopeForSharePlan(name string, apply func(builder query.Condition)) {
	sharePlanGlobalScopes = append(sharePlanGlobalScopes, sharePlanScope{name: name, apply: apply})
}

// AddLocalScopeForSharePlan assign a local scope to a model
func AddLocalScopeForSharePlan(name string, apply func(builder query.Condition)) {
	sharePlanLocalScopes = append(sharePlanLocalScopes, sharePlanScope{name: name, apply: apply})
}

func (m *SharePlanModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range sharePlanGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range sharePlanLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *SharePlanModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *SharePlanModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type SharePlanPlain struct {
	Id           int64
	ShareId      int64
	ShareAt      time.Time
	ShareRoom    string
	PlanDuration int
	RealDuration int
	Note         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (w SharePlanPlain) ToSharePlan(allows ...string) SharePlan {
	if len(allows) == 0 {
		return SharePlan{

			Id:           null.IntFrom(int64(w.Id)),
			ShareId:      null.IntFrom(int64(w.ShareId)),
			ShareAt:      null.TimeFrom(w.ShareAt),
			ShareRoom:    null.StringFrom(w.ShareRoom),
			PlanDuration: null.IntFrom(int64(w.PlanDuration)),
			RealDuration: null.IntFrom(int64(w.RealDuration)),
			Note:         null.StringFrom(w.Note),
			CreatedAt:    null.TimeFrom(w.CreatedAt),
			UpdatedAt:    null.TimeFrom(w.UpdatedAt),
		}
	}

	res := SharePlan{}
	for _, al := range allows {
		switch strcase.ToSnake(al) {

		case "id":
			res.Id = null.IntFrom(int64(w.Id))
		case "share_id":
			res.ShareId = null.IntFrom(int64(w.ShareId))
		case "share_at":
			res.ShareAt = null.TimeFrom(w.ShareAt)
		case "share_room":
			res.ShareRoom = null.StringFrom(w.ShareRoom)
		case "plan_duration":
			res.PlanDuration = null.IntFrom(int64(w.PlanDuration))
		case "real_duration":
			res.RealDuration = null.IntFrom(int64(w.RealDuration))
		case "note":
			res.Note = null.StringFrom(w.Note)
		case "created_at":
			res.CreatedAt = null.TimeFrom(w.CreatedAt)
		case "updated_at":
			res.UpdatedAt = null.TimeFrom(w.UpdatedAt)
		default:
		}
	}

	return res
}

// As convert object to other type
// dst must be a pointer to struct
func (w SharePlanPlain) As(dst interface{}) error {
	return coll.CopyProperties(w, dst)
}

func (w *SharePlan) ToSharePlanPlain() SharePlanPlain {
	return SharePlanPlain{

		Id:           w.Id.Int64,
		ShareId:      w.ShareId.Int64,
		ShareAt:      w.ShareAt.Time,
		ShareRoom:    w.ShareRoom.String,
		PlanDuration: int(w.PlanDuration.Int64),
		RealDuration: int(w.RealDuration.Int64),
		Note:         w.Note.String,
		CreatedAt:    w.CreatedAt.Time,
		UpdatedAt:    w.UpdatedAt.Time,
	}
}

// SharePlanModel is a model which encapsulates the operations of the object
type SharePlanModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var sharePlanTableName = "share_plan"

const (
	SharePlanFieldId           = "id"
	SharePlanFieldShareId      = "share_id"
	SharePlanFieldShareAt      = "share_at"
	SharePlanFieldShareRoom    = "share_room"
	SharePlanFieldPlanDuration = "plan_duration"
	SharePlanFieldRealDuration = "real_duration"
	SharePlanFieldNote         = "note"
	SharePlanFieldCreatedAt    = "created_at"
	SharePlanFieldUpdatedAt    = "updated_at"
)

// SharePlanFields return all fields in SharePlan model
func SharePlanFields() []string {
	return []string{
		"id",
		"share_id",
		"share_at",
		"share_room",
		"plan_duration",
		"real_duration",
		"note",
		"created_at",
		"updated_at",
	}
}

func SetSharePlanTable(tableName string) {
	sharePlanTableName = tableName
}

// NewSharePlanModel create a SharePlanModel
func NewSharePlanModel(db query.Database) *SharePlanModel {
	return &SharePlanModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           sharePlanTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *SharePlanModel) GetDB() query.Database {
	return m.db.GetDB()
}

func (m *SharePlanModel) clone() *SharePlanModel {
	return &SharePlanModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *SharePlanModel) WithoutGlobalScopes(names ...string) *SharePlanModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *SharePlanModel) WithLocalScopes(names ...string) *SharePlanModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Condition add query builder to model
func (m *SharePlanModel) Condition(builder query.SQLBuilder) *SharePlanModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *SharePlanModel) Find(id int64) (SharePlan, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *SharePlanModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *SharePlanModel) Count(builders ...query.SQLBuilder) (int64, error) {
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

func (m *SharePlanModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]SharePlan, query.PaginateMeta, error) {
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
func (m *SharePlanModel) Get(builders ...query.SQLBuilder) ([]SharePlan, error) {
	b := m.query.Merge(builders...).Table(m.tableName).AppendCondition(m.applyScope())
	if len(b.GetFields()) == 0 {
		b = b.Select(
			"id",
			"share_id",
			"share_at",
			"share_room",
			"plan_duration",
			"real_duration",
			"note",
			"created_at",
			"updated_at",
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
		case "share_at":
			selectFields = append(selectFields, f)
		case "share_room":
			selectFields = append(selectFields, f)
		case "plan_duration":
			selectFields = append(selectFields, f)
		case "real_duration":
			selectFields = append(selectFields, f)
		case "note":
			selectFields = append(selectFields, f)
		case "created_at":
			selectFields = append(selectFields, f)
		case "updated_at":
			selectFields = append(selectFields, f)
		}
	}

	var createScanVar = func(fields []query.Expr) (*SharePlan, []interface{}) {
		var sharePlanVar SharePlan
		scanFields := make([]interface{}, 0)

		for _, f := range fields {
			switch strcase.ToSnake(f.Value) {

			case "id":
				scanFields = append(scanFields, &sharePlanVar.Id)
			case "share_id":
				scanFields = append(scanFields, &sharePlanVar.ShareId)
			case "share_at":
				scanFields = append(scanFields, &sharePlanVar.ShareAt)
			case "share_room":
				scanFields = append(scanFields, &sharePlanVar.ShareRoom)
			case "plan_duration":
				scanFields = append(scanFields, &sharePlanVar.PlanDuration)
			case "real_duration":
				scanFields = append(scanFields, &sharePlanVar.RealDuration)
			case "note":
				scanFields = append(scanFields, &sharePlanVar.Note)
			case "created_at":
				scanFields = append(scanFields, &sharePlanVar.CreatedAt)
			case "updated_at":
				scanFields = append(scanFields, &sharePlanVar.UpdatedAt)
			}
		}

		return &sharePlanVar, scanFields
	}

	sqlStr, params := b.Fields(selectFields...).ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sharePlans := make([]SharePlan, 0)
	for rows.Next() {
		sharePlanReal, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		sharePlanReal.original = &sharePlanOriginal{}
		_ = coll.CopyProperties(sharePlanReal, sharePlanReal.original)

		sharePlanReal.SetModel(m)
		sharePlans = append(sharePlans, *sharePlanReal)
	}

	return sharePlans, nil
}

// First return first result for given query
func (m *SharePlanModel) First(builders ...query.SQLBuilder) (SharePlan, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return SharePlan{}, err
	}

	if len(res) == 0 {
		return SharePlan{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new share_plan to database
func (m *SharePlanModel) Create(kv query.KV) (int64, error) {

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

// SaveAll save all share_plans to database
func (m *SharePlanModel) SaveAll(sharePlans []SharePlan) ([]int64, error) {
	ids := make([]int64, 0)
	for _, sharePlan := range sharePlans {
		id, err := m.Save(sharePlan)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a share_plan to database
func (m *SharePlanModel) Save(sharePlan SharePlan, onlyFields ...string) (int64, error) {
	return m.Create(sharePlan.StaledKV(onlyFields...))
}

// SaveOrUpdate save a new share_plan or update it when it has a id > 0
func (m *SharePlanModel) SaveOrUpdate(sharePlan SharePlan, onlyFields ...string) (id int64, updated bool, err error) {
	if sharePlan.Id.Int64 > 0 {
		_, _err := m.UpdateById(sharePlan.Id.Int64, sharePlan, onlyFields...)
		return sharePlan.Id.Int64, true, _err
	}

	_id, _err := m.Save(sharePlan, onlyFields...)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *SharePlanModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
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
func (m *SharePlanModel) Update(sharePlan SharePlan, builders ...query.SQLBuilder) (int64, error) {
	return m.UpdateFields(sharePlan.StaledKV(), builders...)
}

// UpdatePart update a model for given query
func (m *SharePlanModel) UpdatePart(sharePlan SharePlan, onlyFields ...string) (int64, error) {
	return m.UpdateFields(sharePlan.StaledKV(onlyFields...))
}

// UpdateById update a model by id
func (m *SharePlanModel) UpdateById(id int64, sharePlan SharePlan, onlyFields ...string) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).UpdateFields(sharePlan.StaledKV(onlyFields...))
}

// Delete remove a model
func (m *SharePlanModel) Delete(builders ...query.SQLBuilder) (int64, error) {

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// DeleteById remove a model by id
func (m *SharePlanModel) DeleteById(id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Delete()
}
