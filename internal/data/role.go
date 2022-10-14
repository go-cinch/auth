package data

import (
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/utils"
)

type roleRepo struct {
	data   *Data
	action biz.ActionRepo
}

// Role is database fields map
type Role struct {
	Id     uint64 `json:"id,string"` // auto increment id
	Name   string `json:"name"`      // name
	Key    string `json:"key"`       // keyword, must be unique, used as frontend display
	Status uint64 `json:"status"`    // status(0: disabled, 1: enable)
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
	db := ro.data.DB(ctx)
	db.
		Where("`key` = ?", item.Key).
		First(&m)
	if m.Id > constant.UI0 {
		err = biz.DuplicateRoleKey
		return
	}
	copierx.Copy(&m, item)
	m.Id = ro.data.Id(ctx)
	m.Status = constant.UI1
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	err = db.Create(&m).Error
	return
}

func (ro roleRepo) Update(ctx context.Context, item *biz.UpdateRole) (err error) {
	var m Role
	db := ro.data.DB(ctx)
	db.
		Where("id = ?", item.Id).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.RoleNotFound
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
	err = db.
		Model(&m).
		Updates(&change).Error
	return
}
