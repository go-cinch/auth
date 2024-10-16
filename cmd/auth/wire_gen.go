// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/pkg/idempotent"
	"auth/internal/pkg/task"
	"auth/internal/server"
	"auth/internal/service"
	"github.com/go-kratos/kratos/v2"
)

import (
	_ "github.com/go-cinch/common/plugins/gorm/filter"
	_ "github.com/go-cinch/common/plugins/kratos/encoding/yml"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(c *conf.Bootstrap) (*kratos.App, func(), error) {
	universalClient, err := data.NewRedis(c)
	if err != nil {
		return nil, nil, err
	}
	idempotentIdempotent, err := idempotent.New(c, universalClient)
	if err != nil {
		return nil, nil, err
	}
	tenant, err := data.NewDB(c)
	if err != nil {
		return nil, nil, err
	}
	sonyflake, err := data.NewSonyflake(c)
	if err != nil {
		return nil, nil, err
	}
	tracerProvider, err := data.NewTracer(c)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup := data.NewData(universalClient, tenant, sonyflake, tracerProvider)
	hotspotRepo := data.NewHotspotRepo(c, dataData)
	actionRepo := data.NewActionRepo(c, dataData, hotspotRepo)
	userRepo := data.NewUserRepo(dataData, actionRepo)
	transaction := data.NewTransaction(dataData)
	cache := data.NewCache(c, universalClient)
	userUseCase := biz.NewUserUseCase(c, userRepo, hotspotRepo, transaction, cache)
	hotspotUseCase := biz.NewHotspotUseCase(c, hotspotRepo)
	worker, err := task.New(c, userUseCase, hotspotUseCase)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	actionUseCase := biz.NewActionUseCase(c, actionRepo, transaction, cache)
	roleRepo := data.NewRoleRepo(dataData, actionRepo)
	roleUseCase := biz.NewRoleUseCase(c, roleRepo, transaction, cache)
	userGroupRepo := data.NewUserGroupRepo(dataData, actionRepo, userRepo)
	userGroupUseCase := biz.NewUserGroupUseCase(c, userGroupRepo, transaction, cache)
	permissionRepo := data.NewPermissionRepo(dataData, actionRepo, hotspotRepo)
	permissionUseCase := biz.NewPermissionUseCase(c, permissionRepo)
	whitelistRepo := data.NewWhitelistRepo(dataData, actionRepo, hotspotRepo)
	whitelistUseCase := biz.NewWhitelistUseCase(c, whitelistRepo, transaction, cache)
	authService := service.NewAuthService(c, worker, idempotentIdempotent, userUseCase, actionUseCase, roleUseCase, userGroupUseCase, permissionUseCase, whitelistUseCase)
	grpcServer := server.NewGRPCServer(c, universalClient, idempotentIdempotent, authService, whitelistUseCase)
	httpServer := server.NewHTTPServer(c, universalClient, idempotentIdempotent, authService, whitelistUseCase)
	app := newApp(grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
