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

type userGroupRepo struct {
	data   *Data
	action biz.ActionRepo
	user   biz.UserRepo
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
	p := query.Use(ro.data.DB(ctx)).UserGroup
	db := p.WithContext(ctx)
	m := db.GetByCol(p.Word.ColumnName().String(), item.Word)
	if m.ID > constant.UI0 {
		err = biz.ErrDuplicateField(ctx, p.Word.ColumnName().String(), item.Word)
		return
	}
	copierx.Copy(&m, item)
	m.ID = ro.data.ID(ctx)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	if len(item.Users) > 0 {
		m.Users = make([]model.User, 0)
		for _, v := range item.Users {
			err = ro.user.IdExists(ctx, v.Id)
			if err != nil {
				return
			}
			m.Users = append(m.Users, model.User{
				ID: v.Id,
			})
		}
	}
	err = db.Create(&m)
	return
}

func (ro userGroupRepo) Find(ctx context.Context, condition *biz.FindUserGroup) (rp []biz.UserGroup) {
	p := query.Use(ro.data.DB(ctx)).UserGroup
	db := p.WithContext(ctx)
	rp = make([]biz.UserGroup, 0)
	list := make([]model.UserGroup, 0)
	conditions := make([]gen.Condition, 0, 2)
	if condition.Name != nil {
		conditions = append(conditions, p.Name.Like(strings.Join([]string{"%", *condition.Name, "%"}, "")))
	}
	if condition.Word != nil {
		conditions = append(conditions, p.Word.Like(strings.Join([]string{"%", *condition.Word, "%"}, "")))
	}
	if condition.Action != nil {
		conditions = append(conditions, p.Action.Like(strings.Join([]string{"%", *condition.Action, "%"}, "")))
	}
	condition.Page.Primary = p.ID.ColumnName().String()
	condition.Page.
		WithContext(ctx).
		Query(
			db.
				Preload(p.Users).
				Order(p.ID.Desc()).
				Where(conditions...).
				UnderlyingDB().
				Model(&model.UserGroup{}),
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

func (ro userGroupRepo) Update(ctx context.Context, item *biz.UpdateUserGroup) (err error) {
	p := query.Use(ro.data.DB(ctx)).UserGroup
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
	if item.Word != nil && *item.Word != m.Word {
		ok := ro.WordExists(ctx, *item.Word)
		if ok {
			err = biz.ErrDuplicateField(ctx, p.Word.ColumnName().String(), *item.Word)
			return
		}
	}
	if a, ok1 := change["users"]; ok1 {
		if v, ok2 := a.(string); ok2 {
			arr := utils.Str2Uint64Arr(v)
			users := make([]*model.User, 0)
			for _, id := range arr {
				users = append(users, &model.User{
					ID: id,
				})
			}
			err = p.Users.
				Model(&m).
				Replace(users...)
			if err != nil {
				return
			}
			delete(change, "users")
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.Id)).
		Updates(&change)
	return
}

func (ro userGroupRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).UserGroup
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro userGroupRepo) WordExists(ctx context.Context, word string) (ok bool) {
	p := query.Use(ro.data.DB(ctx)).UserGroup
	db := p.WithContext(ctx)
	arr := strings.Split(word, ",")
	for _, item := range arr {
		m := db.GetByCol(p.Word.ColumnName().String(), item)
		if m.ID == constant.UI0 {
			log.
				WithContext(ctx).
				Error("invalid word: %s", item)
			return
		}
	}
	ok = true
	return
}
