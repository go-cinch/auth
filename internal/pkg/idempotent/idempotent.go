package idempotent

import (
	"auth/internal/conf"
	"context"
	"github.com/go-cinch/common/idempotent"
	"github.com/go-cinch/common/log"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(NewIdempotent)

type Idempotent struct {
	idempotent *idempotent.Idempotent
}

func (idt Idempotent) Token(ctx context.Context) string {
	return idt.idempotent.Token(ctx)
}

func (idt Idempotent) Check(ctx context.Context, token string) bool {
	return idt.idempotent.Check(ctx, token)
}

// NewIdempotent is initialize idempotent from config
func NewIdempotent(c *conf.Bootstrap, redis redis.UniversalClient) (idt *Idempotent, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.Errorf("%v", e)
		}
	}()
	ins := idempotent.New(
		idempotent.WithPrefix(c.Name),
		idempotent.WithRedis(redis),
	)
	idt = &Idempotent{
		idempotent: ins,
	}
	log.Info("initialize idempotent success")
	return
}
