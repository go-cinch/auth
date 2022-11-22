package task

import (
	"auth/internal/biz"
	"auth/internal/conf"
	"context"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/go-cinch/common/worker"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

// ProviderSet is task providers.
var ProviderSet = wire.NewSet(NewTask)

type Task struct {
	worker *worker.Worker
}

func (tk Task) Once(options ...func(*worker.RunOptions)) error {
	return tk.worker.Once(options...)
}

func (tk Task) Cron(options ...func(*worker.RunOptions)) error {
	return tk.worker.Cron(options...)
}

// NewTask is initialize task worker from config
func NewTask(c *conf.Bootstrap, user *biz.UserUseCase) (tk *Task, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.Errorf("%v", e)
		}
	}()
	w := worker.New(
		worker.WithRedisUri(c.Data.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			return process(task{
				ctx:     ctx,
				payload: p,
				user:    user,
			})
		}),
	)
	if w.Error != nil {
		err = errors.WithMessage(w.Error, "initialize worker failed")
		return
	}

	tk = &Task{
		worker: w,
	}

	if len(c.Tasks) > 0 {
		// TODO register cron tasks
	}

	log.Info("initialize worker success")
	return
}

type task struct {
	ctx     context.Context
	payload worker.Payload
	user    *biz.UserUseCase
}

func process(t task) (err error) {
	switch t.payload.Group {
	case "login.failed":
		var req biz.LoginTime
		utils.Json2Struct(&req, t.payload.Payload)
		err = t.user.WrongPwd(t.ctx, req)
	case "login.last":
		var req biz.LoginTime
		utils.Json2Struct(&req, t.payload.Payload)
		err = t.user.LastLogin(t.ctx, req.Username)
	}
	return
}
