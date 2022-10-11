package biz

import (
	"auth/internal/conf"
	"context"
)

type Role struct {
	Id     uint64 `json:"id,string"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Status uint64 `json:"status"`
}

type UpdateRole struct {
	Id     *uint64 `json:"id,string,omitempty"`
	Name   *string `json:"name,omitempty"`
	Key    *string `json:"key,omitempty"`
	Status *uint64 `json:"status,omitempty"`
}

type RoleRepo interface {
	Create(ctx context.Context, item *Role) error
	Update(ctx context.Context, item *UpdateRole) error
}

type RoleUseCase struct {
	c     *conf.Bootstrap
	repo  RoleRepo
	tx    Transaction
	cache Cache
}

func NewRoleUseCase(c *conf.Bootstrap, repo RoleRepo, tx Transaction, cache Cache) *RoleUseCase {
	cache.Register("auth_role_cache")
	return &RoleUseCase{c: c, repo: repo, tx: tx, cache: cache}
}

func (uc *RoleUseCase) Create(ctx context.Context, item *Role) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *RoleUseCase) Update(ctx context.Context, item *UpdateRole) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, item)
		})
	})
}
