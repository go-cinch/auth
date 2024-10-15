package service

import (
	"auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
	"github.com/go-cinch/common/idempotent"
	"github.com/go-cinch/common/worker"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthService)

// AuthService is a auth service.
type AuthService struct {
	auth.UnimplementedAuthServer

	c          *conf.Bootstrap
	task       *worker.Worker
	idempotent *idempotent.Idempotent
	user       *biz.UserUseCase
	action     *biz.ActionUseCase
	role       *biz.RoleUseCase
	userGroup  *biz.UserGroupUseCase
	permission *biz.PermissionUseCase
	whitelist  *biz.WhitelistUseCase
}

// NewAuthService new an auth service.
func NewAuthService(
	c *conf.Bootstrap,
	task *worker.Worker,
	idempotent *idempotent.Idempotent,
	user *biz.UserUseCase,
	action *biz.ActionUseCase,
	role *biz.RoleUseCase,
	userGroup *biz.UserGroupUseCase,
	permission *biz.PermissionUseCase,
	whitelist *biz.WhitelistUseCase,
) *AuthService {
	return &AuthService{
		c:          c,
		task:       task,
		idempotent: idempotent,
		user:       user,
		action:     action,
		role:       role,
		userGroup:  userGroup,
		permission: permission,
		whitelist:  whitelist,
	}
}
