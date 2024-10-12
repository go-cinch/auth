package biz

import (
	"context"

	"auth/internal/conf"
)

type HotspotRepo interface {
	Refresh(ctx context.Context) error
	GetUserByCode(ctx context.Context, code string) *User
	GetUserByUsername(ctx context.Context, username string) *User
	GetRoleByID(ctx context.Context, id uint64) *Role
	GetActionByWord(ctx context.Context, word string) *Action
	GetActionByCode(ctx context.Context, code string) *Action
	FindActionByCode(ctx context.Context, codes ...string) []Action
	FindWhitelistResourceByCategory(ctx context.Context, category uint32) []string
	FindUserGroupByUserCode(ctx context.Context, code string) []UserGroup
}

type HotspotUseCase struct {
	c    *conf.Bootstrap
	repo HotspotRepo
}

func NewHotspotUseCase(c *conf.Bootstrap, repo HotspotRepo) *HotspotUseCase {
	return &HotspotUseCase{
		c:    c,
		repo: repo,
	}
}

func (uc *HotspotUseCase) Refresh(ctx context.Context) error {
	return uc.repo.Refresh(ctx)
}
