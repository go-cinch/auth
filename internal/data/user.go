package data

import (
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
)

type userRepo struct {
	data *Data
}

// User is database fields map
type User struct {
	Id           uint64          `json:"id"`           // auto increment id
	CreatedAt    carbon.DateTime `json:"createdAt"`    // create time
	UpdatedAt    carbon.DateTime `json:"updatedAt"`    // update time
	Username     string          `json:"username"`     // user login name
	Password     string          `json:"password"`     // password
	Mobile       string          `json:"mobile"`       // mobile number
	Avatar       string          `json:"avatar"`       // avatar url
	Nickname     string          `json:"nickname"`     // nickname
	Introduction string          `json:"introduction"` // introduction
	Status       uint64          `json:"status"`       // status(0: disabled, 1: enable)
	LastLogin    carbon.DateTime `json:"lastLogin"`    // last login time
	Locked       uint64          `json:"locked"`       // locked(0: unlock, 1: locked)
	LockExpire   int64           `json:"lockExpire"`   // lock expiration time
	Wrong        int64           `json:"wrong"`        // wrong password count
}

// NewUserRepo .
func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}

func (ro userRepo) GetByUsername(ctx context.Context, username string) (item *biz.User, err error) {
	item = &biz.User{}
	var m User
	db := ro.data.DB(ctx)
	db.
		Where("username = ?", username).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.UserNotFound
		return
	}
	copier.Copy(item, m)
	return
}

func (ro userRepo) Create(ctx context.Context, item *biz.User) (err error) {
	var m User
	db := ro.data.DB(ctx)
	db.
		Where("username = ?", item.Username).
		First(&m)
	if m.Id > constant.UI0 {
		err = biz.DuplicateUsername
		return
	}
	copier.Copy(&m, item)
	err = db.Create(&m).Error
	return
}

func (ro userRepo) UpdatePassword(ctx context.Context, item *biz.User) (err error) {
	var m User
	db := ro.data.DB(ctx)
	db.
		Where("username = ?", item.Username).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.UserNotFound
		return
	}
	err = db.
		Model(&User{}).
		Where("id = ?", m.Id).
		Update("password", item.Password).Error
	return
}
