package biz

import (
	"context"

	"auth/internal/conf"
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
	GetByUserCode(ctx context.Context, code string) *Permission
}

type PermissionUseCase struct {
	c    *conf.Bootstrap
	repo PermissionRepo
}

func NewPermissionUseCase(c *conf.Bootstrap, repo PermissionRepo) *PermissionUseCase {
	return &PermissionUseCase{
		c:    c,
		repo: repo,
	}
}

func (uc *PermissionUseCase) Check(ctx context.Context, item CheckPermission) (rp bool) {
	rp = uc.repo.Check(ctx, item)
	return
}

func (uc *PermissionUseCase) GetByUserCode(ctx context.Context, code string) (rp *Permission) {
	rp = uc.repo.GetByUserCode(ctx, code)
	return
}
