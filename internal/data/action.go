package data

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"fmt"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/utils"
	"strings"
)

type actionRepo struct {
	data *Data
}

// Action is database fields map
type Action struct {
	Id       uint64 `json:"id,string"` // auto increment id
	Code     string `json:"code"`      // unique code
	Name     string `json:"name"`      // name
	Word     string `json:"word"`      // keyword, must be unique, used as frontend display
	Resource string `json:"resource"`  // resource array, split by break line str, example: GET,/user+\n+POST,/role+\n+GET,/action
}

func NewActionRepo(data *Data) biz.ActionRepo {
	return &actionRepo{
		data: data,
	}
}

func (ro actionRepo) Create(ctx context.Context, item *biz.Action) (err error) {
	var m Action
	err = ro.WordExists(ctx, item.Word)
	if err == nil {
		err = biz.DuplicateActionWord
		return
	}
	copierx.Copy(&m, item)
	db := ro.data.DB(ctx)
	m.Id = ro.data.Id(ctx)
	m.Code = id.NewCode(m.Id)
	if m.Resource == "" {
		m.Resource = "*"
	}
	err = db.Create(&m).Error
	return
}

func (ro actionRepo) Find(ctx context.Context, condition *biz.FindAction) (rp []biz.Action, err error) {
	db := ro.data.DB(ctx)
	db = db.
		Model(&Action{}).
		Order("id DESC")
	rp = make([]biz.Action, 0)
	list := make([]Action, 0)
	if condition.Name != nil {
		db.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Name))
	}
	if condition.Code != nil {
		db.Where("`code` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Code))
	}
	if condition.Word != nil {
		db.Where("`word` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Word))
	}
	if condition.Resource != nil {
		db.Where("`resource` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Resource))
	}
	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(db).
		Find(&list)
	copierx.Copy(&rp, list)
	return
}

func (ro actionRepo) FindByCode(ctx context.Context, code string) (rp []biz.Action, err error) {
	rp = make([]biz.Action, 0)
	list := make([]Action, 0)
	db := ro.data.DB(ctx)
	arr := strings.Split(code, ",")
	db.
		Model(&Action{}).
		Where("code IN (?)", arr).
		Find(&list)
	copierx.Copy(&rp, list)
	return
}

func (ro actionRepo) Update(ctx context.Context, item *biz.UpdateAction) (err error) {
	var m Action
	db := ro.data.DB(ctx)
	db.
		Where("id = ?", item.Id).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.ActionNotFound
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.DataNotChange
		return
	}
	if item.Word != nil && *item.Word != m.Word {
		err = ro.WordExists(ctx, *item.Word)
		if err == nil {
			err = biz.DuplicateActionWord
			return
		}
	}
	err = db.
		Model(&m).
		Updates(&change).Error
	return
}

func (ro actionRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	db := ro.data.DB(ctx)
	err = db.
		Where("id IN (?)", ids).
		Delete(&Action{}).Error
	if err != nil {
		return
	}
	var count int64
	db.
		Model(&Action{}).
		Count(&count)
	if count == 0 {
		err = biz.KeepLeastOntAction
	}
	return
}

func (ro actionRepo) CodeExists(ctx context.Context, code string) (err error) {
	var m Action
	db := ro.data.DB(ctx)
	arr := strings.Split(code, ",")
	for _, item := range arr {
		db.
			Where("code = ?", item).
			First(&m)
		ok := m.Id > constant.UI1
		if !ok {
			err = v1.ErrorIllegalParameter("%s: %s", biz.ActionNotFound.Message, item)
			return
		}
	}
	return
}

func (ro actionRepo) WordExists(ctx context.Context, word string) (err error) {
	var m Action
	db := ro.data.DB(ctx)
	arr := strings.Split(word, ",")
	for _, item := range arr {
		db.
			Where("word = ?", item).
			First(&m)
		ok := m.Id > constant.UI1
		if !ok {
			err = v1.ErrorIllegalParameter("%s: %s", biz.ActionNotFound.Message, item)
			return
		}
	}
	return
}

func (ro actionRepo) Permission(ctx context.Context, code, resource string) (pass bool) {
	arr := strings.Split(code, ",")
	for _, item := range arr {
		pass = ro.permission(ctx, item, resource)
		if pass {
			return
		}
	}
	return
}

func (ro actionRepo) permission(ctx context.Context, code, resource string) (pass bool) {
	var m Action
	db := ro.data.DB(ctx)
	db.
		Where("code = ?", code).
		First(&m)
	if m.Id == constant.UI0 {
		return
	}
	arr := strings.Split(m.Resource, "\n")
	for _, v := range arr {
		if v == "*" || v == resource {
			pass = true
			return
		}
	}
	return
}
