package biz

import (
	"context"
	"errors"
	"strings"

	"auth/internal/conf"
	"github.com/go-cinch/common/copierx"
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
	rp = uc.repo.Check(ctx, item)
	return
}

func (uc *PermissionUseCase) GetByUserCode(ctx context.Context, code string) (rp *Permission, err error) {
	rp = &Permission{}
	action := strings.Join([]string{"get_by_user_code", code}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.getByUserCode(ctx, action, code)
	})
	if err != nil {
		return
	}
	utils.Json2Struct(&rp, str)
	return
}

func (uc *PermissionUseCase) getByUserCode(ctx context.Context, action string, code string) (res string, err error) {
	// read data from db and write to cache
	rp := &Permission{}
	permission, err := uc.repo.GetByUserCode(ctx, code)
	notFound := errors.Is(err, ErrRecordNotFound(ctx))
	if err != nil && !notFound {
		return
	}
	copierx.Copy(&rp, permission)
	res = utils.Struct2Json(rp)
	uc.cache.Set(ctx, action, res, notFound)
	return
}

func (uc *PermissionUseCase) FlushCache(ctx context.Context) {
	uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
		return
	})
}
