package biz

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"auth/api/reason"
	"auth/internal/conf"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/utils"
)

type Permission struct {
	Resources []string `json:"resources"`
	Menus     []string `json:"menus"`
	Btns      []string `json:"btns"`
}

type CheckPermission struct {
	UserCode string `json:"userCode"`
	Resource string `json:"resource"`
	Method   string `json:"method"`
	URI      string `json:"uri"`
}

type PermissionRepo interface {
	Check(ctx context.Context, item CheckPermission) bool
	GetByUserCode(ctx context.Context, code string) (*Permission, error)
}

type PermissionUseCase struct {
	c     *conf.Bootstrap
	repo  PermissionRepo
	cache Cache
}

func NewPermissionUseCase(c *conf.Bootstrap, repo PermissionRepo, cache Cache) *PermissionUseCase {
	return &PermissionUseCase{
		c:     c,
		repo:  repo,
		cache: cache.WithPrefix("permission"),
	}
}

func (uc *PermissionUseCase) Check(ctx context.Context, item CheckPermission) (rp bool) {
	action := strings.Join([]string{"check", utils.StructMd5(item)}, "_")
	str, ok := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.check(ctx, action, item)
	})
	if ok {
		rp, _ = strconv.ParseBool(str)
	}
	return
}

func (uc *PermissionUseCase) check(ctx context.Context, action string, item CheckPermission) (res string, ok bool) {
	// read data from db and write to cache
	pass := uc.repo.Check(ctx, item)
	res = strconv.FormatBool(pass)
	uc.cache.Set(ctx, action, res, false)
	ok = true
	return
}

func (uc *PermissionUseCase) GetByUserCode(ctx context.Context, code string) (rp *Permission, err error) {
	rp = &Permission{}
	action := strings.Join([]string{"get_by_user_code", code}, "_")
	str, ok := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.getByUserCode(ctx, action, code)
	})
	if ok {
		utils.Json2Struct(&rp, str)
		return
	}
	err = reason.ErrorTooManyRequests(i18n.FromContext(ctx).T(TooManyRequests))
	return
}

func (uc *PermissionUseCase) getByUserCode(ctx context.Context, action string, code string) (res string, ok bool) {
	// read data from db and write to cache
	rp := &Permission{}
	permission, err := uc.repo.GetByUserCode(ctx, code)
	notFound := errors.Is(err, reason.ErrorNotFound(i18n.FromContext(ctx).T(RecordNotFound)))
	if err != nil && !notFound {
		return
	}
	copierx.Copy(&rp, permission)
	res = utils.Struct2Json(rp)
	uc.cache.Set(ctx, action, res, notFound)
	ok = true
	return
}

func (uc *PermissionUseCase) FlushCache(ctx context.Context) {
	uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
		return
	})
}
