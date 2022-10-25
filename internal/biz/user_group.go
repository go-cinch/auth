package biz

import (
	"auth/internal/conf"
	"context"
	"fmt"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
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

type FindUserGroupCache struct {
	Page page.Page   `json:"page"`
	List []UserGroup `json:"list"`
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
	return &UserGroupUseCase{c: c, repo: repo, tx: tx, cache: cache.WithPrefix("auth_user_group")}
}

func (uc *UserGroupUseCase) Create(ctx context.Context, item *UserGroup) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *UserGroupUseCase) Find(ctx context.Context, condition *FindUserGroup) (rp []UserGroup) {
	action := fmt.Sprintf("find_%s", utils.StructMd5(condition))
	str, ok, _, _ := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.find(ctx, action, condition)
	})
	if ok {
		var cache FindUserGroupCache
		utils.Json2Struct(&cache, str)
		condition.Page = cache.Page
		rp = cache.List
	}
	return
}

func (uc *UserGroupUseCase) find(ctx context.Context, action string, condition *FindUserGroup) (res string, ok bool) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindUserGroupCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	ok = true
	return
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
