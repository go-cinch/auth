package data

import (
	"context"
	"strings"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) biz.RoleRepo {
	return &roleRepo{
		data: data,
	}
}

func (ro roleRepo) Create(ctx context.Context, item *biz.Role) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	// Check word uniqueness
	count, err := gorm.G[model.Role](ro.data.DB(ctx)).
		Where("word = ?", item.Word).
		Count(ctx, "*")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("check role word exists failed")
		return err
	}
	if count > 0 {
		return biz.ErrDuplicateField(ctx, "word", item.Word)
	}

	if err = ro.validateActionCodes(ctx, item.Action); err != nil {
		return err
	}

	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}

	var m model.Role
	copierx.Copy(&m, item)
	err = gorm.G[model.Role](ro.data.DB(ctx)).Create(ctx, &m)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("create role failed")
	}
	return err
}

func (ro roleRepo) Update(ctx context.Context, item *biz.UpdateRole) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	db := gorm.G[model.Role](ro.data.DB(ctx))

	m, err := db.Where("id = ?", item.ID).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return biz.ErrRecordNotFound(ctx)
	}
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("get role failed")
		return err
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		return biz.ErrDataNotChange(ctx)
	}

	// Check word uniqueness if word is being updated
	if item.Word != nil && (m.Word == nil || *item.Word != *m.Word) {
		count, err := db.Where("word = ? AND id != ?", *item.Word, item.ID).Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("check role word uniqueness failed")
			return err
		}
		if count > 0 {
			return biz.ErrDuplicateField(ctx, "word", *item.Word)
		}
	}

	// Validate action codes if action is being updated
	if item.Action != nil && (m.Action == nil || *item.Action != *m.Action) {
		if err = ro.validateActionCodes(ctx, *item.Action); err != nil {
			return err
		}
	}

	// Use native DB.Updates for map updates.
	err = ro.data.DB(ctx).
		Model(&model.Role{}).
		Where("id = ?", item.ID).
		Updates(change).
		Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update role failed")
	}
	return err
}

func (ro roleRepo) Delete(ctx context.Context, ids ...int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	_, err = gorm.G[model.Role](ro.data.DB(ctx)).
		Where("id IN ?", ids).
		Delete(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("delete role failed")
	}
	return err
}

func (ro roleRepo) Find(ctx context.Context, condition *biz.FindRole) (rp []biz.Role) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	rp = make([]biz.Role, 0)

	q := gorm.G[model.Role](ro.data.DB(ctx)).Where("1 = 1")
	if condition.Name != nil {
		q = q.Where("name LIKE ?", "%"+*condition.Name+"%")
	}
	if condition.Word != nil {
		q = q.Where("word LIKE ?", "%"+*condition.Word+"%")
	}
	if condition.Action != nil {
		q = q.Where("action LIKE ?", "%"+*condition.Action+"%")
	}

	// Count total before pagination
	if !condition.Page.Disable {
		count, err := q.Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("count role failed")
			return rp
		}
		condition.Page.Total = count
		if count == 0 {
			return rp
		}
	}

	q = q.Order("id DESC")
	if !condition.Page.Disable {
		limit, offset := condition.Page.Limit()
		q = q.Limit(int(limit)).Offset(int(offset))
	}

	list, err := q.Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find role failed")
		return rp
	}
	copierx.Copy(&rp, list)

	for i := range rp {
		actions, err := ro.getActionsByCode(ctx, rp[i].Action)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("get role actions failed")
			continue
		}
		rp[i].Actions = actions
	}
	return rp
}

func (ro roleRepo) validateActionCodes(ctx context.Context, codes string) (err error) {
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

func (ro roleRepo) getActionsByCode(ctx context.Context, code string) (rp []biz.Action, err error) {
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

func (roleRepo) splitCodes(code string) []string {
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
