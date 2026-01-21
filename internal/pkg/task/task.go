package task

import (
	"context"
	"strings"
	"time"

	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/go-cinch/common/worker"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

	"auth/internal/biz"
	"auth/internal/conf"
)

// ProviderSet is task providers.
var ProviderSet = wire.NewSet(New)

// New initializes the task worker from config.
func New(c *conf.Bootstrap, user *biz.UserUseCase, hotspot biz.HotspotRepo) (w *worker.Worker, err error) {
	w = worker.New(
		worker.WithRedisURI(c.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			return process(task{
				ctx:     ctx,
				c:       c,
				payload: p,
				user:    user,
				hotspot: hotspot,
			})
		}),
	)
	if w.Error != nil {
		log.Error(w.Error)
		return nil, errors.New("initialize worker failed")
	}

	for id, item := range c.Task.Cron {
		err = w.Cron(
			context.Background(),
			worker.WithRunUUID(id),
			worker.WithRunGroup(item.Name),
			worker.WithRunExpr(item.Expr),
			worker.WithRunTimeout(int(item.Timeout)),
			worker.WithRunMaxRetry(int(item.Retry)),
		)
		if err != nil {
			log.Error(err)
			return nil, errors.New("initialize worker failed")
		}
	}

	log.Info("initialize worker success")
	// When app restart, clear hotspot (best-effort).
	if c.Task.Group.RefreshHotspotManual != "" {
		_ = w.Once(
			context.Background(),
			worker.WithRunUUID(strings.Join([]string{c.Task.Group.RefreshHotspotManual}, ".")),
			worker.WithRunGroup(c.Task.Group.RefreshHotspotManual),
			worker.WithRunIn(10*time.Second),
			worker.WithRunReplace(true),
		)
	}

	return w, nil
}

type task struct {
	ctx     context.Context
	c       *conf.Bootstrap
	payload worker.Payload
	user    *biz.UserUseCase
	hotspot biz.HotspotRepo
}

func process(t task) (err error) {
	tr := otel.Tracer("task")
	ctx, span := tr.Start(t.ctx, "process")
	defer span.End()

	// Use task group to match tasks instead of UID.
	switch t.payload.Group {
	case t.c.Task.Group.LoginFailed:
		var req biz.LoginTime
		utils.JSON2Struct(&req, t.payload.Payload)
		err = t.user.WrongPwd(ctx, &req)
	case t.c.Task.Group.LoginLast:
		var req biz.LoginTime
		utils.JSON2Struct(&req, t.payload.Payload)
		err = t.user.LastLogin(ctx, req.Username)
	case t.c.Task.Group.RefreshHotspot, t.c.Task.Group.RefreshHotspotManual:
		err = t.hotspot.Refresh(ctx)
	default:
		log.WithContext(ctx).Warn("unknown task group: %s", t.payload.Group)
	}
	return err
}
