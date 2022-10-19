package middleware

import (
	"auth/api/auth"
	"context"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

func permissionWhitelist() selector.MatchFunc {
	whitelist := make(map[string]struct{})
	whitelist[auth.OperationAuthLogin] = struct{}{}
	whitelist[auth.OperationAuthStatus] = struct{}{}
	whitelist[auth.OperationAuthLogout] = struct{}{}
	whitelist[auth.OperationAuthCaptcha] = struct{}{}
	whitelist[auth.OperationAuthRegister] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whitelist[operation]; ok {
			return false
		}
		return true
	}
}

func idempotentBlacklist() selector.MatchFunc {
	blacklist := make(map[string]struct{})
	blacklist[auth.OperationAuthCreateAction] = struct{}{}
	blacklist[auth.OperationAuthCreateRole] = struct{}{}
	blacklist[auth.OperationAuthCreateUserGroup] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := blacklist[operation]; ok {
			return true
		}
		return false
	}
}
