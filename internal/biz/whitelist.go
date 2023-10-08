package biz

import (
	"context"
	"strconv"
	"strings"

	"auth/internal/conf"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
)

const (
	WhitelistPermissionCategory uint32 = iota
	WhitelistJwtCategory
	WhitelistIdempotentCategory
)

type Whitelist struct {
	Id       uint64 `json:"id,string"`
	Category uint32 `json:"category"`
	Resource string `json:"resource"`
}

type FindWhitelist struct {
	Page     page.Page `json:"page"`
	Category *uint32   `json:"category"`
	Resource *string   `json:"resource"`
}

type HasWhitelist struct {
	Category   uint32          `json:"category"`
	Permission CheckPermission `json:"permission"`
}

type FindWhitelistCache struct {
	Page page.Page   `json:"page"`
	List []Whitelist `json:"list"`
}

type HasWhitelistCache struct {
	Has bool `json:"has"`
}

type UpdateWhitelist struct {
	Id       uint64  `json:"id,string"`
	Category *uint32 `json:"category,omitempty"`
	Resource *string `json:"resource,omitempty"`
}

type WhitelistRepo interface {
	Create(ctx context.Context, item *Whitelist) error
	Find(ctx context.Context, condition *FindWhitelist) []Whitelist
	Has(ctx context.Context, condition *HasWhitelist) bool
	Update(ctx context.Context, item *UpdateWhitelist) error
	Delete(ctx context.Context, ids ...uint64) error
}

type WhitelistUseCase struct {
	c          *conf.Bootstrap
	repo       WhitelistRepo
	actionRepo WhitelistRepo
	tx         Transaction
	cache      Cache
}

func NewWhitelistUseCase(c *conf.Bootstrap, repo WhitelistRepo, tx Transaction, cache Cache) *WhitelistUseCase {
	return &WhitelistUseCase{
		c:     c,
		repo:  repo,
		tx:    tx,
		cache: cache.WithPrefix("whitelist"),
	}
}

func (uc *WhitelistUseCase) Create(ctx context.Context, item *Whitelist) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *WhitelistUseCase) Find(ctx context.Context, condition *FindWhitelist) (rp []Whitelist) {
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache FindWhitelistCache
	utils.Json2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
}

func (uc *WhitelistUseCase) find(ctx context.Context, action string, condition *FindWhitelist) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache FindWhitelistCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2Json(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

func (uc *WhitelistUseCase) Has(ctx context.Context, condition *HasWhitelist) (rp bool, err error) {
	action := strings.Join([]string{"has", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.has(ctx, action, condition)
	})
	if err != nil {
		return
	}
	rp, _ = strconv.ParseBool(str)
	return
}

func (uc *WhitelistUseCase) has(ctx context.Context, action string, condition *HasWhitelist) (res string, err error) {
	// read data from db and write to cache
	has := uc.repo.Has(ctx, condition)
	res = strconv.FormatBool(has)
	uc.cache.Set(ctx, action, res, false)
	return
}

func (uc *WhitelistUseCase) Update(ctx context.Context, item *UpdateWhitelist) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *WhitelistUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}
