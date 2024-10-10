package data

import (
	"context"
	"strconv"
	"strings"

	"auth/internal/biz"
	"auth/internal/data/model"
	"auth/internal/data/query"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

type userRepo struct {
	data   *Data
	action biz.ActionRepo
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
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	m, err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(p.Username.Eq(username)).
		First()
	if err != nil || m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		log.
			WithContext(ctx).
			WithError(err).
			Error("invalid username: %s", username)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro userRepo) Find(ctx context.Context, condition *biz.FindUser) (rp []biz.User) {
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	conditions := make([]gen.Condition, 0, 2)
	rp = make([]biz.User, 0)
	list := make([]model.User, 0)
	if condition.StartCreatedAt != nil {
		conditions = append(conditions, p.CreatedAt.Gte(carbon.Parse(*condition.StartCreatedAt)))
	}
	if condition.EndCreatedAt != nil {
		conditions = append(conditions, p.CreatedAt.Lt(carbon.Parse(*condition.EndCreatedAt)))
	}
	if condition.StartUpdatedAt != nil {
		conditions = append(conditions, p.CreatedAt.Gte(carbon.Parse(*condition.StartUpdatedAt)))
	}
	if condition.EndUpdatedAt != nil {
		conditions = append(conditions, p.CreatedAt.Lt(carbon.Parse(*condition.EndUpdatedAt)))
	}
	if condition.Username != nil {
		conditions = append(conditions, p.Username.Like(strings.Join([]string{"%", *condition.Username, "%"}, "")))
	}
	if condition.Code != nil {
		conditions = append(conditions, p.Code.Like(strings.Join([]string{"%", *condition.Code, "%"}, "")))
	}
	if condition.Platform != nil {
		conditions = append(conditions, p.Platform.Like(strings.Join([]string{"%", *condition.Platform, "%"}, "")))
	}
	if condition.Locked != nil {
		conditions = append(conditions, p.Locked.Is(*condition.Locked))
	}
	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(
			db.
				Preload(p.Role).
				Order(p.CreatedAt.Desc()).
				Where(conditions...).
				UnderlyingDB().
				Model(&model.User{}),
		).
		Find(&list)
	copierx.Copy(&rp, list)
	timestamp := carbon.Now().Timestamp()
	for i, item := range rp {
		rp[i].Actions = make([]biz.Action, 0)
		arr := ro.action.FindByCode(ctx, item.Action)
		copierx.Copy(&rp[i].Actions, arr)
		if !item.Locked || (item.LockExpire > constant.I0 && timestamp > item.LockExpire) {
			rp[i].Locked = false
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
		ms := make([]string, 0)
		if hours < 24 {
			if hours > 0 {
				ms = append(ms, strconv.FormatInt(hours, 10), "h")
			}
			if minutes > 0 {
				ms = append(ms, strconv.FormatInt(minutes, 10), "m")
			}
			if seconds > 0 {
				ms = append(ms, strconv.FormatInt(seconds, 10), "s")
			}
		} else {
			ms = append(ms, carbon.CreateFromTimestamp(item.LockExpire).ToDateTimeString())
		}
		rp[i].LockMsg = strings.Join(ms, "")
	}
	return
}

func (ro userRepo) Create(ctx context.Context, item *biz.User) (err error) {
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	m := db.GetByCol("username", item.Username)
	if m.ID > constant.UI0 {
		err = biz.ErrDuplicateField(ctx, "username", item.Username)
		return
	}
	copierx.Copy(&m, item)
	m.ID = ro.data.ID(ctx)
	m.Code = id.NewCode(m.ID)
	if m.Action != "" {
		err = ro.action.CodeExists(ctx, m.Action)
		if err != nil {
			return
		}
	}
	err = db.Create(&m)
	return
}

func (ro userRepo) Update(ctx context.Context, item *biz.UpdateUser) (err error) {
	p := query.Use(ro.data.DB(ctx)).User
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
	// check lock or unlock
	if locked, ok1 := change["locked"]; ok1 {
		if v1, ok2 := locked.(uint64); ok2 {
			var lockExpire int64
			if expire, ok3 := change["lock_expire"]; ok3 {
				if v2, ok4 := expire.(int64); ok4 {
					lockExpire = v2
				}
			}
			if m.Locked && v1 == constant.UI0 {
				change["lock_expire"] = constant.I0
			} else if !m.Locked && v1 == constant.UI1 {
				change["lock_expire"] = lockExpire
			}
		}
	}
	if username, ok1 := change["username"]; ok1 {
		if v, ok2 := username.(string); ok2 {
			_, err = ro.GetByUsername(ctx, v)
			if err == nil {
				err = biz.ErrDuplicateField(ctx, "username", v)
				return
			}
		}
	}
	if roleId, ok1 := change["role_id"]; ok1 {
		if v, ok2 := roleId.(string); ok2 && v != "0" {
			pRole := query.Use(ro.data.DB(ctx)).Role
			dbRole := pRole.WithContext(ctx)
			mRole := dbRole.GetByID(*item.RoleId)
			if mRole.ID == constant.UI0 {
				err = biz.ErrRecordNotFound(ctx)
				log.
					WithContext(ctx).
					WithError(err).
					Error("invalid roleId: %s", v)
				return
			}
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.Id)).
		Updates(&change)
	return
}

func (ro userRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro userRepo) LastLogin(ctx context.Context, username string) (err error) {
	fields := make(map[string]interface{})
	fields["wrong"] = constant.I0
	fields["last_login"] = carbon.Now()
	fields["locked"] = constant.UI0
	fields["lock_expire"] = constant.I0
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.Username.Eq(username)).
		Updates(&fields)
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
	change := make(map[string]interface{})
	newWrong := oldItem.Wrong + 1
	if req.Wrong > 0 {
		newWrong = req.Wrong
	}
	if newWrong >= 5 {
		change["locked"] = constant.UI1
		if newWrong == 5 {
			change["lock_expire"] = carbon.Now().AddDuration("5m").StdTime().Unix()
		} else if newWrong == 10 {
			change["lock_expire"] = carbon.Now().AddDuration("30m").StdTime().Unix()
		} else if newWrong == 20 {
			change["lock_expire"] = carbon.Now().AddDuration("24h").StdTime().Unix()
		} else if newWrong >= 30 {
			// forever lock
			change["lock_expire"] = 0
		}
	}
	change["wrong"] = newWrong

	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.Eq(oldItem.Id)).
		Where(p.Wrong.Eq(oldItem.Wrong)).
		Updates(&change)
	return
}

func (ro userRepo) UpdatePassword(ctx context.Context, item *biz.User) (err error) {
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	m := db.GetByCol("username", item.Username)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	fields := make(map[string]interface{})
	fields["password"] = item.Password
	fields["wrong"] = constant.I0
	fields["locked"] = constant.UI0
	fields["lock_expire"] = constant.I0
	_, err = db.
		Where(p.ID.Eq(m.ID)).
		Updates(&fields)
	return
}

func (ro userRepo) IdExists(ctx context.Context, id uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	return
}

func (ro userRepo) GetByCode(ctx context.Context, code string) (item *biz.User, err error) {
	item = &biz.User{}
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	m, err := db.
		Preload(p.Role).
		Where(p.Code.Eq(code)).
		First()
	if err != nil || m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		log.
			WithContext(ctx).
			WithError(err).
			Error("invalid code: %s", code)
		return
	}
	copierx.Copy(&item, m)
	return
}
