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

// User is a User object
type User struct {
	original  *userOriginal
	userModel *UserModel

	Id        null.Int
	Uuid      null.String
	Name      null.String
	CreatedAt null.Time
	UpdatedAt null.Time
}

// As convert object to other type
// dst must be a pointer to struct
func (inst *User) As(dst interface{}) error {
	return coll.CopyProperties(inst, dst)
}

// SetModel set model for User
func (inst *User) SetModel(userModel *UserModel) {
	inst.userModel = userModel
}

// userOriginal is an object which stores original User from database
type userOriginal struct {
	Id        null.Int
	Uuid      null.String
	Name      null.String
	CreatedAt null.Time
	UpdatedAt null.Time
}

// Staled identify whether the object has been modified
func (inst *User) Staled() bool {
	if inst.original == nil {
		inst.original = &userOriginal{}
	}

	if inst.Id != inst.original.Id {
		return true
	}
	if inst.Uuid != inst.original.Uuid {
		return true
	}
	if inst.Name != inst.original.Name {
		return true
	}
	if inst.CreatedAt != inst.original.CreatedAt {
		return true
	}
	if inst.UpdatedAt != inst.original.UpdatedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (inst *User) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if inst.original == nil {
		inst.original = &userOriginal{}
	}

	if inst.Id != inst.original.Id {
		kv["id"] = inst.Id
	}
	if inst.Uuid != inst.original.Uuid {
		kv["uuid"] = inst.Uuid
	}
	if inst.Name != inst.original.Name {
		kv["name"] = inst.Name
	}
	if inst.CreatedAt != inst.original.CreatedAt {
		kv["created_at"] = inst.CreatedAt
	}
	if inst.UpdatedAt != inst.original.UpdatedAt {
		kv["updated_at"] = inst.UpdatedAt
	}

	return kv
}

// Save create a new model or update it
func (inst *User) Save() error {
	if inst.userModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := inst.userModel.SaveOrUpdate(*inst)
	if err != nil {
		return err
	}

	inst.Id = null.IntFrom(id)
	return nil
}

// Delete remove a user
func (inst *User) Delete() error {
	if inst.userModel == nil {
		return query.ErrModelNotSet
	}

	_, err := inst.userModel.DeleteById(inst.Id.Int64)
	if err != nil {
		return err
	}

	return nil
}

// String convert instance to json string
func (inst *User) String() string {
	rs, _ := json.Marshal(inst)
	return string(rs)
}

type userScope struct {
	name  string
	apply func(builder query.Condition)
}

var userGlobalScopes = make([]userScope, 0)
var userLocalScopes = make([]userScope, 0)

// AddGlobalScopeForUser assign a global scope to a model
func AddGlobalScopeForUser(name string, apply func(builder query.Condition)) {
	userGlobalScopes = append(userGlobalScopes, userScope{name: name, apply: apply})
}

// AddLocalScopeForUser assign a local scope to a model
func AddLocalScopeForUser(name string, apply func(builder query.Condition)) {
	userLocalScopes = append(userLocalScopes, userScope{name: name, apply: apply})
}

func (m *UserModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range userGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range userLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *UserModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *UserModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type UserPlain struct {
	Id        int64
	Uuid      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w UserPlain) ToUser() User {
	return User{

		Id:        null.IntFrom(int64(w.Id)),
		Uuid:      null.StringFrom(w.Uuid),
		Name:      null.StringFrom(w.Name),
		CreatedAt: null.TimeFrom(w.CreatedAt),
		UpdatedAt: null.TimeFrom(w.UpdatedAt),
	}
}

// As convert object to other type
// dst must be a pointer to struct
func (w UserPlain) As(dst interface{}) error {
	return coll.CopyProperties(w, dst)
}

func (w *User) ToUserPlain() UserPlain {
	return UserPlain{

		Id:        w.Id.Int64,
		Uuid:      w.Uuid.String,
		Name:      w.Name.String,
		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}
}

// UserModel is a model which encapsulates the operations of the object
type UserModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var userTableName = "user"

func SetUserTable(tableName string) {
	userTableName = tableName
}

// NewUserModel create a UserModel
func NewUserModel(db query.Database) *UserModel {
	return &UserModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           userTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *UserModel) GetDB() query.Database {
	return m.db.GetDB()
}

func (m *UserModel) clone() *UserModel {
	return &UserModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *UserModel) WithoutGlobalScopes(names ...string) *UserModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *UserModel) WithLocalScopes(names ...string) *UserModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Condition add query builder to model
func (m *UserModel) Condition(builder query.SQLBuilder) *UserModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *UserModel) Find(id int64) (User, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *UserModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *UserModel) Count(builders ...query.SQLBuilder) (int64, error) {
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

func (m *UserModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]User, query.PaginateMeta, error) {
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
func (m *UserModel) Get(builders ...query.SQLBuilder) ([]User, error) {
	b := m.query.Merge(builders...).Table(m.tableName).AppendCondition(m.applyScope())
	if len(b.GetFields()) == 0 {
		b = b.Select(
			"id",
			"uuid",
			"name",
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
		case "uuid":
			selectFields = append(selectFields, f)
		case "name":
			selectFields = append(selectFields, f)
		case "created_at":
			selectFields = append(selectFields, f)
		case "updated_at":
			selectFields = append(selectFields, f)
		}
	}

	var createScanVar = func(fields []query.Expr) (*User, []interface{}) {
		var userVar User
		scanFields := make([]interface{}, 0)

		for _, f := range fields {
			switch strcase.ToSnake(f.Value) {

			case "id":
				scanFields = append(scanFields, &userVar.Id)
			case "uuid":
				scanFields = append(scanFields, &userVar.Uuid)
			case "name":
				scanFields = append(scanFields, &userVar.Name)
			case "created_at":
				scanFields = append(scanFields, &userVar.CreatedAt)
			case "updated_at":
				scanFields = append(scanFields, &userVar.UpdatedAt)
			}
		}

		return &userVar, scanFields
	}

	sqlStr, params := b.Fields(selectFields...).ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		userReal, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		userReal.SetModel(m)
		users = append(users, *userReal)
	}

	return users, nil
}

// First return first result for given query
func (m *UserModel) First(builders ...query.SQLBuilder) (User, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return User{}, err
	}

	if len(res) == 0 {
		return User{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new user to database
func (m *UserModel) Create(kv query.KV) (int64, error) {

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

// SaveAll save all users to database
func (m *UserModel) SaveAll(users []User) ([]int64, error) {
	ids := make([]int64, 0)
	for _, user := range users {
		id, err := m.Save(user)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a user to database
func (m *UserModel) Save(user User) (int64, error) {
	return m.Create(user.StaledKV())
}

// SaveOrUpdate save a new user or update it when it has a id > 0
func (m *UserModel) SaveOrUpdate(user User) (id int64, updated bool, err error) {
	if user.Id.Int64 > 0 {
		_, _err := m.UpdateById(user.Id.Int64, user)
		return user.Id.Int64, true, _err
	}

	_id, _err := m.Save(user)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *UserModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
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
func (m *UserModel) Update(user User, builders ...query.SQLBuilder) (int64, error) {
	return m.UpdateFields(user.StaledKV(), builders...)
}

// UpdateById update a model by id
func (m *UserModel) UpdateById(id int64, user User) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Update(user)
}

// Delete remove a model
func (m *UserModel) Delete(builders ...query.SQLBuilder) (int64, error) {

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// DeleteById remove a model by id
func (m *UserModel) DeleteById(id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Delete()
}
