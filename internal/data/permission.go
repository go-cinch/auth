package data

import (
	"auth/internal/biz"
	"context"
	"fmt"
)

type permissionRepo struct {
	data      *Data
	action    biz.ActionRepo
	user      biz.UserRepo
	userGroup biz.UserGroupRepo
}

// NewPermissionRepo .
func NewPermissionRepo(data *Data, action biz.ActionRepo, user biz.UserRepo, userGroup biz.UserGroupRepo) biz.PermissionRepo {
	return &permissionRepo{
		data:      data,
		action:    action,
		user:      user,
		userGroup: userGroup,
	}
}

func (ro permissionRepo) Check(ctx context.Context, item *biz.Permission) (pass bool) {
	user, err := ro.user.GetByCode(ctx, item.UserCode)
	if err != nil {
		return
	}
	resource := fmt.Sprintf("%s,%s", item.Method, item.Uri)
	// 1. check user permission
	pass = ro.action.Permission(ctx, user.Action, resource)
	if pass {
		return
	}
	// 2. check role permission
	pass = ro.action.Permission(ctx, user.Role.Action, resource)
	if pass {
		return
	}
	// 3. check user group permission
	groups, err := ro.userGroup.FindGroupByUserCode(ctx, user.Code)
	if err != nil {
		return
	}
	for _, group := range groups {
		pass = ro.action.Permission(ctx, group.Action, resource)
		if pass {
			return
		}
	}
	return
}
