package data

import (
	"auth/internal/biz"
	"context"
	"fmt"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/utils"
	"strings"
)

type roleRepo struct {
	data   *Data
	action biz.ActionRepo
}

// Role is database fields map
type Role struct {
	Id     uint64 `json:"id,string"` // auto increment id
	Name   string `json:"name"`      // name
	Word   string `json:"word"`      // keyword, must be unique, used as frontend display
	Action string `json:"action"`    // role action code array
}

func NewRoleRepo(data *Data, action biz.ActionRepo) biz.RoleRepo {
	return &roleRepo{
		data:   data,
		action: action,
	}
}

func (ro roleRepo) Create(ctx context.Context, item *biz.Role) (err error) {
	var m Role
	err = ro.WordExists(ctx, item.Word)
	if err == nil {
		err = biz.IllegalParameter("%s `word`: %s", biz.DuplicateField.Message, item.Word)
		return
	}
	copierx.Copy(&m, item)
	db := ro.data.DB(ctx)
	m.Id = ro.data.Id(ctx)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	err = db.Create(&m).Error
	return
}

func (ro roleRepo) Find(ctx context.Context, condition *biz.FindRole) (rp []biz.Role) {
	db := ro.data.DB(ctx)
	db = db.
		Model(&Role{}).
		Order("id DESC")
	rp = make([]biz.Role, 0)
	list := make([]Role, 0)
	if condition.Name != nil {
		db.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Name))
	}
	if condition.Word != nil {
		db.Where("`word` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Word))
	}
	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(db).
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
	var m Role
	db := ro.data.DB(ctx)
	db.
		Where("`id` = ?", item.Id).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.NotFound("%s Role.id: %d", biz.RecordNotFound.Message, item.Id)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.DataNotChange
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
		err = ro.WordExists(ctx, *item.Word)
		if err == nil {
			err = biz.IllegalParameter("%s `word`: %s", biz.DuplicateField.Message, *item.Word)
			return
		}
	}
	err = db.
		Model(&m).
		Updates(&change).Error
	return
}

func (ro roleRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	db := ro.data.DB(ctx)
	err = db.
		Where("`id` IN (?)", ids).
		Delete(&Role{}).Error
	return
}

func (ro roleRepo) WordExists(ctx context.Context, word string) (err error) {
	var m Role
	db := ro.data.DB(ctx)
	arr := strings.Split(word, ",")
	for _, item := range arr {
		db.
			Where("`word` = ?", item).
			First(&m)
		if m.Id == constant.UI0 {
			err = biz.NotFound("%s Role.word: %s", biz.RecordNotFound.Message, item)
			return
		}
	}
	return
}
