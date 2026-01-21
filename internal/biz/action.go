package biz

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/page/v2"
	"github.com/go-cinch/common/utils"

	"auth/internal/conf"
)

// Action defines a fine-grained permission rule for a resource/menu/button.
type Action struct {
	ID       int64     `json:"id,string"`
	Code     *string   `json:"code,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Word     *string   `json:"word,omitempty"`
	Resource *string   `json:"resource,omitempty"`
	Menu     *string   `json:"menu,omitempty"`
	Btn      *string   `json:"btn,omitempty"`
	Children []*Action `json:"children,omitempty"`
}

type FindAction struct {
	Page     page.Page `json:"page"`
	Code     *string   `json:"code,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Word     *string   `json:"word,omitempty"`
	Resource *string   `json:"resource,omitempty"`
}

type FindActionCache struct {
	Page page.Page `json:"page"`
	List []*Action `json:"list"`
}

type ActionUseCase struct {
	c     *conf.Bootstrap
	repo  ActionRepo
	tx    Transaction
	cache Cache
}

func NewActionUseCase(c *conf.Bootstrap, repo ActionRepo, tx Transaction, cache Cache) *ActionUseCase {
	return &ActionUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "action",
		}, "_")),
	}
}

func (uc *ActionUseCase) Create(ctx context.Context, item *Action) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *ActionUseCase) Find(ctx context.Context, condition *FindAction) (rp []Action, err error) {
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
	var cache FindActionCache
	utils.JSON2Struct(&cache, str)
	condition.Page = cache.Page
	rp = make([]Action, 0, len(cache.List))
	for _, item := range cache.List {
		if item != nil {
			rp = append(rp, *item)
		}
	}
	return rp, nil
}

func (uc *ActionUseCase) find(ctx context.Context, action string, condition *FindAction) (res string, err error) {
	list := uc.repo.Find(ctx, condition)
	var cache FindActionCache
	cache.Page = condition.Page
	cache.List = make([]*Action, 0, len(list))
	for i := range list {
		cache.List = append(cache.List, &list[i])
	}
	res = utils.Struct2JSON(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return res, nil
}

func (uc *ActionUseCase) Update(ctx context.Context, item *Action) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *ActionUseCase) Delete(ctx context.Context, ids ...int64) error {
	tr := otel.Tracer("biz")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			return uc.repo.Delete(ctx, ids)
		})
	})
}
