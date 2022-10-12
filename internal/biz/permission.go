package biz

import (
	"auth/internal/conf"
	"context"
)

type Permission struct {
	UserCode string `json:"userCode"`
	Resource string `json:"resource"`
}

type PermissionRepo interface {
	Check(ctx context.Context, item *Permission) bool
}

type PermissionUseCase struct {
	c     *conf.Bootstrap
	repo  PermissionRepo
	cache Cache
}

func NewPermissionUseCase(c *conf.Bootstrap, repo PermissionRepo, cache Cache) *PermissionUseCase {
	cache.Register("auth_permission_cache")
	return &PermissionUseCase{c: c, repo: repo, cache: cache}
}

func (uc *PermissionUseCase) Check(ctx context.Context, item *Permission) bool {
	return uc.repo.Check(ctx, item)
}
