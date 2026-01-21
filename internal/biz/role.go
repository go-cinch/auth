package biz

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/page/v2"
	"github.com/go-cinch/common/utils"

	"auth/internal/conf"
)

type Role struct {
	ID      int64    `json:"id,string"`
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
	ID     int64   `json:"id,string"`
	Name   *string `json:"name,omitempty"`
	Word   *string `json:"word,omitempty"`
	Action *string `json:"action,omitempty"`
}

type RoleUseCase struct {
	c     *conf.Bootstrap
	repo  RoleRepo
	tx    Transaction
	cache Cache
}

func NewRoleUseCase(c *conf.Bootstrap, repo RoleRepo, tx Transaction, cache Cache) *RoleUseCase {
	return &RoleUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "role",
		}, "_")),
	}
}

func (uc *RoleUseCase) Create(ctx context.Context, item *Role) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *RoleUseCase) Find(ctx context.Context, condition *FindRole) (rp []Role, err error) {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return nil, err
	}
	var cache FindRoleCache
	utils.JSON2Struct(&cache, str)
	condition.Page = cache.Page
	return cache.List, nil
}

func (uc *RoleUseCase) find(ctx context.Context, action string, condition *FindRole) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindRoleCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2JSON(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return res, nil
}

func (uc *RoleUseCase) Update(ctx context.Context, item *UpdateRole) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *RoleUseCase) Delete(ctx context.Context, ids ...int64) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Delete(ctx, ids...)
		})
	})
}
