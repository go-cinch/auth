package idempotent

import (
	"auth/internal/conf"
	"github.com/go-cinch/common/idempotent"
	"github.com/go-cinch/common/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(New)

// New is initialize idempotent from config
func New(c *conf.Bootstrap, redis redis.UniversalClient) (idt *idempotent.Idempotent, err error) {
	idt = idempotent.New(
		idempotent.WithPrefix(c.Name),
		idempotent.WithRedis(redis),
	)
	log.Info("initialize idempotent success")
	return
}
