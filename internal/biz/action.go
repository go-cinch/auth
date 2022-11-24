package biz

import (
	"auth/internal/conf"
	"context"
	"fmt"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
)

type Action struct {
	Id       uint64 `json:"id,string"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Word     string `json:"word"`
	Resource string `json:"resource"`
	Menu     string `json:"menu"`
	Btn      string `json:"btn"`
}

type FindAction struct {
	Page     page.Page `json:"page"`
	Code     *string   `json:"code"`
	Name     *string   `json:"name"`
	Word     *string   `json:"word"`
	Resource *string   `json:"resource"`
}

type FindActionCache struct {
	Page page.Page `json:"page"`
	List []Action  `json:"list"`
}

type UpdateAction struct {
	Id       *uint64 `json:"id,string,omitempty"`
	Name     *string `json:"name,omitempty"`
	Word     *string `json:"word,omitempty"`
	Resource *string `json:"resource,omitempty"`
	Menu     *string `json:"menu,omitempty"`
	Btn      *string `json:"btn,omitempty"`
}

type ActionRepo interface {
	Create(ctx context.Context, item *Action) error
	Find(ctx context.Context, condition *FindAction) []Action
	FindByCode(ctx context.Context, code string) []Action
	Update(ctx context.Context, item *UpdateAction) error
	Delete(ctx context.Context, ids ...uint64) error
	CodeExists(ctx context.Context, code string) error
	Permission(ctx context.Context, code, resource string) bool
}

type ActionUseCase struct {
	c     *conf.Bootstrap
	repo  ActionRepo
	tx    Transaction
	cache Cache
}

func NewActionUseCase(c *conf.Bootstrap, repo ActionRepo, tx Transaction, cache Cache) *ActionUseCase {
	return &ActionUseCase{c: c, repo: repo, tx: tx, cache: cache.WithPrefix("auth_action")}
}

func (uc *ActionUseCase) Create(ctx context.Context, item *Action) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *ActionUseCase) Find(ctx context.Context, condition *FindAction) (rp []Action) {
	action := fmt.Sprintf("find_%s", utils.StructMd5(condition))
	str, ok := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.find(ctx, action, condition)
	})
	if ok {
		var cache FindActionCache
		utils.Json2Struct(&cache, str)
		condition.Page = cache.Page
		rp = cache.List
	}
	return
}

func (uc *ActionUseCase) find(ctx context.Context, action string, condition *FindAction) (res string, ok bool) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindActionCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	ok = true
	return
}

func (uc *ActionUseCase) Update(ctx context.Context, item *UpdateAction) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *ActionUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
