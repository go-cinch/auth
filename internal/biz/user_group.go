package biz

import (
	"auth/internal/conf"
	"context"
	"github.com/go-cinch/common/page"
)

type UserGroup struct {
	Id      uint64   `json:"id,string"`
	Users   []User   `json:"users"`
	Name    string   `json:"name"`
	Word    string   `json:"word"`
	Action  string   `json:"action"`
	Actions []Action `json:"actions"`
}

type FindUserGroup struct {
	Page   page.Page `json:"page"`
	Code   *string   `json:"code"`
	Name   *string   `json:"name"`
	Word   *string   `json:"word"`
	Action *string   `json:"action"`
}

type UpdateUserGroup struct {
	Id     *uint64 `json:"id,string,omitempty"`
	Name   *string `json:"name,omitempty"`
	Word   *string `json:"word,omitempty"`
	Action *string `json:"action,omitempty"`
	Users  *string `json:"users,omitempty"`
}

type UserGroupRepo interface {
	Create(ctx context.Context, item *UserGroup) error
	Find(ctx context.Context, condition *FindUserGroup) []UserGroup
	Update(ctx context.Context, item *UpdateUserGroup) error
	Delete(ctx context.Context, ids ...uint64) error
	FindGroupByUserCode(ctx context.Context, code string) []UserGroup
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

func (uc *UserGroupUseCase) Find(ctx context.Context, condition *FindUserGroup) []UserGroup {
	return uc.repo.Find(ctx, condition)
}

func (uc *UserGroupUseCase) Update(ctx context.Context, item *UpdateUserGroup) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *UserGroupUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
