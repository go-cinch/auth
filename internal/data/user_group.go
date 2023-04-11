package data

import (
	"auth/api/reason"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/utils"
	"strings"
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
	Word   string `json:"word"`                                           // keyword, must be unique, used as frontend display
	Action string `json:"action"`                                         // user group action code array
}

type UserUserGroupRelation struct {
	UserId      uint64 `json:"userId,string"`
	UserGroupId uint64 `json:"userGroupId,string"`
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
		Where("`word` = ?", item.Word).
		First(&m)
	if m.Id > constant.UI0 {
		err = reason.ErrorIllegalParameter("%s `word`: %s", i18n.FromContext(ctx).T(biz.DuplicateField), item.Word)
		return
	}
	copierx.Copy(&m, item)
	m.Id = ro.data.Id(ctx)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	if len(item.Users) > 0 {
		m.Users = make([]User, 0)
		for _, v := range item.Users {
			err = ro.user.IdExists(ctx, v.Id)
			if err != nil {
				return
			}
			m.Users = append(m.Users, User{
				Id: v.Id,
			})
		}
	}
	err = db.Create(&m).Error
	return
}

func (ro userGroupRepo) FindGroupByUserCode(ctx context.Context, code string) (list []biz.UserGroup) {
	list = make([]biz.UserGroup, 0)
	user, err := ro.user.GetByCode(ctx, code)
	if err != nil {
		return
	}
	db := ro.data.DB(ctx)
	groupIds := make([]uint64, 0)
	db.
		Model(&UserUserGroupRelation{}).
		Where("`user_id` = ?", user.Id).
		Pluck("`user_group_id`", &groupIds)
	if len(groupIds) == 0 {
		return
	}
	groups := make([]UserGroup, 0)
	db.
		Model(&UserGroup{}).
		Where("`id` IN (?)", groupIds).
		Find(&groups)
	copierx.Copy(&list, groups)
	return
}

func (ro userGroupRepo) Find(ctx context.Context, condition *biz.FindUserGroup) (rp []biz.UserGroup) {
	db := ro.data.DB(ctx)
	db = db.
		Model(&UserGroup{}).
		Preload("Users").
		Order("id DESC")
	rp = make([]biz.UserGroup, 0)
	list := make([]UserGroup, 0)
	if condition.Name != nil {
		db.Where("`name` LIKE ?", strings.Join([]string{"%", *condition.Name, "%"}, ""))
	}
	if condition.Code != nil {
		db.Where("`code` LIKE ?", strings.Join([]string{"%", *condition.Code, "%"}, ""))
	}
	if condition.Word != nil {
		db.Where("`word` LIKE ?", strings.Join([]string{"%", *condition.Word, "%"}, ""))
	}
	if condition.Action != nil {
		db.Where("`action` LIKE ?", strings.Join([]string{"%", *condition.Action, "%"}, ""))
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

func (ro userGroupRepo) Update(ctx context.Context, item *biz.UpdateUserGroup) (err error) {
	var m UserGroup
	db := ro.data.DB(ctx)
	db.
		Where("`id` = ?", item.Id).
		First(&m)
	if m.Id == constant.UI0 {
		err = reason.ErrorNotFound("%s UserGroup.id: %d", i18n.FromContext(ctx).T(biz.RecordNotFound), item.Id)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.DataNotChange))
		return
	}
	if item.Word != nil && *item.Word != m.Word {
		err = ro.WordExists(ctx, *item.Word)
		if err == nil {
			err = reason.ErrorIllegalParameter("%s `word`: %s", i18n.FromContext(ctx).T(biz.DuplicateField), *item.Word)
			return
		}
	}
	if a, ok1 := change["users"]; ok1 {
		if v, ok2 := a.(string); ok2 {
			arr := utils.Str2Uint64Arr(v)
			users := make([]User, 0)
			for _, id := range arr {
				users = append(users, User{
					Id: id,
				})
			}
			err = db.
				Model(&m).
				Association("Users").
				Replace(users)
			if err != nil {
				return
			}
			delete(change, "users")
		}
	}
	err = db.
		Model(&m).
		Updates(&change).Error
	return
}

func (ro userGroupRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	db := ro.data.DB(ctx)
	err = db.
		Where("`id` IN (?)", ids).
		Delete(&UserGroup{}).Error
	return
}

func (ro userGroupRepo) WordExists(ctx context.Context, word string) (err error) {
	var m UserGroup
	db := ro.data.DB(ctx)
	arr := strings.Split(word, ",")
	for _, item := range arr {
		db.
			Where("`word` = ?", item).
			First(&m)
		if m.Id == constant.UI0 {
			err = reason.ErrorNotFound("%s UserGroup.code: %s", i18n.FromContext(ctx).T(biz.RecordNotFound), item)
			return
		}
	}
	return
}
