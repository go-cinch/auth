package data

import (
	"context"
	"strings"

	"auth/internal/biz"
	"github.com/go-cinch/common/utils"
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

func (ro permissionRepo) Check(ctx context.Context, item biz.CheckPermission) (pass bool) {
	user, err := ro.user.GetByCode(ctx, item.UserCode)
	if err != nil {
		return
	}
	// 1. check default permission
	defaultAction := ro.action.GetDefault(ctx)
	pass = ro.action.Permission(ctx, defaultAction.Code, item)
	if pass {
		return
	}
	// 2. check user permission
	pass = ro.action.Permission(ctx, user.Action, item)
	if pass {
		return
	}
	// 3. check role permission
	pass = ro.action.Permission(ctx, user.Role.Action, item)
	if pass {
		return
	}
	// 4. check user group permission
	groups := ro.userGroup.FindGroupByUserCode(ctx, user.Code)
	if err != nil {
		return
	}
	for _, group := range groups {
		pass = ro.action.Permission(ctx, group.Action, item)
		if pass {
			return
		}
	}
	return
}

func (ro permissionRepo) GetByUserCode(ctx context.Context, code string) (rp *biz.Permission, err error) {
	rp = &biz.Permission{}
	rp.Resources = make([]string, 0)
	user, err := ro.user.GetByCode(ctx, code)
	if err != nil {
		return
	}
	// 1. user action
	actions := make([]string, 0)
	if user.Action != "" {
		arr := strings.Split(user.Action, ",")
		actions = append(actions, arr...)
	}
	// 2. role action
	if user.Role.Action != "" {
		arr := strings.Split(user.Role.Action, ",")
		actions = append(actions, arr...)
	}
	// 3. user group action
	groups := ro.userGroup.FindGroupByUserCode(ctx, user.Code)
	for _, group := range groups {
		if group.Action != "" {
			arr := strings.Split(group.Action, ",")
			actions = append(actions, arr...)
		}
	}
	actions = utils.RemoveRepeat(actions)
	if len(actions) > 0 {
		list := ro.action.FindByCode(ctx, strings.Join(actions, ","))
		for _, item := range list {
			if item.Resource != "" {
				rp.Resources = append(rp.Resources, strings.Split(item.Resource, "\n")...)
			}
			if item.Menu != "" {
				rp.Menus = append(rp.Menus, strings.Split(item.Menu, "\n")...)
			}
			if item.Btn != "" {
				rp.Btns = append(rp.Btns, strings.Split(item.Btn, "\n")...)
			}
		}
	}
	rp.Resources = utils.RemoveRepeat(rp.Resources)
	rp.Menus = utils.RemoveRepeat(rp.Menus)
	rp.Btns = utils.RemoveRepeat(rp.Btns)
	return
}
