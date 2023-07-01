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

func newAction(db *gorm.DB, opts ...gen.DOOption) action {
	_action := action{}

	_action.actionDo.UseDB(db, opts...)
	_action.actionDo.UseModel(&model.Action{})

	tableName := _action.actionDo.TableName()
	_action.ALL = field.NewAsterisk(tableName)
	_action.ID = field.NewUint64(tableName, "id")
	_action.Name = field.NewString(tableName, "name")
	_action.Code = field.NewString(tableName, "code")
	_action.Word = field.NewString(tableName, "word")
	_action.Resource = field.NewString(tableName, "resource")
	_action.Menu = field.NewString(tableName, "menu")
	_action.Btn = field.NewString(tableName, "btn")

	_action.fillFieldMap()

	return _action
}

type action struct {
	actionDo actionDo

	ALL  field.Asterisk
	ID   field.Uint64 // auto increment id
	Name field.String // name
	Code field.String // code
	Word field.String // keyword, must be unique, used as frontend display
	/*
		resource array, split by break line str, example: GET|/user+
		+PUT,PATCH|/role/*+
		+GET|/action
	*/
	Resource field.String
	Menu     field.String // menu array, split by break line str
	Btn      field.String // btn array, split by break line str

	fieldMap map[string]field.Expr
}

func (a action) Table(newTableName string) *action {
	a.actionDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a action) As(alias string) *action {
	a.actionDo.DO = *(a.actionDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *action) updateTableName(table string) *action {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewUint64(table, "id")
	a.Name = field.NewString(table, "name")
	a.Code = field.NewString(table, "code")
	a.Word = field.NewString(table, "word")
	a.Resource = field.NewString(table, "resource")
	a.Menu = field.NewString(table, "menu")
	a.Btn = field.NewString(table, "btn")

	a.fillFieldMap()

	return a
}

func (a *action) WithContext(ctx context.Context) *actionDo { return a.actionDo.WithContext(ctx) }

func (a action) TableName() string { return a.actionDo.TableName() }

func (a action) Alias() string { return a.actionDo.Alias() }

func (a *action) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *action) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 7)
	a.fieldMap["id"] = a.ID
	a.fieldMap["name"] = a.Name
	a.fieldMap["code"] = a.Code
	a.fieldMap["word"] = a.Word
	a.fieldMap["resource"] = a.Resource
	a.fieldMap["menu"] = a.Menu
	a.fieldMap["btn"] = a.Btn
}

func (a action) clone(db *gorm.DB) action {
	a.actionDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a action) replaceDB(db *gorm.DB) action {
	a.actionDo.ReplaceDB(db)
	return a
}

type actionDo struct{ gen.DO }

// SELECT * FROM @@table WHERE id = @id LIMIT 1
func (a actionDo) GetByID(id uint64) (result model.Action) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM action WHERE id = ? LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
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
func (a actionDo) GetByCol(col string, val string) (result model.Action) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM action ")
	var whereSQL0 strings.Builder
	if val != "" {
		if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') ")
		} else if strings.HasPrefix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) ")
		} else if strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') ")
		} else {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " = ? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	generateSQL.WriteString("LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
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
func (a actionDo) GetByCols(cols []string, vals []string) (result model.Action) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM action ")
	if len(cols) == len(vals) {
		var whereSQL0 strings.Builder
		for i, col := range cols {
			for j, val := range vals {
				if i == j {
					if val != "" {
						if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') AND ")
						} else if strings.HasPrefix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) AND ")
						} else if strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') AND ")
						} else {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " = ? AND ")
						}
					}
				}
			}
		}
		helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	}
	generateSQL.WriteString("LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
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
func (a actionDo) FindByCol(col string, val string) (result []model.Action) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM action ")
	var whereSQL0 strings.Builder
	if val != "" {
		if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') ")
		} else if strings.HasPrefix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) ")
		} else if strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') ")
		} else {
			params = append(params, val)
			whereSQL0.WriteString(a.Quote(col) + " = ? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
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
func (a actionDo) FindByCols(cols []string, vals []string) (result []model.Action) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM action ")
	if len(cols) == len(vals) {
		var whereSQL0 strings.Builder
		for i, col := range cols {
			for j, val := range vals {
				if i == j {
					if val != "" {
						if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') AND ")
						} else if strings.HasPrefix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) AND ")
						} else if strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') AND ")
						} else {
							params = append(params, val)
							whereSQL0.WriteString(a.Quote(col) + " = ? AND ")
						}
					}
				}
			}
		}
		helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	}

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

func (a actionDo) Debug() *actionDo {
	return a.withDO(a.DO.Debug())
}

func (a actionDo) WithContext(ctx context.Context) *actionDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a actionDo) ReadDB() *actionDo {
	return a.Clauses(dbresolver.Read)
}

func (a actionDo) WriteDB() *actionDo {
	return a.Clauses(dbresolver.Write)
}

func (a actionDo) Session(config *gorm.Session) *actionDo {
	return a.withDO(a.DO.Session(config))
}

func (a actionDo) Clauses(conds ...clause.Expression) *actionDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a actionDo) Returning(value interface{}, columns ...string) *actionDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a actionDo) Not(conds ...gen.Condition) *actionDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a actionDo) Or(conds ...gen.Condition) *actionDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a actionDo) Select(conds ...field.Expr) *actionDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a actionDo) Where(conds ...gen.Condition) *actionDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a actionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *actionDo {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a actionDo) Order(conds ...field.Expr) *actionDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a actionDo) Distinct(cols ...field.Expr) *actionDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a actionDo) Omit(cols ...field.Expr) *actionDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a actionDo) Join(table schema.Tabler, on ...field.Expr) *actionDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a actionDo) LeftJoin(table schema.Tabler, on ...field.Expr) *actionDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a actionDo) RightJoin(table schema.Tabler, on ...field.Expr) *actionDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a actionDo) Group(cols ...field.Expr) *actionDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a actionDo) Having(conds ...gen.Condition) *actionDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a actionDo) Limit(limit int) *actionDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a actionDo) Offset(offset int) *actionDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a actionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *actionDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a actionDo) Unscoped() *actionDo {
	return a.withDO(a.DO.Unscoped())
}

func (a actionDo) Create(values ...*model.Action) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a actionDo) CreateInBatches(values []*model.Action, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a actionDo) Save(values ...*model.Action) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a actionDo) First() (*model.Action, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Action), nil
	}
}

func (a actionDo) Take() (*model.Action, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Action), nil
	}
}

func (a actionDo) Last() (*model.Action, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Action), nil
	}
}

func (a actionDo) Find() ([]*model.Action, error) {
	result, err := a.DO.Find()
	return result.([]*model.Action), err
}

func (a actionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Action, err error) {
	buf := make([]*model.Action, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a actionDo) FindInBatches(result *[]*model.Action, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a actionDo) Attrs(attrs ...field.AssignExpr) *actionDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a actionDo) Assign(attrs ...field.AssignExpr) *actionDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a actionDo) Joins(fields ...field.RelationField) *actionDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a actionDo) Preload(fields ...field.RelationField) *actionDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a actionDo) FirstOrInit() (*model.Action, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Action), nil
	}
}

func (a actionDo) FirstOrCreate() (*model.Action, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Action), nil
	}
}

func (a actionDo) FindByPage(offset int, limit int) (result []*model.Action, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a actionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a actionDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a actionDo) Delete(models ...*model.Action) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *actionDo) withDO(do gen.Dao) *actionDo {
	a.DO = *do.(*gen.DO)
	return a
}
