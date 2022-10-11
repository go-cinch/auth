package data

import (
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/jinzhu/copier"
)

type userGroupRepo struct {
	data   *Data
	action biz.ActionRepo
	user   biz.UserRepo
}

// UserGroup is database fields map
type UserGroup struct {
	Id     uint64 `json:"id,string"`                                      // auto increment id
	Users  []User `gorm:"many2many:user_user_group_relation" json:"user"` // User and UserGroup many2many relation
	Name   string `json:"name"`                                           // name
	Key    string `json:"key"`                                            // keyword, must be unique, used as frontend display
	Action string `json:"action"`                                         // user group action code array
}

func NewUserGroupRepo(data *Data, action biz.ActionRepo, user biz.UserRepo) biz.UserGroupRepo {
	return &userGroupRepo{
		data:   data,
		action: action,
		user:   user,
	}
}

func (ro userGroupRepo) Create(ctx context.Context, item *biz.UserGroup) (err error) {
	var m UserGroup
	db := ro.data.DB(ctx)
	db.
		Where("`key` = ?", item.Key).
		First(&m)
	if m.Id > constant.UI0 {
		err = biz.DuplicateUserGroupKey
		return
	}
	copier.Copy(&m, item)
	m.Id = ro.data.Id(ctx)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	if len(item.Users) > 0 {
		m.Users = make([]User, 0)
		for _, id := range item.Users {
			err = ro.user.IdExists(ctx, id)
			if err != nil {
				return
			}
			m.Users = append(m.Users, User{
				Id: id,
			})
		}
	}
	err = db.Create(&m).Error
	return
}
