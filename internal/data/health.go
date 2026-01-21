package data

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"

	"auth/internal/biz"
)

type healthRepo struct {
	data  *Data
	redis redis.UniversalClient
}

func NewHealthRepo(data *Data, redis redis.UniversalClient) biz.HealthRepo {
	return &healthRepo{
		data:  data,
		redis: redis,
	}
}

func (r *healthRepo) PingDB(ctx context.Context) error {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "PingDB")
	defer span.End()

	sqlDB, err := r.data.DB(ctx).DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (r *healthRepo) PingRedis(ctx context.Context) error {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "PingRedis")
	defer span.End()

	if r.redis == nil {
		return nil
	}
	return r.redis.Ping(ctx).Err()
}
