package service

import (
	"context"
	"strings"

	"github.com/go-cinch/common/worker"
)

func (s *AuthService) flushCache(ctx context.Context) {
	// clear user info cache
	s.user.FlushCache(ctx)

	// delay clear hotspot cache
	_ = s.task.Once(
		worker.WithRunCtx(ctx),
		worker.WithRunUUID(strings.Join([]string{s.c.Task.Group.RefreshHotspotManual}, ".")),
		worker.WithRunGroup(s.c.Task.Group.RefreshHotspotManual),
		worker.WithRunNow(true),
		worker.WithRunReplace(true),
	)
}
