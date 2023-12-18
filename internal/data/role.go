package data

import (
	"context"
	"strings"

	"auth/internal/biz"
	"auth/internal/data/model"
	"auth/internal/data/query"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"gorm.io/gen"
)

type roleRepo struct {
	data   *Data
	action biz.ActionRepo
}

func NewRoleRepo(data *Data, action biz.ActionRepo) biz.RoleRepo {
	return &roleRepo{
		data:   data,
		action: action,
	}
}

func (ro roleRepo) Create(ctx context.Context, item *biz.Role) (err error) {
	ok := ro.WordExists(ctx, item.Word)
	if ok {
		err = biz.ErrDuplicateField(ctx, "word", item.Word)
		return
	}
	var m model.Role
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).Role
	db := p.WithContext(ctx)
	m.ID = ro.data.Id(ctx)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	err = db.Create(&m)
	return
}

func (ro roleRepo) Find(ctx context.Context, condition *biz.FindRole) (rp []biz.Role) {
	p := query.Use(ro.data.DB(ctx)).Role
	db := p.WithContext(ctx)
	rp = make([]biz.Role, 0)
	list := make([]model.Role, 0)
	conditions := make([]gen.Condition, 0, 2)
	if condition.Name != nil {
		conditions = append(conditions, p.Name.Like(strings.Join([]string{"%", *condition.Name, "%"}, "")))
	}
	if condition.Word != nil {
		conditions = append(conditions, p.Word.Like(strings.Join([]string{"%", *condition.Word, "%"}, "")))
	}
	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(
			db.
				Order(p.ID.Desc()).
				Where(conditions...).
				UnderlyingDB(),
		).
		Find(&list)
	copierx.Copy(&rp, list)
	for i, item := range rp {
		rp[i].Actions = make([]biz.Action, 0)
		arr := ro.action.FindByCode(ctx, item.Action)
		copierx.Copy(&rp[i].Actions, arr)
	}
	return
}

func (ro roleRepo) Update(ctx context.Context, item *biz.UpdateRole) (err error) {
	p := query.Use(ro.data.DB(ctx)).Role
	db := p.WithContext(ctx)
	m := db.GetByID(item.Id)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.ErrDataNotChange(ctx)
		return
	}
	if a, ok1 := change["action"]; ok1 {
		if v, ok2 := a.(string); ok2 {
			err = ro.action.CodeExists(ctx, v)
			if err != nil {
				return
			}
		}
	}
	if item.Word != nil && *item.Word != m.Word {
		ok := ro.WordExists(ctx, *item.Word)
		if ok {
			err = biz.ErrDuplicateField(ctx, "word", *item.Word)
			return
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.Id)).
		Updates(&change)
	return
}

func (ro roleRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).Role
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro roleRepo) WordExists(ctx context.Context, word string) (ok bool) {
	p := query.Use(ro.data.DB(ctx)).Role
	db := p.WithContext(ctx)
	arr := strings.Split(word, ",")
	for _, item := range arr {
		m := db.GetByCol("word", item)
		if m.ID == constant.UI0 {
			log.Error("invalid word: %s", item)
			return
		}
	}
	ok = true
	return
}
