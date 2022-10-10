package biz

import (
	"auth/internal/conf"
	"context"
)

type Action struct {
	Id   uint64 `json:"id,string"`
	Code string `json:"code"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Path string `json:"path"`
}

type ActionRepo interface {
	Create(ctx context.Context, item *Action) error
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