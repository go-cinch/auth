package data

import (
	"context"
	"strings"

	"auth/api/reason"
	"auth/internal/biz"
	"auth/internal/data/model"
	"auth/internal/data/query"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/utils"
	"gorm.io/gen"
)

type whitelistRepo struct {
	data   *Data
	action biz.ActionRepo
}

func NewWhitelistRepo(data *Data, action biz.ActionRepo) biz.WhitelistRepo {
	return &whitelistRepo{
		data:   data,
		action: action,
	}
}

func (ro whitelistRepo) Create(ctx context.Context, item *biz.Whitelist) (err error) {
	var m model.Whitelist
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).Whitelist
	db := p.WithContext(ctx)
	m.ID = ro.data.Id(ctx)
	err = db.Create(&m)
	return
}

func (ro whitelistRepo) Find(ctx context.Context, condition *biz.FindWhitelist) (rp []biz.Whitelist) {
	p := query.Use(ro.data.DB(ctx)).Whitelist
	db := p.WithContext(ctx)
	rp = make([]biz.Whitelist, 0)
	list := make([]model.Whitelist, 0)
	conditions := make([]gen.Condition, 0, 2)
	if condition.Category != nil {
		conditions = append(conditions, p.Category.Eq(*condition.Category))
	}
	if condition.Resource != nil {
		conditions = append(conditions, p.Resource.Like(strings.Join([]string{"%", *condition.Resource, "%"}, "")))
	}
	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(
			db.
				Order(p.ID.Desc()).
				Where(conditions...).
				UnderlyingDB(),
		).
		Find(&list)
	copierx.Copy(&rp, list)
	return
}

func (ro whitelistRepo) Has(ctx context.Context, condition *biz.HasWhitelist) (has bool) {
	p := query.Use(ro.data.DB(ctx)).Whitelist
	db := p.WithContext(ctx)
	resources := make([]string, 0)
	db.
		Where(p.Category.Eq(condition.Category)).
		Pluck(p.Resource, &resources)
	for _, item := range resources {
		has = ro.action.MatchResource(ctx, item, condition.Permission)
		if has {
			return
		}
	}
	return
}

func (ro whitelistRepo) Update(ctx context.Context, item *biz.UpdateWhitelist) (err error) {
	p := query.Use(ro.data.DB(ctx)).Whitelist
	db := p.WithContext(ctx)
	m := db.GetByID(item.Id)
	if m.ID == constant.UI0 {
		err = reason.ErrorNotFound("%s Whitelist.id: %d", i18n.FromContext(ctx).T(biz.RecordNotFound), item.Id)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.DataNotChange))
		return
	}
	_, err = db.
		Where(p.ID.Eq(item.Id)).
		Updates(&change)
	return
}

func (ro whitelistRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).Whitelist
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}
