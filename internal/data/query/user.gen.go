// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"

	"gorm.io/plugin/dbresolver"

	"auth/internal/data/model"
)

func newUser(db *gorm.DB, opts ...gen.DOOption) user {
	_user := user{}

	_user.userDo.UseDB(db, opts...)
	_user.userDo.UseModel(&model.User{})

	tableName := _user.userDo.TableName()
	_user.ALL = field.NewAsterisk(tableName)
	_user.ID = field.NewUint64(tableName, "id")
	_user.CreatedAt = field.NewField(tableName, "created_at")
	_user.UpdatedAt = field.NewField(tableName, "updated_at")
	_user.RoleID = field.NewUint64(tableName, "role_id")
	_user.Action = field.NewString(tableName, "action")
	_user.Username = field.NewString(tableName, "username")
	_user.Code = field.NewString(tableName, "code")
	_user.Password = field.NewString(tableName, "password")
	_user.Platform = field.NewString(tableName, "platform")
	_user.LastLogin = field.NewField(tableName, "last_login")
	_user.Locked = field.NewBool(tableName, "locked")
	_user.LockExpire = field.NewUint64(tableName, "lock_expire")
	_user.Wrong = field.NewUint64(tableName, "wrong")
	_user.Role = userHasOneRole{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Role", "model.Role"),
	}

	_user.fillFieldMap()

	return _user
}

type user struct {
	userDo userDo

	ALL        field.Asterisk
	ID         field.Uint64 // auto increment id
	CreatedAt  field.Field  // create time
	UpdatedAt  field.Field  // update time
	RoleID     field.Uint64 // role id
	Action     field.String // user action code array
	Username   field.String // user login name
	Code       field.String // user code
	Password   field.String // password
	Platform   field.String // device platform: pc/android/ios/mini...
	LastLogin  field.Field  // last login time
	Locked     field.Bool   // locked(0: unlock, 1: locked)
	LockExpire field.Uint64 // lock expiration time
	Wrong      field.Uint64 // type wrong password count
	Role       userHasOneRole

	fieldMap map[string]field.Expr
}

func (u user) Table(newTableName string) *user {
	u.userDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u user) As(alias string) *user {
	u.userDo.DO = *(u.userDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *user) updateTableName(table string) *user {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewUint64(table, "id")
	u.CreatedAt = field.NewField(table, "created_at")
	u.UpdatedAt = field.NewField(table, "updated_at")
	u.RoleID = field.NewUint64(table, "role_id")
	u.Action = field.NewString(table, "action")
	u.Username = field.NewString(table, "username")
	u.Code = field.NewString(table, "code")
	u.Password = field.NewString(table, "password")
	u.Platform = field.NewString(table, "platform")
	u.LastLogin = field.NewField(table, "last_login")
	u.Locked = field.NewBool(table, "locked")
	u.LockExpire = field.NewUint64(table, "lock_expire")
	u.Wrong = field.NewUint64(table, "wrong")

	u.fillFieldMap()

	return u
}

func (u *user) WithContext(ctx context.Context) *userDo { return u.userDo.WithContext(ctx) }

func (u user) TableName() string { return u.userDo.TableName() }

func (u user) Alias() string { return u.userDo.Alias() }

func (u *user) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *user) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 14)
	u.fieldMap["id"] = u.ID
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["role_id"] = u.RoleID
	u.fieldMap["action"] = u.Action
	u.fieldMap["username"] = u.Username
	u.fieldMap["code"] = u.Code
	u.fieldMap["password"] = u.Password
	u.fieldMap["platform"] = u.Platform
	u.fieldMap["last_login"] = u.LastLogin
	u.fieldMap["locked"] = u.Locked
	u.fieldMap["lock_expire"] = u.LockExpire
	u.fieldMap["wrong"] = u.Wrong

}

func (u user) clone(db *gorm.DB) user {
	u.userDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u user) replaceDB(db *gorm.DB) user {
	u.userDo.ReplaceDB(db)
	return u
}

type userHasOneRole struct {
	db *gorm.DB

	field.RelationField
}

func (a userHasOneRole) Where(conds ...field.Expr) *userHasOneRole {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userHasOneRole) WithContext(ctx context.Context) *userHasOneRole {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userHasOneRole) Session(session *gorm.Session) *userHasOneRole {
	a.db = a.db.Session(session)
	return &a
}

func (a userHasOneRole) Model(m *model.User) *userHasOneRoleTx {
	return &userHasOneRoleTx{a.db.Model(m).Association(a.Name())}
}

type userHasOneRoleTx struct{ tx *gorm.Association }

func (a userHasOneRoleTx) Find() (result *model.Role, err error) {
	return result, a.tx.Find(&result)
}

func (a userHasOneRoleTx) Append(values ...*model.Role) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userHasOneRoleTx) Replace(values ...*model.Role) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userHasOneRoleTx) Delete(values ...*model.Role) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userHasOneRoleTx) Clear() error {
	return a.tx.Clear()
}

func (a userHasOneRoleTx) Count() int64 {
	return a.tx.Count()
}

type userDo struct{ gen.DO }

// SELECT * FROM @@table WHERE id = @id LIMIT 1
func (u userDo) GetByID(id uint64) (result model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM user WHERE id = ? LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// SELECT * FROM @@table
// {{where}}
//
//	{{if val != ""}}
//	  {{if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%")}}
//	    @@col LIKE concat('%', TRIM(BOTH '%' FROM @val), '%')
//	  {{else if strings.HasPrefix(val, "%")}}
//	    @@col LIKE concat('%', TRIM(BOTH '%' FROM @val))
//	  {{else if strings.HasSuffix(val, "%")}}
//	    @@col LIKE concat(TRIM(BOTH '%' FROM @val), '%')
//	  {{else}}
//	    @@col = @val
//	  {{end}}
//	{{end}}
//
// {{end}}
// LIMIT 1
func (u userDo) GetByCol(col string, val string) (result model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user ")
	var whereSQL0 strings.Builder
	if val != "" {
		if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') ")
		} else if strings.HasPrefix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) ")
		} else if strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') ")
		} else {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " = ? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	generateSQL.WriteString("LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// SELECT * FROM @@table
// {{if len(cols) == len(vals)}}
// {{where}}
//
//	  {{for i, col := range cols}}
//	    {{for j, val := range vals}}
//	      {{if i == j}}
//	        {{if val != ""}}
//	          {{if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%")}}
//	            @@col LIKE concat('%', TRIM(BOTH '%' FROM @val), '%') AND
//	          {{else if strings.HasPrefix(val, "%")}}
//	            @@col LIKE concat('%', TRIM(BOTH '%' FROM @val)) AND
//	          {{else if strings.HasSuffix(val, "%")}}
//	            @@col LIKE concat(TRIM(BOTH '%' FROM @val), '%') AND
//	          {{else}}
//	            @@col = @val AND
//	          {{end}}
//	        {{end}}
//	      {{end}}
//	    {{end}}
//	  {{end}}
//	{{end}}
//
// {{end}}
// LIMIT 1
func (u userDo) GetByCols(cols []string, vals []string) (result model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user ")
	if len(cols) == len(vals) {
		var whereSQL0 strings.Builder
		for i, col := range cols {
			for j, val := range vals {
				if i == j {
					if val != "" {
						if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') AND ")
						} else if strings.HasPrefix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) AND ")
						} else if strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') AND ")
						} else {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " = ? AND ")
						}
					}
				}
			}
		}
		helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	}
	generateSQL.WriteString("LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// SELECT * FROM @@table
// {{where}}
//
//	{{if val != ""}}
//	  {{if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%")}}
//	    @@col LIKE concat('%', TRIM(BOTH '%' FROM @val), '%')
//	  {{else if strings.HasPrefix(val, "%")}}
//	    @@col LIKE concat('%', TRIM(BOTH '%' FROM @val))
//	  {{else if strings.HasSuffix(val, "%")}}
//	    @@col LIKE concat(TRIM(BOTH '%' FROM @val), '%')
//	  {{else}}
//	    @@col = @val
//	  {{end}}
//	{{end}}
//
// {{end}}
func (u userDo) FindByCol(col string, val string) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user ")
	var whereSQL0 strings.Builder
	if val != "" {
		if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') ")
		} else if strings.HasPrefix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) ")
		} else if strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') ")
		} else {
			params = append(params, val)
			whereSQL0.WriteString(u.Quote(col) + " = ? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// SELECT * FROM @@table
// {{if len(cols) == len(vals)}}
// {{where}}
//
//	  {{for i, col := range cols}}
//	    {{for j, val := range vals}}
//	      {{if i == j}}
//	        {{if val != ""}}
//	          {{if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%")}}
//	            @@col LIKE concat('%', TRIM(BOTH '%' FROM @val), '%') AND
//	          {{else if strings.HasPrefix(val, "%")}}
//	            @@col LIKE concat('%', TRIM(BOTH '%' FROM @val)) AND
//	          {{else if strings.HasSuffix(val, "%")}}
//	            @@col LIKE concat(TRIM(BOTH '%' FROM @val), '%') AND
//	          {{else}}
//	            @@col = @val AND
//	          {{end}}
//	        {{end}}
//	      {{end}}
//	    {{end}}
//	  {{end}}
//	{{end}}
//
// {{end}}
func (u userDo) FindByCols(cols []string, vals []string) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user ")
	if len(cols) == len(vals) {
		var whereSQL0 strings.Builder
		for i, col := range cols {
			for j, val := range vals {
				if i == j {
					if val != "" {
						if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') AND ")
						} else if strings.HasPrefix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) AND ")
						} else if strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') AND ")
						} else {
							params = append(params, val)
							whereSQL0.WriteString(u.Quote(col) + " = ? AND ")
						}
					}
				}
			}
		}
		helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	}

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

func (u userDo) Debug() *userDo {
	return u.withDO(u.DO.Debug())
}

func (u userDo) WithContext(ctx context.Context) *userDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userDo) ReadDB() *userDo {
	return u.Clauses(dbresolver.Read)
}

func (u userDo) WriteDB() *userDo {
	return u.Clauses(dbresolver.Write)
}

func (u userDo) Session(config *gorm.Session) *userDo {
	return u.withDO(u.DO.Session(config))
}

func (u userDo) Clauses(conds ...clause.Expression) *userDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userDo) Returning(value interface{}, columns ...string) *userDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userDo) Not(conds ...gen.Condition) *userDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userDo) Or(conds ...gen.Condition) *userDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userDo) Select(conds ...field.Expr) *userDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userDo) Where(conds ...gen.Condition) *userDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userDo) Order(conds ...field.Expr) *userDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userDo) Distinct(cols ...field.Expr) *userDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userDo) Omit(cols ...field.Expr) *userDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userDo) Join(table schema.Tabler, on ...field.Expr) *userDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userDo) RightJoin(table schema.Tabler, on ...field.Expr) *userDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userDo) Group(cols ...field.Expr) *userDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userDo) Having(conds ...gen.Condition) *userDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userDo) Limit(limit int) *userDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userDo) Offset(offset int) *userDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userDo) Unscoped() *userDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userDo) Create(values ...*model.User) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userDo) CreateInBatches(values []*model.User, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userDo) Save(values ...*model.User) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userDo) First() (*model.User, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Take() (*model.User, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Last() (*model.User, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Find() ([]*model.User, error) {
	result, err := u.DO.Find()
	return result.([]*model.User), err
}

func (u userDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.User, err error) {
	buf := make([]*model.User, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userDo) FindInBatches(result *[]*model.User, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userDo) Attrs(attrs ...field.AssignExpr) *userDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userDo) Assign(attrs ...field.AssignExpr) *userDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userDo) Joins(fields ...field.RelationField) *userDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userDo) Preload(fields ...field.RelationField) *userDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userDo) FirstOrInit() (*model.User, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) FirstOrCreate() (*model.User, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) FindByPage(offset int, limit int) (result []*model.User, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userDo) Delete(models ...*model.User) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userDo) withDO(do gen.Dao) *userDo {
	u.DO = *do.(*gen.DO)
	return u
}
