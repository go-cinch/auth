package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/idempotent"
	"auth/internal/task"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthService)

// AuthService is a greeter service.
type AuthService struct {
	v1.UnimplementedAuthServer

	task       *task.Task
	idempotent *idempotent.Idempotent
	user       *biz.UserUseCase
	action     *biz.ActionUseCase
	role       *biz.RoleUseCase
	userGroup  *biz.UserGroupUseCase
	permission *biz.PermissionUseCase
}

// NewAuthService new an auth service.
func NewAuthService(task *task.Task, idempotent *idempotent.Idempotent, user *biz.UserUseCase, action *biz.ActionUseCase, role *biz.RoleUseCase, userGroup *biz.UserGroupUseCase, permission *biz.PermissionUseCase) *AuthService {
	return &AuthService{task: task, idempotent: idempotent, user: user, action: action, role: role, userGroup: userGroup, permission: permission}
}
