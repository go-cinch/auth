//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/server"
	"auth/internal/service"
)

// wireApp initializes the Kratos application.
func wireApp(*conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
