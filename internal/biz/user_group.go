package biz

import (
	"auth/internal/conf"
	"context"
)

type UserGroup struct {
	Id     uint64   `json:"id,string"`
	Users  []uint64 `json:"users"`
	Name   string   `json:"name"`
	Word   string   `json:"word"`
	Action string   `json:"action"`
}

type UserGroupRepo interface {
	Create(ctx context.Context, item *UserGroup) error
	FindGroupByUserCode(ctx context.Context, code string) ([]UserGroup, error)
}

type UserGroupUseCase struct {
	c     *conf.Bootstrap
	repo  UserGroupRepo
	tx    Transaction
	cache Cache
}

func NewUserGroupUseCase(c *conf.Bootstrap, repo UserGroupRepo, tx Transaction, cache Cache) *UserGroupUseCase {
	cache.Register("auth_user_group_cache")
	return &UserGroupUseCase{c: c, repo: repo, tx: tx, cache: cache}
}

func (uc *UserGroupUseCase) Create(ctx context.Context, item *UserGroup) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}
