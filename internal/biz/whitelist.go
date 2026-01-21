package biz

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/page/v2"
	"github.com/go-cinch/common/utils"

	"auth/internal/conf"
)

const (
	// WhitelistPermissionCategory is used by the permission middleware to skip permission checks.
	WhitelistPermissionCategory int16 = iota
	// WhitelistJwtCategory is used by the auth middleware to skip JWT validation/parsing.
	WhitelistJwtCategory
)

type Whitelist struct {
	ID       int64  `json:"id,string"`
	Category int16  `json:"category"`
	Resource string `json:"resource"`
}

type FindWhitelist struct {
	Page     page.Page `json:"page"`
	Category *int16    `json:"category"`
	Resource *string   `json:"resource"`
}

type FindWhitelistCache struct {
	Page page.Page   `json:"page"`
	List []Whitelist `json:"list"`
}

type UpdateWhitelist struct {
	ID       int64   `json:"id,string"`
	Category *int16  `json:"category,omitempty"`
	Resource *string `json:"resource,omitempty"`
}

type WhitelistUseCase struct {
	c     *conf.Bootstrap
	repo  WhitelistRepo
	tx    Transaction
	cache Cache
}

func NewWhitelistUseCase(c *conf.Bootstrap, repo WhitelistRepo, tx Transaction, cache Cache) *WhitelistUseCase {
	return &WhitelistUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "whitelist",
		}, "_")),
	}
}

func (uc *WhitelistUseCase) Create(ctx context.Context, item *Whitelist) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *WhitelistUseCase) Find(ctx context.Context, condition *FindWhitelist) (rp []Whitelist, err error) {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return nil, err
	}

	var cache FindWhitelistCache
	utils.JSON2Struct(&cache, str)
	condition.Page = cache.Page
	return cache.List, nil
}

func (uc *WhitelistUseCase) find(ctx context.Context, action string, condition *FindWhitelist) (res string, err error) {
	list := uc.repo.Find(ctx, condition)
	cache := FindWhitelistCache{
		List: list,
		Page: condition.Page,
	}
	res = utils.Struct2JSON(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return res, nil
}

func (uc *WhitelistUseCase) Match(ctx context.Context, category int16, resource string) (ok bool, err error) {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Match")
	defer span.End()

	resource = strings.TrimSpace(resource)
	if resource == "" {
		return false, nil
	}

	// Match sits on the hot path (middleware); cache results briefly.
	type matchKey struct {
		Category int16  `json:"category"`
		Resource string `json:"resource"`
	}
	action := strings.Join([]string{
		"match",
		utils.StructMd5(matchKey{Category: category, Resource: resource}),
	}, "_")

	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		ok, err := uc.repo.Match(ctx, category, resource)
		if err != nil {
			return "", err
		}
		res := utils.Struct2JSON(ok)
		uc.cache.Set(ctx, action, res, !ok)
		return res, nil
	})
	if err != nil {
		return false, err
	}

	// Cache stores a JSON boolean: true/false.
	utils.JSON2Struct(&ok, str)
	return ok, nil
}

func (uc *WhitelistUseCase) Update(ctx context.Context, item *UpdateWhitelist) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *WhitelistUseCase) Delete(ctx context.Context, ids ...int64) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Delete(ctx, ids...)
		})
	})
}
