package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/task"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthService)

// AuthService is a greeter service.
type AuthService struct {
	v1.UnimplementedAuthServer

	task *task.Task
	user *biz.UserUseCase
}

// NewAuthService new an auth service.
func NewAuthService(task *task.Task, user *biz.UserUseCase) *AuthService {
	return &AuthService{task: task, user: user}
}
