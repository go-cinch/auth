package task

import (
	"context"

	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/go-cinch/common/worker"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

// ProviderSet is task providers.
var ProviderSet = wire.NewSet(New)

// New is initialize task worker from config
func New(c *conf.Bootstrap, user *biz.UserUseCase, hotspot *biz.HotspotUseCase) (w *worker.Worker, err error) {
	w = worker.New(
		worker.WithRedisURI(c.Data.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			return process(task{
				ctx:     ctx,
				payload: p,
				user:    user,
				hotspot: hotspot,
			})
		}),
	)
	if w.Error != nil {
		log.Error(w.Error)
		err = errors.New("initialize worker failed")
		return
	}

	for id, item := range c.Task {
		err = w.Cron(
			worker.WithRunUUID(id),
			worker.WithRunGroup(item.Name),
			worker.WithRunExpr(item.Expr),
			worker.WithRunTimeout(int(item.Timeout)),
			worker.WithRunMaxRetry(int(item.Retry)),
		)
		if err != nil {
			log.Error(err)
			err = errors.New("initialize worker failed")
			return
		}
	}

	log.Info("initialize worker success")
	return
}

type task struct {
	ctx     context.Context
	payload worker.Payload
	user    *biz.UserUseCase
	hotspot *biz.HotspotUseCase
}

func process(t task) (err error) {
	tr := otel.Tracer("task")
	ctx, span := tr.Start(t.ctx, "Task")
	defer span.End()
	switch t.payload.Group {
	case "login.failed":
		var req biz.LoginTime
		utils.Json2Struct(&req, t.payload.Payload)
		err = t.user.WrongPwd(ctx, req)
	case "login.last":
		var req biz.LoginTime
		utils.Json2Struct(&req, t.payload.Payload)
		err = t.user.LastLogin(ctx, req.Username)
	case "refresh.hotspot":
		err = t.hotspot.Refresh(ctx)
	}
	return
}
