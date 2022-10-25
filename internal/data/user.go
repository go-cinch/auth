package data

import (
	"auth/internal/biz"
	"context"
	"fmt"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/utils"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm/clause"
)

type userRepo struct {
	data   *Data
	action biz.ActionRepo
}

// User is database fields map
type User struct {
	Id         uint64          `json:"id,string"`                     // auto increment id
	CreatedAt  carbon.DateTime `json:"createdAt"`                     // create time
	UpdatedAt  carbon.DateTime `json:"updatedAt"`                     // update time
	RoleId     uint64          `json:"roleId,string"`                 // role id
	Role       Role            `json:"role" gorm:"foreignKey:RoleId"` // role
	Action     string          `json:"action"`                        // user action code array
	Username   string          `json:"username"`                      // user login name
	Code       string          `json:"code"`                          // user code
	Password   string          `json:"password"`                      // password
	Platform   string          `json:"platform"`                      // device platform: pc/android/ios/mini...
	LastLogin  carbon.DateTime `json:"lastLogin"`                     // last login time
	Locked     uint64          `json:"locked"`                        // locked(0: unlock, 1: locked)
	LockExpire int64           `json:"lockExpire"`                    // lock expiration time
	Wrong      int64           `json:"wrong"`                         // wrong password count
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
		Where("`username` = ?", username).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.NotFound("%s User.username: %s", biz.RecordNotFound.Message, username)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro userRepo) Find(ctx context.Context, condition *biz.FindUser) (rp []biz.User) {
	db := ro.data.DB(ctx)
	db = db.
		Model(&User{}).
		Order("created_at DESC")
	rp = make([]biz.User, 0)
	list := make([]User, 0)
	if condition.StartCreatedAt != nil {
		db.Where("`created_at` >= ?", condition.StartCreatedAt)
	}
	if condition.EndCreatedAt != nil {
		db.Where("`created_at` < ?", condition.EndCreatedAt)
	}
	if condition.StartUpdatedAt != nil {
		db.Where("`updated_at` >= ?", condition.StartUpdatedAt)
	}
	if condition.EndUpdatedAt != nil {
		db.Where("`updated_at` < ?", condition.EndCreatedAt)
	}
	if condition.Username != nil {
		db.Where("`username` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Username))
	}
	if condition.Code != nil {
		db.Where("`code` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Code))
	}
	if condition.Platform != nil {
		db.Where("`platform` LIKE ?", fmt.Sprintf("%%%s%%", *condition.Platform))
	}
	if condition.Locked != nil {
		db.Where("`locked` = ?", *condition.Locked)
	}
	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(db).
		Find(&list)
	copierx.Copy(&rp, list)
	timestamp := carbon.Now().Timestamp()
	for i, item := range rp {
		rp[i].Actions = make([]biz.Action, 0)
		arr := ro.action.FindByCode(ctx, item.Action)
		copierx.Copy(&rp[i].Actions, arr)
		if item.Locked == constant.UI0 || (item.LockExpire > constant.I0 && timestamp > item.LockExpire) {
			rp[i].Locked = constant.UI0
			continue
		}
		if item.LockExpire == constant.I0 {
			rp[i].LockMsg = "forever"
			continue
		}
		diff := item.LockExpire - timestamp
		hours := diff / 3600
		minutes := diff % 3600 / 60
		seconds := diff % 3600 % 60
		msg := ""
		if hours < 24 {
			if hours > 0 {
				msg += fmt.Sprintf("%dh", hours)
			}
			if minutes > 0 {
				msg += fmt.Sprintf("%dm", minutes)
			}
			if seconds > 0 {
				msg += fmt.Sprintf("%ds", seconds)
			}
		} else {
			msg = carbon.CreateFromTimestamp(item.LockExpire).ToDateTimeString()
		}
		rp[i].LockMsg = msg
	}
	return
}

func (ro userRepo) Create(ctx context.Context, item *biz.User) (err error) {
	var m User
	db := ro.data.DB(ctx)
	db.
		Where("`username` = ?", item.Username).
		First(&m)
	if m.Id > constant.UI0 {
		err = biz.IllegalParameter("%s `username`: %s", biz.DuplicateField.Message, item.Username)
		return
	}
	copierx.Copy(&m, item)
	m.Id = ro.data.Id(ctx)
	m.Code = id.NewCode(m.Id)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	err = db.Create(&m).Error
	return
}

func (ro userRepo) Update(ctx context.Context, item *biz.UpdateUser) (err error) {
	var m User
	db := ro.data.DB(ctx)
	db.
		Where("`id` = ?", item.Id).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.NotFound("%s User.id: %d", biz.RecordNotFound.Message, item.Id)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.DataNotChange
		return
	}
	// check lock or unlock
	if locked, ok1 := change["locked"]; ok1 {
		if v1, ok2 := locked.(uint64); ok2 {
			var lockExpire int64
			if expire, ok3 := change["lock_expire"]; ok3 {
				if v2, ok4 := expire.(int64); ok4 {
					lockExpire = v2
				}
			}
			if m.Locked == constant.UI1 && v1 == constant.UI0 {
				change["lock_expire"] = constant.I0
			} else if m.Locked == constant.UI0 && v1 == constant.UI1 {
				change["lock_expire"] = lockExpire
			}
		}
	}
	if username, ok1 := change["username"]; ok1 {
		if v, ok2 := username.(string); ok2 {
			_, err = ro.GetByUsername(ctx, v)
			if err == nil {
				err = biz.IllegalParameter("%s `username`: %s", biz.DuplicateField.Message, v)
				return
			}
		}
	}
	err = db.
		Model(&m).
		Updates(&change).Error
	return
}

func (ro userRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	db := ro.data.DB(ctx)
	err = db.
		Where("`id` IN (?)", ids).
		Delete(&User{}).Error
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
		Where("`id` = ?", id).
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
		Where("`id` = ?", oldItem.Id).
		Where("`wrong` = ?", oldItem.Wrong).
		Updates(&m).Error
	return
}

func (ro userRepo) UpdatePassword(ctx context.Context, item *biz.User) (err error) {
	var m User
	db := ro.data.DB(ctx)
	db.
		Where("`username` = ?", item.Username).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.NotFound("%s User.username: %s", biz.RecordNotFound.Message, item.Username)
		return
	}
	fields := make(map[string]interface{})
	fields["password"] = item.Password
	fields["wrong"] = constant.I0
	fields["locked"] = constant.UI0
	fields["lock_expire"] = constant.I0
	err = db.
		Model(&User{}).
		Where("`id` = ?", m.Id).
		Updates(&fields).Error
	return
}

func (ro userRepo) IdExists(ctx context.Context, id uint64) (err error) {
	var m User
	ro.data.DB(ctx).
		Where("`id` = ?", id).
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.NotFound("%s User.id: %d", biz.RecordNotFound.Message, id)
		return
	}
	return
}

func (ro userRepo) GetByCode(ctx context.Context, code string) (item *biz.User, err error) {
	item = &biz.User{}
	var m User
	ro.data.DB(ctx).
		Where("`code` = ?", code).
		Preload("Role").
		First(&m)
	if m.Id == constant.UI0 {
		err = biz.NotFound("%s User.code: %s", biz.RecordNotFound.Message, code)
		return
	}
	copierx.Copy(&item, m)
	return
}
