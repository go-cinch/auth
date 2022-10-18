package biz

import (
	"auth/internal/conf"
	"context"
)

type Permission struct {
	Resources []string `json:"resources"`
	Menus     []string `json:"menus"`
	Btns      []string `json:"btns"`
}

type CheckPermission struct {
	UserCode string `json:"userCode"`
	Resource string `json:"resource"`
}

type PermissionRepo interface {
	Check(ctx context.Context, item *CheckPermission) bool
	GetByUserCode(ctx context.Context, code string) (*Permission, error)
}

type PermissionUseCase struct {
	c     *conf.Bootstrap
	repo  PermissionRepo
	cache Cache
}

func NewPermissionUseCase(c *conf.Bootstrap, repo PermissionRepo, cache Cache) *PermissionUseCase {
	return &PermissionUseCase{c: c, repo: repo, cache: cache.WithPrefix("auth_permission")}
}

func (uc *PermissionUseCase) Check(ctx context.Context, item *CheckPermission) bool {
	return uc.repo.Check(ctx, item)
}

func (uc *PermissionUseCase) GetByUserCode(ctx context.Context, code string) (*Permission, error) {
	return uc.repo.GetByUserCode(ctx, code)
}
