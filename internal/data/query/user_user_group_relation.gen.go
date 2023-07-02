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

func newUserUserGroupRelation(db *gorm.DB, opts ...gen.DOOption) userUserGroupRelation {
	_userUserGroupRelation := userUserGroupRelation{}

	_userUserGroupRelation.userUserGroupRelationDo.UseDB(db, opts...)
	_userUserGroupRelation.userUserGroupRelationDo.UseModel(&model.UserUserGroupRelation{})

	tableName := _userUserGroupRelation.userUserGroupRelationDo.TableName()
	_userUserGroupRelation.ALL = field.NewAsterisk(tableName)
	_userUserGroupRelation.UserID = field.NewUint64(tableName, "user_id")
	_userUserGroupRelation.UserGroupID = field.NewUint64(tableName, "user_group_id")

	_userUserGroupRelation.fillFieldMap()

	return _userUserGroupRelation
}

type userUserGroupRelation struct {
	userUserGroupRelationDo userUserGroupRelationDo

	ALL         field.Asterisk
	UserID      field.Uint64 // auto increment id
	UserGroupID field.Uint64 // auto increment id

	fieldMap map[string]field.Expr
}

func (u userUserGroupRelation) Table(newTableName string) *userUserGroupRelation {
	u.userUserGroupRelationDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userUserGroupRelation) As(alias string) *userUserGroupRelation {
	u.userUserGroupRelationDo.DO = *(u.userUserGroupRelationDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userUserGroupRelation) updateTableName(table string) *userUserGroupRelation {
	u.ALL = field.NewAsterisk(table)
	u.UserID = field.NewUint64(table, "user_id")
	u.UserGroupID = field.NewUint64(table, "user_group_id")

	u.fillFieldMap()

	return u
}

func (u *userUserGroupRelation) WithContext(ctx context.Context) *userUserGroupRelationDo {
	return u.userUserGroupRelationDo.WithContext(ctx)
}

func (u userUserGroupRelation) TableName() string { return u.userUserGroupRelationDo.TableName() }

func (u userUserGroupRelation) Alias() string { return u.userUserGroupRelationDo.Alias() }

func (u *userUserGroupRelation) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userUserGroupRelation) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 2)
	u.fieldMap["user_id"] = u.UserID
	u.fieldMap["user_group_id"] = u.UserGroupID
}

func (u userUserGroupRelation) clone(db *gorm.DB) userUserGroupRelation {
	u.userUserGroupRelationDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userUserGroupRelation) replaceDB(db *gorm.DB) userUserGroupRelation {
	u.userUserGroupRelationDo.ReplaceDB(db)
	return u
}

type userUserGroupRelationDo struct{ gen.DO }

// SELECT * FROM @@table WHERE id = @id LIMIT 1
func (u userUserGroupRelationDo) GetByID(id uint64) (result model.UserUserGroupRelation) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM user_user_group_relation WHERE id = ? LIMIT 1 ")

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
func (u userUserGroupRelationDo) GetByCol(col string, val string) (result model.UserUserGroupRelation) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user_user_group_relation ")
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
func (u userUserGroupRelationDo) GetByCols(cols []string, vals []string) (result model.UserUserGroupRelation) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user_user_group_relation ")
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
func (u userUserGroupRelationDo) FindByCol(col string, val string) (result []model.UserUserGroupRelation) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user_user_group_relation ")
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
func (u userUserGroupRelationDo) FindByCols(cols []string, vals []string) (result []model.UserUserGroupRelation) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user_user_group_relation ")
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

func (u userUserGroupRelationDo) Debug() *userUserGroupRelationDo {
	return u.withDO(u.DO.Debug())
}

func (u userUserGroupRelationDo) WithContext(ctx context.Context) *userUserGroupRelationDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userUserGroupRelationDo) ReadDB() *userUserGroupRelationDo {
	return u.Clauses(dbresolver.Read)
}

func (u userUserGroupRelationDo) WriteDB() *userUserGroupRelationDo {
	return u.Clauses(dbresolver.Write)
}

func (u userUserGroupRelationDo) Session(config *gorm.Session) *userUserGroupRelationDo {
	return u.withDO(u.DO.Session(config))
}

func (u userUserGroupRelationDo) Clauses(conds ...clause.Expression) *userUserGroupRelationDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userUserGroupRelationDo) Returning(value interface{}, columns ...string) *userUserGroupRelationDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userUserGroupRelationDo) Not(conds ...gen.Condition) *userUserGroupRelationDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userUserGroupRelationDo) Or(conds ...gen.Condition) *userUserGroupRelationDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userUserGroupRelationDo) Select(conds ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userUserGroupRelationDo) Where(conds ...gen.Condition) *userUserGroupRelationDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userUserGroupRelationDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userUserGroupRelationDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userUserGroupRelationDo) Order(conds ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userUserGroupRelationDo) Distinct(cols ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userUserGroupRelationDo) Omit(cols ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userUserGroupRelationDo) Join(table schema.Tabler, on ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userUserGroupRelationDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userUserGroupRelationDo) RightJoin(table schema.Tabler, on ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userUserGroupRelationDo) Group(cols ...field.Expr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userUserGroupRelationDo) Having(conds ...gen.Condition) *userUserGroupRelationDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userUserGroupRelationDo) Limit(limit int) *userUserGroupRelationDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userUserGroupRelationDo) Offset(offset int) *userUserGroupRelationDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userUserGroupRelationDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userUserGroupRelationDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userUserGroupRelationDo) Unscoped() *userUserGroupRelationDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userUserGroupRelationDo) Create(values ...*model.UserUserGroupRelation) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userUserGroupRelationDo) CreateInBatches(values []*model.UserUserGroupRelation, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userUserGroupRelationDo) Save(values ...*model.UserUserGroupRelation) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userUserGroupRelationDo) First() (*model.UserUserGroupRelation, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserUserGroupRelation), nil
	}
}

func (u userUserGroupRelationDo) Take() (*model.UserUserGroupRelation, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserUserGroupRelation), nil
	}
}

func (u userUserGroupRelationDo) Last() (*model.UserUserGroupRelation, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserUserGroupRelation), nil
	}
}

func (u userUserGroupRelationDo) Find() ([]*model.UserUserGroupRelation, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserUserGroupRelation), err
}

func (u userUserGroupRelationDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserUserGroupRelation, err error) {
	buf := make([]*model.UserUserGroupRelation, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userUserGroupRelationDo) FindInBatches(result *[]*model.UserUserGroupRelation, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userUserGroupRelationDo) Attrs(attrs ...field.AssignExpr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userUserGroupRelationDo) Assign(attrs ...field.AssignExpr) *userUserGroupRelationDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userUserGroupRelationDo) Joins(fields ...field.RelationField) *userUserGroupRelationDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userUserGroupRelationDo) Preload(fields ...field.RelationField) *userUserGroupRelationDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userUserGroupRelationDo) FirstOrInit() (*model.UserUserGroupRelation, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserUserGroupRelation), nil
	}
}

func (u userUserGroupRelationDo) FirstOrCreate() (*model.UserUserGroupRelation, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserUserGroupRelation), nil
	}
}

func (u userUserGroupRelationDo) FindByPage(offset int, limit int) (result []*model.UserUserGroupRelation, count int64, err error) {
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

func (u userUserGroupRelationDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userUserGroupRelationDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userUserGroupRelationDo) Delete(models ...*model.UserUserGroupRelation) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userUserGroupRelationDo) withDO(do gen.Dao) *userUserGroupRelationDo {
	u.DO = *do.(*gen.DO)
	return u
}