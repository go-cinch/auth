//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/idempotent"
	"auth/internal/server"
	"auth/internal/service"
	"auth/internal/task"
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(c *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, task.ProviderSet, idempotent.ProviderSet, service.ProviderSet, newApp))
}
