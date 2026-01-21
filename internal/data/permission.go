package data

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type permissionRepo struct {
	data *Data
}

func NewPermissionRepo(data *Data) biz.PermissionRepo {
	return &permissionRepo{
		data: data,
	}
}

// FindUserPermissions returns all actions a user can access.
// Permissions are aggregated from:
// - default action (word = "default") if present
// - user.action
// - role.action (via user.role_id)
// - user_group.action (via user_user_group_relation) (enabled).
func (ro permissionRepo) FindUserPermissions(ctx context.Context, userID int64) (rp []*biz.Action, err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "FindUserPermissions")
	defer span.End()

	rp = make([]*biz.Action, 0)

	codes, err := ro.getUserActionCodes(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return rp, nil
	}

	list, err := gorm.G[model.Action](ro.data.DB(ctx)).
		Where("code IN ?", codes).
		Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find user permission actions failed")
		return nil, err
	}

	// Preserve code order.
	byCode := make(map[string]model.Action, len(list))
	for _, m := range list {
		c := strings.TrimSpace(m.Code)
		if c == "" {
			continue
		}
		byCode[c] = m
	}

	rp = make([]*biz.Action, 0, len(codes))
	for _, code := range codes {
		m, ok := byCode[code]
		if !ok {
			continue
		}
		it := new(biz.Action)
		if err := copierx.Copy(it, m); err != nil {
			log.WithContext(ctx).WithError(err).Error("copy action failed")
			return nil, err
		}
		rp = append(rp, it)
	}
	return rp, nil
}

// CheckPermission checks whether a user is allowed to access a resource.
// If method is empty, resource is treated as a gRPC operation (exact match).
// If method is not empty, resource is treated as an HTTP URI, and action rules
// are matched using the "METHODS|URI_PATTERN" format (see action.resource).
func (ro permissionRepo) CheckPermission(ctx context.Context, userID int64, resource, method string) (ok bool, err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "CheckPermission")
	defer span.End()

	resource = strings.TrimSpace(resource)
	method = strings.TrimSpace(method)

	if userID == 0 {
		return false, biz.ErrIllegalParameter(ctx, "userID")
	}
	if resource == "" {
		return false, nil
	}

	codes, err := ro.getUserActionCodes(ctx, userID)
	if err != nil {
		return false, err
	}
	if len(codes) == 0 {
		return false, nil
	}

	req := permissionMatchReq{}
	if method == "" {
		req.Resource = resource
	} else {
		req.Method = method
		req.URI = resource
	}

	list, err := gorm.G[model.Action](ro.data.DB(ctx)).
		Select("code, resource").
		Where("code IN ?", codes).
		Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find actions for permission check failed")
		return false, err
	}

	for _, a := range list {
		patterns := strings.TrimSpace(permissionDerefString(a.Resource))
		if patterns == "" {
			continue
		}
		if permissionMatchResource(patterns, req) {
			return true, nil
		}
	}
	return false, nil
}

func (ro permissionRepo) getUserActionCodes(ctx context.Context, userID int64) (codes []string, err error) {
	if userID == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "userID")
	}

	codes = make([]string, 0)
	seen := make(map[string]struct{})

	addUnique := func(arr []string) {
		for _, v := range arr {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			codes = append(codes, v)
		}
	}

	// 0) Default action (optional)
	def, defErr := gorm.G[model.Action](ro.data.DB(ctx)).
		Where("word = ?", "default").
		First(ctx)
	if defErr == nil && def.Code != "" {
		addUnique([]string{def.Code})
	} else if defErr != nil && !errors.Is(defErr, gorm.ErrRecordNotFound) {
		log.WithContext(ctx).WithError(defErr).Error("get default action failed")
		return nil, defErr
	}

	// 1) User
	u, qErr := gorm.G[model.User](ro.data.DB(ctx)).
		Where("id = ?", userID).
		First(ctx)
	if errors.Is(qErr, gorm.ErrRecordNotFound) {
		return nil, biz.ErrRecordNotFound(ctx)
	}
	if qErr != nil {
		log.WithContext(ctx).WithError(qErr).Error("get user failed")
		return nil, qErr
	}
	if u.Action != nil {
		addUnique(permissionSplitComma(*u.Action))
	}

	// 2) Role (optional)
	if u.RoleID != nil && *u.RoleID != 0 {
		r, rErr := gorm.G[model.Role](ro.data.DB(ctx)).
			Where("id = ?", *u.RoleID).
			First(ctx)
		if rErr != nil && !errors.Is(rErr, gorm.ErrRecordNotFound) {
			log.WithContext(ctx).WithError(rErr).Error("get role failed")
			return nil, rErr
		}
		if rErr == nil && r.Action != nil {
			addUnique(permissionSplitComma(*r.Action))
		}
	}
	// 3) User groups (optional)
	rels, rErr := gorm.G[model.UserUserGroupRelation](ro.data.DB(ctx)).
		Where("user_id = ?", userID).
		Find(ctx)
	if rErr != nil {
		log.WithContext(ctx).WithError(rErr).Error("find user_user_group_relation failed")
		return nil, rErr
	}
	if len(rels) == 0 {
		return codes, nil
	}

	groupIDs := make([]int64, 0, len(rels))
	for _, rel := range rels {
		if rel.UserGroupID == 0 {
			continue
		}
		groupIDs = append(groupIDs, rel.UserGroupID)
	}
	groupIDs = permissionUniqueInt64(groupIDs)
	if len(groupIDs) == 0 {
		return codes, nil
	}

	groups, gErr := gorm.G[model.UserGroup](ro.data.DB(ctx)).
		Where("id IN ?", groupIDs).
		Find(ctx)
	if gErr != nil {
		log.WithContext(ctx).WithError(gErr).Error("find user_group failed")
		return nil, gErr
	}
	for _, g := range groups {
		if g.Action == nil {
			continue
		}
		addUnique(permissionSplitComma(*g.Action))
	}

	return codes, nil
}

type permissionMatchReq struct {
	Resource string
	Method   string
	URI      string
}

func permissionMatchResource(patterns string, req permissionMatchReq) bool {
	for _, line := range strings.Split(patterns, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if permissionMatchLine(line, req) {
			return true
		}
	}
	return false
}

func permissionMatchLine(line string, req permissionMatchReq) bool {
	if line == "*" {
		return true
	}
	parts := strings.Split(line, "|")
	switch len(parts) {
	case 1:
		want := strings.TrimSpace(parts[0])
		return want != "" && req.Resource != "" && req.Resource == want
	case 2:
		return permissionMatchHTTP(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), req)
	default:
		methods := strings.TrimSpace(parts[0])
		uriPattern := strings.TrimSpace(parts[1])
		grpcRes := strings.TrimSpace(strings.Join(parts[2:], "|"))
		if grpcRes != "" && req.Resource != "" && req.Resource == grpcRes {
			return true
		}
		return permissionMatchHTTP(methods, uriPattern, req)
	}
}

func permissionMatchHTTP(methodSpec, uriPattern string, req permissionMatchReq) bool {
	if req.Method == "" || req.URI == "" || methodSpec == "" || uriPattern == "" {
		return false
	}
	okMethod := false
	for _, m := range strings.Split(methodSpec, ",") {
		if strings.TrimSpace(m) == req.Method {
			okMethod = true
			break
		}
	}
	if !okMethod {
		return false
	}
	return permissionMatchGlob(uriPattern, req.URI)
}

func permissionMatchGlob(pattern, s string) bool {
	if pattern == "*" {
		return true
	}
	// Fast path: exact match when no wildcards are present.
	if !strings.ContainsAny(pattern, "*?") {
		return pattern == s
	}
	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*`, `.*`)
	re = strings.ReplaceAll(re, `\?`, `.`)
	re = "^" + re + "$"
	return regexp.MustCompile(re).MatchString(s)
}

func permissionDerefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func permissionSplitComma(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, it := range raw {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		if _, ok := seen[it]; ok {
			continue
		}
		seen[it] = struct{}{}
		out = append(out, it)
	}
	return out
}

func permissionUniqueInt64(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	out := make([]int64, 0, len(ids))
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
