package data

import (
	"context"
	"strings"

	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data/model"
	"auth/internal/data/query"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/gobwas/glob"
	"gorm.io/gen"
)

type actionRepo struct {
	c    *conf.Bootstrap
	data *Data
}

func NewActionRepo(c *conf.Bootstrap, data *Data) biz.ActionRepo {
	return &actionRepo{
		c:    c,
		data: data,
	}
}

func (ro actionRepo) Create(ctx context.Context, item *biz.Action) (err error) {
	ok := ro.WordExists(ctx, item.Word)
	if ok {
		err = biz.ErrDuplicateField(ctx, "word", item.Word)
		return
	}
	var m model.Action
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	m.Code = id.NewCode(m.ID)
	if m.Resource == "" {
		m.Resource = "*"
	}
	err = db.Create(&m)
	return
}

func (ro actionRepo) GetDefault(ctx context.Context) (rp biz.Action) {
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	m := db.GetByCol("word", "default")
	copierx.Copy(&rp, m)
	return
}

func (ro actionRepo) Find(ctx context.Context, condition *biz.FindAction) (rp []biz.Action) {
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	rp = make([]biz.Action, 0)
	list := make([]model.Action, 0)
	conditions := make([]gen.Condition, 0, 2)
	if condition.Name != nil {
		conditions = append(conditions, p.Name.Like(strings.Join([]string{"%", *condition.Name, "%"}, "")))
	}
	if condition.Code != nil {
		conditions = append(conditions, p.Code.Like(strings.Join([]string{"%", *condition.Code, "%"}, "")))
	}
	if condition.Word != nil {
		conditions = append(conditions, p.Word.Like(strings.Join([]string{"%", *condition.Word, "%"}, "")))
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

func (ro actionRepo) FindByCode(ctx context.Context, code string) (rp []biz.Action) {
	rp = make([]biz.Action, 0)
	if code == "" {
		return
	}
	list := make([]*model.Action, 0)
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	arr := strings.Split(code, ",")
	list, _ = db.
		Where(p.Code.In(arr...)).
		Find()
	copierx.Copy(&rp, list)
	return
}

func (ro actionRepo) Update(ctx context.Context, item *biz.UpdateAction) (err error) {
	p := query.Use(ro.data.DB(ctx)).Action
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
			err = biz.ErrDuplicateField(ctx, "word", *item.Word)
			return
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.Id)).
		Updates(&change)
	return
}

func (ro actionRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	if err != nil {
		return
	}
	count, _ := db.Count()
	if count == 0 {
		err = biz.ErrKeepLeastOneAction(ctx)
	}
	return
}

func (ro actionRepo) CodeExists(ctx context.Context, code string) (err error) {
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	arr := strings.Split(code, ",")
	for _, item := range arr {
		m := db.GetByCol("code", item)
		if m.ID == constant.UI0 {
			err = biz.ErrRecordNotFound(ctx)
			log.
				WithError(err).
				Error("invalid code: %s", code)
			return
		}
	}
	return
}

func (ro actionRepo) WordExists(ctx context.Context, word string) (ok bool) {
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	arr := strings.Split(word, ",")
	for _, item := range arr {
		m := db.GetByCol("word", item)
		if m.ID == constant.UI0 {
			log.Error("invalid word: %s", item)
			return
		}
	}
	ok = true
	return
}

func (ro actionRepo) Permission(ctx context.Context, code string, req biz.CheckPermission) (pass bool) {
	arr := strings.Split(code, ",")
	for _, item := range arr {
		pass = ro.permission(ctx, item, req)
		if pass {
			return
		}
	}
	return
}

func (ro actionRepo) permission(ctx context.Context, code string, req biz.CheckPermission) (pass bool) {
	if code == "" {
		return
	}
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	m := db.GetByCol("code", code)
	if m.ID == constant.UI0 {
		return
	}
	return ro.MatchResource(ctx, m.Resource, req)
}

func (actionRepo) MatchResource(_ context.Context, resource string, req biz.CheckPermission) (pass bool) {
	if resource == "" {
		// empty resource no need match
		return
	}
	arr1 := strings.Split(resource, "\n")
	for _, v1 := range arr1 {
		if v1 == "*" {
			pass = true
			return
		}
		arr2 := strings.Split(v1, "|")
		switch len(arr2) {
		case 1:
			// only grpc resource
			if req.Resource == arr2[0] {
				pass = true
				return
			}
		case 2:
			// only http method / http uri
			methods := strings.Split(arr2[0], ",")
			g, err := glob.Compile(arr2[1])
			if err != nil {
				return
			}
			matched := g.Match(req.URI)
			if matched && utils.Contains[string](methods, req.Method) {
				pass = true
				return
			}
		case 3:
			// grpc resource / http method / http uri
			// match one means has permission
			if req.Resource == arr2[2] {
				pass = true
				return
			}
			methods := strings.Split(arr2[0], ",")
			g, err := glob.Compile(arr2[1])
			if err != nil {
				return
			}
			matched := g.Match(req.URI)
			if matched && utils.Contains[string](methods, req.Method) {
				pass = true
				return
			}
		}
	}
	return
}
