package biz

import (
	"context"
	"strings"

	"auth/internal/conf"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
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

type FindRoleCache struct {
	Page page.Page `json:"page"`
	List []Role    `json:"list"`
}

type UpdateRole struct {
	Id     uint64  `json:"id,string"`
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
	return &RoleUseCase{
		c:     c,
		repo:  repo,
		tx:    tx,
		cache: cache.WithPrefix("role"),
	}
}

func (uc *RoleUseCase) Create(ctx context.Context, item *Role) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *RoleUseCase) Find(ctx context.Context, condition *FindRole) (rp []Role) {
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, ok := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.find(ctx, action, condition)
	})
	if ok {
		var cache FindRoleCache
		utils.Json2Struct(&cache, str)
		condition.Page = cache.Page
		rp = cache.List
	}
	return
}

func (uc *RoleUseCase) find(ctx context.Context, action string, condition *FindRole) (res string, ok bool) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindRoleCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	ok = true
	return
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
