package biz

import (
	"auth/internal/conf"
	"context"
	"github.com/go-cinch/common/page"
)

type Action struct {
	Id       uint64 `json:"id,string"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Word     string `json:"word"`
	Resource string `json:"resource"`
}

type FindAction struct {
	Page     page.Page `json:"page"`
	Code     *string   `json:"code"`
	Name     *string   `json:"name"`
	Word     *string   `json:"word"`
	Resource *string   `json:"resource"`
}

type ActionRepo interface {
	Create(ctx context.Context, item *Action) error
	Find(ctx context.Context, condition *FindAction) ([]Action, error)
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
	cache.Register("auth_action_cache")
	return &ActionUseCase{c: c, repo: repo, tx: tx, cache: cache}
}

func (uc *ActionUseCase) Create(ctx context.Context, item *Action) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *ActionUseCase) Find(ctx context.Context, condition *FindAction) ([]Action, error) {
	return uc.repo.Find(ctx, condition)
}
