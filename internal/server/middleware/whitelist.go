package middleware

import (
	v1 "auth/api/auth/v1"
	"context"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

func permissionWhitelist() selector.MatchFunc {
	whitelist := make(map[string]struct{})
	whitelist[v1.OperationAuthLogin] = struct{}{}
	whitelist[v1.OperationAuthStatus] = struct{}{}
	whitelist[v1.OperationAuthLogout] = struct{}{}
	whitelist[v1.OperationAuthCaptcha] = struct{}{}
	whitelist[v1.OperationAuthRegister] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whitelist[operation]; ok {
			return false
		}
		return true
	}
}

func idempotentBlacklist() selector.MatchFunc {
	blacklist := make(map[string]struct{})
	blacklist[v1.OperationAuthCreateAction] = struct{}{}
	blacklist[v1.OperationAuthCreateRole] = struct{}{}
	blacklist[v1.OperationAuthCreateUserGroup] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := blacklist[operation]; ok {
			return true
		}
		return false
	}
}
