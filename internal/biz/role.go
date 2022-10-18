package biz

import (
	"auth/internal/conf"
	"context"
	"github.com/go-cinch/common/page"
)

type Role struct {
	Id      uint64   `json:"id,string"`
	Name    string   `json:"name"`
	Word    string   `json:"word"`
	Action  string   `json:"action"`
	Actions []Action `json:"actions"`
}

type FindRole struct {
	Page   page.Page `json:"page"`
	Name   *string   `json:"name"`
	Word   *string   `json:"word"`
	Action *string   `json:"action"`
}

type UpdateRole struct {
	Id     *uint64 `json:"id,string,omitempty"`
	Name   *string `json:"name,omitempty"`
	Word   *string `json:"word,omitempty"`
	Action *string `json:"action,omitempty"`
}

type RoleRepo interface {
	Create(ctx context.Context, item *Role) error
	Find(ctx context.Context, condition *FindRole) []Role
	Update(ctx context.Context, item *UpdateRole) error
	Delete(ctx context.Context, ids ...uint64) error
}

type RoleUseCase struct {
	c     *conf.Bootstrap
	repo  RoleRepo
	tx    Transaction
	cache Cache
}

func NewRoleUseCase(c *conf.Bootstrap, repo RoleRepo, tx Transaction, cache Cache) *RoleUseCase {
	return &RoleUseCase{c: c, repo: repo, tx: tx, cache: cache.WithPrefix("auth_role")}
}

func (uc *RoleUseCase) Create(ctx context.Context, item *Role) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *RoleUseCase) Find(ctx context.Context, condition *FindRole) []Role {
	return uc.repo.Find(ctx, condition)
}

func (uc *RoleUseCase) Update(ctx context.Context, item *UpdateRole) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *RoleUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
