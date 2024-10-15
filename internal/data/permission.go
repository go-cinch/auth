package data

import (
	"context"
	"strings"

	"auth/internal/biz"
	"github.com/samber/lo"
)

type permissionRepo struct {
	data    *Data
	action  biz.ActionRepo
	hotspot biz.HotspotRepo
}

// NewPermissionRepo .
func NewPermissionRepo(data *Data, action biz.ActionRepo, hotspot biz.HotspotRepo) biz.PermissionRepo {
	return &permissionRepo{
		data:    data,
		action:  action,
		hotspot: hotspot,
	}
}

func (ro permissionRepo) Check(ctx context.Context, item *biz.CheckPermission) (pass bool) {
	user := ro.hotspot.GetUserByCode(ctx, item.UserCode)
	// 1. check default permission
	defaultAction := ro.hotspot.GetActionByWord(ctx, "default")
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
	groups := ro.hotspot.FindUserGroupByUserCode(ctx, user.Code)
	for _, group := range groups {
		pass = ro.action.Permission(ctx, group.Action, item)
		if pass {
			return
		}
	}
	return
}

func (ro permissionRepo) GetByUserCode(ctx context.Context, code string) (rp *biz.Permission) {
	rp = &biz.Permission{}
	rp.Resources = make([]string, 0)
	user := ro.hotspot.GetUserByCode(ctx, code)
	// 1. user action
	actions := make([]string, 0)
	defaultAction := ro.hotspot.GetActionByWord(ctx, "default")
	actions = append(actions, defaultAction.Code)
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
	groups := ro.hotspot.FindUserGroupByUserCode(ctx, user.Code)
	for _, group := range groups {
		if group.Action != "" {
			arr := strings.Split(group.Action, ",")
			actions = append(actions, arr...)
		}
	}
	actions = lo.Uniq(actions)
	if len(actions) > 0 {
		list := ro.hotspot.FindActionByCode(ctx, actions...)
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
	rp.Resources = lo.Uniq(rp.Resources)
	rp.Menus = lo.Uniq(rp.Menus)
	rp.Btns = lo.Uniq(rp.Btns)
	return
}
