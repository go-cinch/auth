package data

import (
	"context"
	"regexp"
	"strings"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type whitelistRepo struct {
	data *Data
}

func NewWhitelistRepo(data *Data) biz.WhitelistRepo {
	return &whitelistRepo{
		data: data,
	}
}

func (ro whitelistRepo) Create(ctx context.Context, item *biz.Whitelist) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	db := gorm.G[model.Whitelist](ro.data.DB(ctx))

	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}

	var m model.Whitelist
	copierx.Copy(&m, item)

	// Set values directly since model fields are non-pointer.
	m.Category = item.Category
	m.Resource = item.Resource

	err = db.Create(ctx, &m)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("create whitelist failed")
	}
	return err
}

func (ro whitelistRepo) Update(ctx context.Context, item *biz.UpdateWhitelist) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	db := gorm.G[model.Whitelist](ro.data.DB(ctx))

	m, err := db.Where("id = ?", item.ID).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return biz.ErrRecordNotFound(ctx)
	}
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("get whitelist failed")
		return err
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		return biz.ErrDataNotChange(ctx)
	}

	// Note: Use native DB.Updates for map updates, gorm.G.Updates expects struct type.
	err = ro.data.DB(ctx).
		Model(&model.Whitelist{}).
		Where("id = ?", item.ID).
		Updates(change).
		Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update whitelist failed")
	}
	return err
}

func (ro whitelistRepo) Delete(ctx context.Context, ids ...int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	if len(ids) == 0 {
		return nil
	}
	db := gorm.G[model.Whitelist](ro.data.DB(ctx))

	_, err = db.Where("id IN ?", ids).Delete(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("delete whitelist failed")
	}
	return err
}

func (ro whitelistRepo) Find(ctx context.Context, condition *biz.FindWhitelist) (rp []biz.Whitelist) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	rp = make([]biz.Whitelist, 0)
	q := gorm.G[model.Whitelist](ro.data.DB(ctx)).Where("1 = 1")

	// Apply filters.
	if condition.Category != nil {
		q = q.Where("category = ?", *condition.Category)
	}
	if condition.Resource != nil {
		q = q.Where("resource LIKE ?", "%"+*condition.Resource+"%")
	}

	// Count total before pagination.
	if !condition.Page.Disable {
		count, err := q.Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("count whitelist failed")
			return rp
		}
		condition.Page.Total = count
		if count == 0 {
			return rp
		}
	}

	// Apply ordering + pagination.
	q = q.Order("id DESC")
	if !condition.Page.Disable {
		limit, offset := condition.Page.Limit()
		q = q.Limit(int(limit)).Offset(int(offset))
	}

	list, err := q.Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find whitelist failed")
		return rp
	}

	copierx.Copy(&rp, list)
	return rp
}

func (ro whitelistRepo) Match(ctx context.Context, category int16, resource string) (ok bool, err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Match")
	defer span.End()

	resource = strings.TrimSpace(resource)
	if resource == "" {
		return false, nil
	}

	db := gorm.G[model.Whitelist](ro.data.DB(ctx))
	list, err := db.Select("resource").Where("category = ?", category).Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find whitelist failed")
		return false, err
	}

	req := parseMatchResource(resource)
	for _, item := range list {
		if item.Resource == "" {
			continue
		}
		if matchWhitelistResource(item.Resource, req) {
			return true, nil
		}
	}
	return false, nil
}

type matchReq struct {
	Resource string
	Method   string
	URI      string
}

func parseMatchResource(s string) matchReq {
	parts := strings.Split(s, "|")
	switch len(parts) {
	case 1:
		return matchReq{Resource: strings.TrimSpace(parts[0])}
	case 2:
		return matchReq{
			Method: strings.TrimSpace(parts[0]),
			URI:    strings.TrimSpace(parts[1]),
		}
	default:
		return matchReq{
			Method:   strings.TrimSpace(parts[0]),
			URI:      strings.TrimSpace(parts[1]),
			Resource: strings.TrimSpace(strings.Join(parts[2:], "|")),
		}
	}
}

func matchWhitelistResource(patterns string, req matchReq) bool {
	for _, line := range strings.Split(patterns, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if matchWhitelistLine(line, req) {
			return true
		}
	}
	return false
}

func matchWhitelistLine(line string, req matchReq) bool {
	if line == "*" {
		return true
	}
	parts := strings.Split(line, "|")
	switch len(parts) {
	case 1:
		want := strings.TrimSpace(parts[0])
		return want != "" && req.Resource != "" && req.Resource == want
	case 2:
		return matchHTTP(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), req)
	default:
		methods := strings.TrimSpace(parts[0])
		uriPattern := strings.TrimSpace(parts[1])
		grpcRes := strings.TrimSpace(strings.Join(parts[2:], "|"))
		if grpcRes != "" && req.Resource != "" && req.Resource == grpcRes {
			return true
		}
		return matchHTTP(methods, uriPattern, req)
	}
}

func matchHTTP(methodSpec, uriPattern string, req matchReq) bool {
	if req.Method == "" || req.URI == "" || methodSpec == "" || uriPattern == "" {
		return false
	}
	okMethod := false
	for _, m := range strings.Split(methodSpec, ",") {
		if strings.TrimSpace(m) == req.Method {
			okMethod = true
			break
		}
	}
	if !okMethod {
		return false
	}
	return matchGlob(uriPattern, req.URI)
}

func matchGlob(pattern, s string) bool {
	if pattern == "*" {
		return true
	}
	// Fast path: exact match when no wildcards are present.
	if !strings.ContainsAny(pattern, "*?") {
		return pattern == s
	}
	// Convert a tiny glob subset to regex:
	// - '*' => '.*'
	// - '?' => '.'
	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*`, `.*`)
	re = strings.ReplaceAll(re, `\?`, `.`)
	re = "^" + re + "$"
	return regexp.MustCompile(re).MatchString(s)
}
