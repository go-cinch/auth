package data

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/golang-module/carbon/v2"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type userRepo struct {
	data *Data
}

// NewUserRepo creates a new User repository.
func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}

func (ro userRepo) GetByUsername(ctx context.Context, username string) (item *biz.User, err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetByUsername")
	defer span.End()

	item = &biz.User{}

	db := gorm.G[model.User](ro.data.DB(ctx), clause.Locking{Strength: "UPDATE"})
	m, qErr := db.Where("username = ?", username).First(ctx)
	if qErr != nil || m.ID == 0 {
		// Keep old behavior: treat any error as "record not found".
		err = biz.ErrRecordNotFound(ctx)
		log.WithContext(ctx).WithError(err).Error("invalid username: %s", username)
		return
	}

	if err = copierx.Copy(item, m); err != nil {
		log.WithContext(ctx).WithError(err).Error("copy user failed")
		return
	}
	item.Id = m.ID
	item.Locked = m.Locked != nil && *m.Locked != 0
	return
}

func (ro userRepo) Find(ctx context.Context, condition *biz.FindUser) (rp []biz.User) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	rp = make([]biz.User, 0)

	q := gorm.G[model.User](ro.data.DB(ctx)).Where("1 = 1")

	if condition.StartCreatedAt != nil {
		q = q.Where("created_at >= ?", carbon.Parse(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where("created_at < ?", carbon.Parse(*condition.EndCreatedAt))
	}
	// Keep old behavior: "updated" filters were mistakenly applied to created_at.
	if condition.StartUpdatedAt != nil {
		q = q.Where("created_at >= ?", carbon.Parse(*condition.StartUpdatedAt))
	}
	if condition.EndUpdatedAt != nil {
		q = q.Where("created_at < ?", carbon.Parse(*condition.EndUpdatedAt))
	}
	if condition.Username != nil {
		q = q.Where("username LIKE ?", "%"+*condition.Username+"%")
	}
	if condition.Code != nil {
		q = q.Where("code LIKE ?", "%"+*condition.Code+"%")
	}
	if condition.Platform != nil {
		q = q.Where("platform LIKE ?", "%"+*condition.Platform+"%")
	}
	if condition.Locked != nil {
		q = q.Where("locked = ?", boolToInt16(*condition.Locked))
	}

	// Count total before pagination.
	if !condition.Page.Disable {
		count, err := q.Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("count user failed")
			return rp
		}
		condition.Page.Total = count
		if count == 0 {
			return rp
		}
	}

	// Apply ordering + pagination.
	q = q.Order("created_at DESC")
	if !condition.Page.Disable {
		limit, offset := condition.Page.Limit()
		q = q.Limit(limit).Offset(offset)
	}

	list, err := q.Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find user failed")
		return rp
	}

	// Batch load roles since generated model doesn't include association fields.
	roleMap := make(map[int64]model.Role)
	roleIDs := make([]int64, 0, len(list))
	seenRole := make(map[int64]struct{}, len(list))
	for _, u := range list {
		if u.RoleID == nil || *u.RoleID == 0 {
			continue
		}
		rid := *u.RoleID
		if _, ok := seenRole[rid]; ok {
			continue
		}
		seenRole[rid] = struct{}{}
		roleIDs = append(roleIDs, rid)
	}
	if len(roleIDs) > 0 {
		rq := gorm.G[model.Role](ro.data.DB(ctx)).Where("id IN ?", roleIDs)
		roles, rErr := rq.Find(ctx)
		if rErr != nil {
			log.WithContext(ctx).WithError(rErr).Error("find role for users failed")
			return rp
		}
		for _, r := range roles {
			roleMap[r.ID] = r
		}
	}
	timestamp := carbon.Now().Timestamp()
	rp = make([]biz.User, 0, len(list))
	for _, m := range list {
		var u biz.User
		if err := copierx.Copy(&u, m); err != nil {
			log.WithContext(ctx).WithError(err).Error("copy user failed")
			return rp
		}
		u.Id = m.ID
		u.Locked = m.Locked != nil && *m.Locked != 0

		// Fill role info if requested/available.
		if m.RoleID != nil {
			if r, ok := roleMap[*m.RoleID]; ok {
				copierx.Copy(&u.Role, r)
				u.Role.ID = r.ID
			}
		}

		// Populate actions (keep old behavior).
		u.Actions, _ = ro.getActionsByCode(ctx, u.Action)
		// Lock status/message (keep old behavior).
		if !u.Locked || (u.LockExpire > constant.I0 && timestamp > u.LockExpire) {
			u.Locked = false
			rp = append(rp, u)
			continue
		}
		if u.LockExpire == constant.I0 {
			u.LockMsg = "forever"
			rp = append(rp, u)
			continue
		}
		diff := u.LockExpire - timestamp
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
			ms = append(ms, carbon.CreateFromTimestamp(u.LockExpire).ToDateTimeString())
		}
		u.LockMsg = strings.Join(ms, "")

		rp = append(rp, u)
	}

	return rp
}

func (ro userRepo) Create(ctx context.Context, item *biz.User) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	db := gorm.G[model.User](ro.data.DB(ctx))

	count, err := db.Where("username = ?", item.Username).Count(ctx, "*")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("check username exists failed")
		return err
	}
	if count > 0 {
		err = biz.ErrDuplicateField(ctx, "username", item.Username)
		return err
	}

	var m model.User
	if err = copierx.Copy(&m, item); err != nil {
		log.WithContext(ctx).WithError(err).Error("copy user failed")
		return err
	}
	m.ID = ro.data.ID(ctx)
	m.Code = id.NewCode(uint64(m.ID))
	m.Locked = ptrInt16(boolToInt16(item.Locked))

	if m.Action != nil && *m.Action != "" {
		if err = ro.validateActionCodes(ctx, *m.Action); err != nil {
			return err
		}
	}

	err = db.Create(ctx, &m)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("create user failed")
	}
	return err
}

func (ro userRepo) Update(ctx context.Context, item *biz.UpdateUser) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	db := gorm.G[model.User](ro.data.DB(ctx))

	m, qErr := db.Where("id = ?", item.Id).First(ctx)
	if qErr == gorm.ErrRecordNotFound || m.ID == 0 {
		err = biz.ErrRecordNotFound(ctx)
		return err
	}
	if qErr != nil {
		log.WithContext(ctx).WithError(qErr).Error("get user failed")
		err = qErr
		return err
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.ErrDataNotChange(ctx)
		return err
	}

	// check lock or unlock
	if locked, ok := change["locked"]; ok {
		if v, ok := toInt16(locked); ok {
			var lockExpire int64
			if expire, ok := change["lock_expire"]; ok {
				if v2, ok := toInt64(expire); ok {
					lockExpire = v2
				}
			}

			oldLocked := m.Locked != nil && *m.Locked != 0
			newLocked := v != 0
			if oldLocked && !newLocked {
				change["lock_expire"] = constant.I0
			} else if !oldLocked && newLocked {
				change["lock_expire"] = lockExpire
			}
			change["locked"] = v
		}
	}

	if username, ok := change["username"]; ok {
		if v, ok := username.(string); ok {
			_, e := ro.GetByUsername(ctx, v)
			if e == nil {
				err = biz.ErrDuplicateField(ctx, "username", v)
				return err
			}
		}
	}

	if _, ok := change["role_id"]; ok {
		roleID := int64(0)
		if item.RoleId != nil {
			roleID = *item.RoleId
		} else if v, ok := toInt64(change["role_id"]); ok {
			roleID = v
		}
		if roleID != 0 {
			rdb := gorm.G[model.Role](ro.data.DB(ctx))
			rm, rErr := rdb.Where("id = ?", roleID).First(ctx)
			if rErr == gorm.ErrRecordNotFound || rm.ID == 0 {
				err = biz.ErrRecordNotFound(ctx)
				log.WithContext(ctx).WithError(err).Error("invalid roleId: %d", roleID)
				return err
			}
			if rErr != nil {
				log.WithContext(ctx).WithError(rErr).Error("get role failed")
				err = rErr
				return err
			}
		}
		// Make sure role_id is a numeric value for Updates(map).
		if v, ok := toInt64(change["role_id"]); ok {
			change["role_id"] = v
		}
	}

	// Normalize numeric fields that may come from JSON maps (float64/string).
	if v, ok := change["lock_expire"]; ok {
		if vv, ok := toInt64(v); ok {
			change["lock_expire"] = vv
		}
	}
	if v, ok := change["wrong"]; ok {
		if vv, ok := toInt64(v); ok {
			change["wrong"] = vv
		}
	}

	err = ro.data.DB(ctx).Model(&model.User{}).Where("id = ?", item.Id).Updates(change).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update user failed")
	}
	return err
}

func (ro userRepo) Delete(ctx context.Context, ids ...int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	db := gorm.G[model.User](ro.data.DB(ctx))
	_, err = db.Where("id IN ?", ids).Delete(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("delete user failed")
	}
	return
}

func (ro userRepo) LastLogin(ctx context.Context, username string) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "LastLogin")
	defer span.End()

	fields := map[string]interface{}{
		"last_login":  carbon.Now(),
		"wrong":       constant.I0,
		"locked":      int16(0),
		"lock_expire": constant.I0,
	}
	err = ro.data.DB(ctx).Model(&model.User{}).Where("username = ?", username).Updates(fields).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update last login failed")
	}
	return
}

func (ro userRepo) WrongPwd(ctx context.Context, req *biz.LoginTime) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "WrongPwd")
	defer span.End()
	oldItem, err := ro.GetByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if oldItem.LastLogin.Gt(req.LastLogin.Carbon) {
		// already login success, skip set wrong count
		return err
	}

	change := make(map[string]interface{})
	newWrong := oldItem.Wrong + 1
	if req.Wrong > 0 {
		newWrong = req.Wrong
	}
	if newWrong >= 5 {
		change["locked"] = int16(1)
		if newWrong == 5 {
			change["lock_expire"] = carbon.Now().AddDuration("5m").StdTime().Unix()
		} else if newWrong == 10 {
			change["lock_expire"] = carbon.Now().AddDuration("30m").StdTime().Unix()
		} else if newWrong == 20 {
			change["lock_expire"] = carbon.Now().AddDuration("24h").StdTime().Unix()
		} else if newWrong >= 30 {
			// forever lock
			change["lock_expire"] = int64(0)
		}
	}
	change["wrong"] = newWrong

	err = ro.data.DB(ctx).
		Model(&model.User{}).
		Where("id = ? AND wrong = ?", oldItem.Id, oldItem.Wrong).
		Updates(change).
		Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update wrong password failed")
	}
	return err
}

func (ro userRepo) UpdatePassword(ctx context.Context, item *biz.User) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "UpdatePassword")
	defer span.End()

	db := gorm.G[model.User](ro.data.DB(ctx))

	m, qErr := db.Where("username = ?", item.Username).First(ctx)
	if qErr == gorm.ErrRecordNotFound || m.ID == 0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	if qErr != nil {
		log.WithContext(ctx).WithError(qErr).Error("get user failed")
		err = qErr
		return
	}

	fields := map[string]interface{}{
		"password":    item.Password,
		"wrong":       constant.I0,
		"locked":      int16(0),
		"lock_expire": constant.I0,
	}
	err = ro.data.DB(ctx).Model(&model.User{}).Where("id = ?", m.ID).Updates(fields).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update password failed")
	}
	return
}

func (ro userRepo) IdExists(ctx context.Context, id int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "IdExists")
	defer span.End()

	db := gorm.G[model.User](ro.data.DB(ctx))

	m, qErr := db.Where("id = ?", id).First(ctx)
	if qErr == gorm.ErrRecordNotFound || m.ID == 0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	if qErr != nil {
		log.WithContext(ctx).WithError(qErr).Error("get user failed")
		err = qErr
		return
	}
	return
}

func (ro userRepo) validateActionCodes(ctx context.Context, codes string) (err error) {
	codes = strings.TrimSpace(codes)
	if codes == "" {
		return nil
	}
	arr := ro.splitCodes(codes)
	if len(arr) == 0 {
		return nil
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	for _, code := range arr {
		count, err := db.Where("code = ?", code).Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("check action code exists failed")
			return err
		}
		if count == 0 {
			err = biz.ErrRecordNotFound(ctx)
			log.WithContext(ctx).WithError(err).Error("invalid code: %s", code)
			return err
		}
	}
	return nil
}

func (ro userRepo) getActionsByCode(ctx context.Context, code string) (rp []biz.Action, err error) {
	rp = make([]biz.Action, 0)
	code = strings.TrimSpace(code)
	if code == "" {
		return rp, nil
	}
	codes := ro.splitCodes(code)
	if len(codes) == 0 {
		return rp, nil
	}

	list, err := gorm.G[model.Action](ro.data.DB(ctx)).
		Where("code IN ?", codes).
		Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find actions by code failed")
		return nil, err
	}
	copierx.Copy(&rp, list)
	return rp, nil
}

func (userRepo) splitCodes(code string) []string {
	arr := strings.Split(code, ",")
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

func boolToInt16(v bool) int16 {
	if v {
		return 1
	}
	return 0
}

func ptrInt16(v int16) *int16 { return &v }

func toInt16(v interface{}) (int16, bool) {
	switch x := v.(type) {
	case int16:
		return x, true
	case int:
		return int16(x), true
	case int64:
		return int16(x), true
	case float64:
		return int16(x), true
	default:
		return 0, false
	}
}

func toInt64(v interface{}) (int64, bool) {
	switch x := v.(type) {
	case int64:
		return x, true
	case int:
		return int64(x), true
	case uint64:
		return int64(x), true
	case float64:
		return int64(x), true
	case float32:
		return int64(x), true
	case string:
		if x == "" {
			return 0, true
		}
		i, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

func toBool(v interface{}) (bool, bool) {
	switch x := v.(type) {
	case bool:
		return x, true
	case int:
		return x != 0, true
	case int64:
		return x != 0, true
	case uint64:
		return x != 0, true
	case float64:
		return x != 0, true
	case string:
		if x == "" {
			return false, true
		}
		b, err := strconv.ParseBool(x)
		if err != nil {
			return false, false
		}
		return b, true
	default:
		return false, false
	}
}
