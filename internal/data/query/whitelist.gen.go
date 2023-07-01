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

func newWhitelist(db *gorm.DB, opts ...gen.DOOption) whitelist {
	_whitelist := whitelist{}

	_whitelist.whitelistDo.UseDB(db, opts...)
	_whitelist.whitelistDo.UseModel(&model.Whitelist{})

	tableName := _whitelist.whitelistDo.TableName()
	_whitelist.ALL = field.NewAsterisk(tableName)
	_whitelist.ID = field.NewUint64(tableName, "id")
	_whitelist.Category = field.NewUint32(tableName, "category")
	_whitelist.Resource = field.NewString(tableName, "resource")

	_whitelist.fillFieldMap()

	return _whitelist
}

type whitelist struct {
	whitelistDo whitelistDo

	ALL      field.Asterisk
	ID       field.Uint64 // auto increment id
	Category field.Uint32 // category(0:permission, 1:jwt, 2:idempotent)
	/*
		resource array, split by break line str, example: GET|/user+
		+PUT,PATCH|/role/*+
		+GET|/action
	*/
	Resource field.String

	fieldMap map[string]field.Expr
}

func (w whitelist) Table(newTableName string) *whitelist {
	w.whitelistDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w whitelist) As(alias string) *whitelist {
	w.whitelistDo.DO = *(w.whitelistDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *whitelist) updateTableName(table string) *whitelist {
	w.ALL = field.NewAsterisk(table)
	w.ID = field.NewUint64(table, "id")
	w.Category = field.NewUint32(table, "category")
	w.Resource = field.NewString(table, "resource")

	w.fillFieldMap()

	return w
}

func (w *whitelist) WithContext(ctx context.Context) *whitelistDo {
	return w.whitelistDo.WithContext(ctx)
}

func (w whitelist) TableName() string { return w.whitelistDo.TableName() }

func (w whitelist) Alias() string { return w.whitelistDo.Alias() }

func (w *whitelist) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *whitelist) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 3)
	w.fieldMap["id"] = w.ID
	w.fieldMap["category"] = w.Category
	w.fieldMap["resource"] = w.Resource
}

func (w whitelist) clone(db *gorm.DB) whitelist {
	w.whitelistDo.ReplaceConnPool(db.Statement.ConnPool)
	return w
}

func (w whitelist) replaceDB(db *gorm.DB) whitelist {
	w.whitelistDo.ReplaceDB(db)
	return w
}

type whitelistDo struct{ gen.DO }

// SELECT * FROM @@table WHERE id = @id LIMIT 1
func (w whitelistDo) GetByID(id uint64) (result model.Whitelist) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM whitelist WHERE id = ? LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = w.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
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
func (w whitelistDo) GetByCol(col string, val string) (result model.Whitelist) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM whitelist ")
	var whereSQL0 strings.Builder
	if val != "" {
		if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') ")
		} else if strings.HasPrefix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) ")
		} else if strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') ")
		} else {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " = ? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	generateSQL.WriteString("LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = w.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
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
func (w whitelistDo) GetByCols(cols []string, vals []string) (result model.Whitelist) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM whitelist ")
	if len(cols) == len(vals) {
		var whereSQL0 strings.Builder
		for i, col := range cols {
			for j, val := range vals {
				if i == j {
					if val != "" {
						if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') AND ")
						} else if strings.HasPrefix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) AND ")
						} else if strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') AND ")
						} else {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " = ? AND ")
						}
					}
				}
			}
		}
		helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	}
	generateSQL.WriteString("LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = w.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
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
func (w whitelistDo) FindByCol(col string, val string) (result []model.Whitelist) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM whitelist ")
	var whereSQL0 strings.Builder
	if val != "" {
		if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') ")
		} else if strings.HasPrefix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) ")
		} else if strings.HasSuffix(val, "%") {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') ")
		} else {
			params = append(params, val)
			whereSQL0.WriteString(w.Quote(col) + " = ? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = w.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
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
func (w whitelistDo) FindByCols(cols []string, vals []string) (result []model.Whitelist) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM whitelist ")
	if len(cols) == len(vals) {
		var whereSQL0 strings.Builder
		for i, col := range cols {
			for j, val := range vals {
				if i == j {
					if val != "" {
						if strings.HasPrefix(val, "%") && strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?), '%') AND ")
						} else if strings.HasPrefix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " LIKE concat('%', TRIM(BOTH '%' FROM ?)) AND ")
						} else if strings.HasSuffix(val, "%") {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " LIKE concat(TRIM(BOTH '%' FROM ?), '%') AND ")
						} else {
							params = append(params, val)
							whereSQL0.WriteString(w.Quote(col) + " = ? AND ")
						}
					}
				}
			}
		}
		helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	}

	var executeSQL *gorm.DB
	executeSQL = w.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

func (w whitelistDo) Debug() *whitelistDo {
	return w.withDO(w.DO.Debug())
}

func (w whitelistDo) WithContext(ctx context.Context) *whitelistDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w whitelistDo) ReadDB() *whitelistDo {
	return w.Clauses(dbresolver.Read)
}

func (w whitelistDo) WriteDB() *whitelistDo {
	return w.Clauses(dbresolver.Write)
}

func (w whitelistDo) Session(config *gorm.Session) *whitelistDo {
	return w.withDO(w.DO.Session(config))
}

func (w whitelistDo) Clauses(conds ...clause.Expression) *whitelistDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w whitelistDo) Returning(value interface{}, columns ...string) *whitelistDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w whitelistDo) Not(conds ...gen.Condition) *whitelistDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w whitelistDo) Or(conds ...gen.Condition) *whitelistDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w whitelistDo) Select(conds ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w whitelistDo) Where(conds ...gen.Condition) *whitelistDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w whitelistDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *whitelistDo {
	return w.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (w whitelistDo) Order(conds ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w whitelistDo) Distinct(cols ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w whitelistDo) Omit(cols ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w whitelistDo) Join(table schema.Tabler, on ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w whitelistDo) LeftJoin(table schema.Tabler, on ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w whitelistDo) RightJoin(table schema.Tabler, on ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w whitelistDo) Group(cols ...field.Expr) *whitelistDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w whitelistDo) Having(conds ...gen.Condition) *whitelistDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w whitelistDo) Limit(limit int) *whitelistDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w whitelistDo) Offset(offset int) *whitelistDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w whitelistDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *whitelistDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w whitelistDo) Unscoped() *whitelistDo {
	return w.withDO(w.DO.Unscoped())
}

func (w whitelistDo) Create(values ...*model.Whitelist) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w whitelistDo) CreateInBatches(values []*model.Whitelist, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w whitelistDo) Save(values ...*model.Whitelist) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w whitelistDo) First() (*model.Whitelist, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Whitelist), nil
	}
}

func (w whitelistDo) Take() (*model.Whitelist, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Whitelist), nil
	}
}

func (w whitelistDo) Last() (*model.Whitelist, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Whitelist), nil
	}
}

func (w whitelistDo) Find() ([]*model.Whitelist, error) {
	result, err := w.DO.Find()
	return result.([]*model.Whitelist), err
}

func (w whitelistDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Whitelist, err error) {
	buf := make([]*model.Whitelist, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w whitelistDo) FindInBatches(result *[]*model.Whitelist, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w whitelistDo) Attrs(attrs ...field.AssignExpr) *whitelistDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w whitelistDo) Assign(attrs ...field.AssignExpr) *whitelistDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w whitelistDo) Joins(fields ...field.RelationField) *whitelistDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w whitelistDo) Preload(fields ...field.RelationField) *whitelistDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w whitelistDo) FirstOrInit() (*model.Whitelist, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Whitelist), nil
	}
}

func (w whitelistDo) FirstOrCreate() (*model.Whitelist, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Whitelist), nil
	}
}

func (w whitelistDo) FindByPage(offset int, limit int) (result []*model.Whitelist, count int64, err error) {
	result, err = w.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = w.Offset(-1).Limit(-1).Count()
	return
}

func (w whitelistDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w whitelistDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w whitelistDo) Delete(models ...*model.Whitelist) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *whitelistDo) withDO(do gen.Dao) *whitelistDo {
	w.DO = *do.(*gen.DO)
	return w
}
