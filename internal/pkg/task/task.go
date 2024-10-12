package task

import (
	"context"
	"strings"
	"time"

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
		worker.WithHandlerNeedWorker(func(ctx context.Context, w worker.Worker, p worker.Payload) error {
			return process(task{
				ctx:     ctx,
				c:       c,
				w:       w,
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

	for id, item := range c.Task.Cron {
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
	// when app restart, clear hotspot
	_ = w.Once(
		worker.WithRunUUID(strings.Join([]string{c.Task.Group.RefreshHotspotManual}, ".")),
		worker.WithRunGroup(c.Task.Group.RefreshHotspotManual),
		worker.WithRunIn(10*time.Second),
		worker.WithRunReplace(true),
	)
	return
}

type task struct {
	ctx     context.Context
	c       *conf.Bootstrap
	w       worker.Worker
	payload worker.Payload
	user    *biz.UserUseCase
	hotspot *biz.HotspotUseCase
}

func process(t task) (err error) {
	tr := otel.Tracer("task")
	ctx, span := tr.Start(t.ctx, "Task")
	defer span.End()
	switch t.payload.Group {
	case t.c.Task.Group.LoginFailed:
		var req biz.LoginTime
		utils.Json2Struct(&req, t.payload.Payload)
		err = t.user.WrongPwd(ctx, req)
	case t.c.Task.Group.LoginLast:
		var req biz.LoginTime
		utils.Json2Struct(&req, t.payload.Payload)
		err = t.user.LastLogin(ctx, req.Username)
	case t.c.Task.Group.RefreshHotspot:
		err = t.hotspot.Refresh(ctx)
	case t.c.Task.Group.RefreshHotspotManual:
		err = t.hotspot.Refresh(ctx)
	}
	return
}
