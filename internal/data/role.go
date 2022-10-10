package data

import (
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/jinzhu/copier"
)

type roleRepo struct {
	data *Data
}

// Role is database fields map
type Role struct {
	Id     uint64 `json:"id,string"` // auto increment id
	Name   string `json:"name"`      // name
	Key    string `json:"key"`       // keyword, must be unique, used as frontend display
	Status uint64 `json:"status"`    // status(0: disabled, 1: enable)
}

func NewRoleRepo(data *Data) biz.RoleRepo {
	return &roleRepo{
		data: data,
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
	copier.Copy(&m, item)
	m.Id = ro.data.Id(ctx)
	m.Status = constant.UI1
	err = db.Create(&m).Error
	return
}
