package data

import (
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/id"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
)

type userRepo struct {
	data   *Data
	action biz.ActionRepo
}

// User is database fields map
type User struct {
	Id           uint64          `json:"id,string"`     // auto increment id
	CreatedAt    carbon.DateTime `json:"createdAt"`     // create time
	UpdatedAt    carbon.DateTime `json:"updatedAt"`     // update time
	RoleId       uint64          `json:"roleId,string"` // role id
	Role         Role            `json:"role"`          // role
	Action       string          `json:"action"`        // user action code array
	Username     string          `json:"username"`      // user login name
	UserCode     string          `json:"userCode"`      // user code
	Password     string          `json:"password"`      // password
	Mobile       string          `json:"mobile"`        // mobile number
	Avatar       string          `json:"avatar"`        // avatar url
	Nickname     string          `json:"nickname"`      // nickname
	Introduction string          `json:"introduction"`  // introduction
	Status       uint64          `json:"status"`        // status(0: disabled, 1: enable)
	LastLogin    carbon.DateTime `json:"lastLogin"`     // last login time
	Locked       uint64          `json:"locked"`        // locked(0: unlock, 1: locked)
	LockExpire   int64           `json:"lockExpire"`    // lock expiration time
	Wrong        int64           `json:"wrong"`         // wrong password count
}

// NewUserRepo .
func NewUserRepo(data *Data, action biz.ActionRepo) biz.UserRepo {
	return &userRepo{
		data:   data,
		action: action,
	}
}

func (ro userRepo) GetByUsername(ctx context.Context, username string) (item *biz.User, err error) {
	item = &biz.User{}
	var m User
	ro.data.DB(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
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
	m.Id = ro.data.Id(ctx)
	m.UserCode = id.NewCode(m.Id)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	err = db.Create(&m).Error
	return
}

func (ro userRepo) LastLogin(ctx context.Context, id uint64) (err error) {
	fields := make(map[string]interface{})
	fields["wrong"] = constant.I0
	fields["last_login"] = carbon.Now()
	fields["locked"] = constant.UI0
	fields["lock_expire"] = constant.I0
	err = ro.data.DB(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Updates(&fields).Error
	return
}

func (ro userRepo) WrongPwd(ctx context.Context, req biz.LoginTime) (err error) {
	oldItem, err := ro.GetByUsername(ctx, req.Username)
	if err != nil {
		return
	}
	if oldItem.LastLogin.Gt(req.LastLogin.Carbon) {
		// already login success, skip set wrong count
		return
	}
	m := make(map[string]interface{})
	newWrong := oldItem.Wrong + 1
	if newWrong >= 5 {
		m["locked"] = constant.UI1
		if newWrong == 5 {
			m["lock_expire"] = carbon.Now().AddDuration("5m").Carbon2Time().Unix()
		} else if newWrong == 10 {
			m["lock_expire"] = carbon.Now().AddDuration("30m").Carbon2Time().Unix()
		} else if newWrong == 20 {
			m["lock_expire"] = carbon.Now().AddDuration("24h").Carbon2Time().Unix()
		} else if newWrong >= 30 {
			// forever lock
			m["lock_expire"] = 0
		}
	}
	m["wrong"] = newWrong
	err = ro.data.DB(ctx).
		Model(&User{}).
		Where("id = ?", oldItem.Id).
		Where("wrong = ?", oldItem.Wrong).
		Updates(&m).Error
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
	fields := make(map[string]interface{})
	fields["password"] = item.Password
	fields["wrong"] = constant.I0
	fields["locked"] = constant.UI0
	fields["lock_expire"] = constant.I0
	err = db.
		Model(&User{}).
		Where("id = ?", m.Id).
		Updates(&fields).Error
	return
}

func (ro userRepo) IdExists(ctx context.Context, id uint64) (err error) {
	var m User
	ro.data.DB(ctx).
		Where("id = ?", id).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.UserNotFound
		return
	}
	return
}
