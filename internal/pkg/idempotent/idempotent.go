package idempotent

import (
	"context"

	"auth/internal/conf"
	"github.com/go-cinch/common/idempotent"
	"github.com/go-cinch/common/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(New)

type Idempotent struct {
	idempotent *idempotent.Idempotent
}

func (idt Idempotent) Token(ctx context.Context) string {
	return idt.idempotent.Token(ctx)
}

func (idt Idempotent) Check(ctx context.Context, token string) bool {
	return idt.idempotent.Check(ctx, token)
}

// New is initialize idempotent from config
func New(c *conf.Bootstrap, redis redis.UniversalClient) (idt *Idempotent, err error) {
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
