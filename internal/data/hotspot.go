package data

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	gocache "github.com/patrickmn/go-cache"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"auth/internal/biz"
	"auth/internal/data/model"
)

// In-memory hotspot cache (per-process).
//
// Notes:
// - Keep user status cache TTL short to avoid stale wrong/lock values affecting login flow.
// - Use Refresh (triggered by task/service) to clear caches after mutating operations.
const (
	hotspotUserByCodeTTL     = 10 * time.Minute
	hotspotUserByUsernameTTL = 3 * time.Second
	hotspotActionTTL         = 30 * time.Minute
	hotspotPermissionTTL     = 10 * time.Minute

	hotspotCleanupInterval = 5 * time.Minute
)

type hotspotRepo struct {
	data *Data

	userByCode     *gocache.Cache // key: user.code
	userByUsername *gocache.Cache // key: user.username

	actionByCode    *gocache.Cache // key: action.code
	userActionCodes *gocache.Cache // key: userID(string) -> []string
	userPermissions *gocache.Cache // key: userID(string) -> []biz.Action
}

func NewHotspotRepo(data *Data) biz.HotspotRepo {
	return &hotspotRepo{
		data: data,

		userByCode:     gocache.New(hotspotUserByCodeTTL, hotspotCleanupInterval),
		userByUsername: gocache.New(hotspotUserByUsernameTTL, hotspotCleanupInterval),

		actionByCode:    gocache.New(hotspotActionTTL, hotspotCleanupInterval),
		userActionCodes: gocache.New(hotspotPermissionTTL, hotspotCleanupInterval),
		userPermissions: gocache.New(hotspotPermissionTTL, hotspotCleanupInterval),
	}
}

func (ro *hotspotRepo) Refresh(ctx context.Context) error {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Refresh")
	defer span.End()

	// Clear all in-memory caches (best-effort warmup happens below).
	ro.userByCode.Flush()
	ro.userByUsername.Flush()
	ro.actionByCode.Flush()
	ro.userActionCodes.Flush()
	ro.userPermissions.Flush()

	// Warm action cache since it is usually small and heavily used by permission checks.
	list, err := gorm.G[model.Action](ro.data.DB(ctx)).Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("refresh hotspot: load actions failed")
		return err
	}
	for _, m := range list {
		code := strings.TrimSpace(m.Code)
		if code == "" {
			continue
		}

		var a biz.Action
		if err := copierx.Copy(&a, m); err != nil {
			log.WithContext(ctx).WithError(err).Error("refresh hotspot: copy action failed")
			continue
		}
		ro.actionByCode.Set(code, a, hotspotActionTTL)
	}
	return nil
}

// getUserByField is a helper to fetch user by a specific field (code or username).
func (ro *hotspotRepo) getUserByField(
	ctx context.Context,
	field, value string,
	primaryCache *gocache.Cache,
	primaryTTL time.Duration,
	secondaryCache *gocache.Cache,
	secondaryTTL time.Duration,
	getSecondaryKey func(*biz.User) string,
) *biz.User {
	item := &biz.User{}
	value = strings.TrimSpace(value)
	if value == "" {
		return item
	}

	// Check primary cache
	if v, ok := primaryCache.Get(value); ok {
		if u, ok := v.(biz.User); ok {
			uc := u
			return &uc
		}
	}

	// Query database
	db := gorm.G[model.User](ro.data.DB(ctx))
	m, err := db.Where(field+" = ?", value).First(ctx)
	if err != nil {
		return item
	}

	copierx.Copy(item, m)
	item.Id = m.ID
	item.Locked = m.Locked != nil && *m.Locked != 0

	// Update both caches
	primaryCache.Set(value, *item, primaryTTL)
	if secondaryKey := getSecondaryKey(item); strings.TrimSpace(secondaryKey) != "" {
		secondaryCache.Set(secondaryKey, *item, secondaryTTL)
	}
	return item
}

func (ro *hotspotRepo) GetUserByCode(ctx context.Context, code string) *biz.User {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetUserByCode")
	defer span.End()

	return ro.getUserByField(
		ctx, "code", code,
		ro.userByCode, hotspotUserByCodeTTL,
		ro.userByUsername, hotspotUserByUsernameTTL,
		func(u *biz.User) string { return u.Username },
	)
}

func (ro *hotspotRepo) GetUserByUsername(ctx context.Context, username string) *biz.User {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetUserByUsername")
	defer span.End()

	return ro.getUserByField(
		ctx, "username", username,
		ro.userByUsername, hotspotUserByUsernameTTL,
		ro.userByCode, hotspotUserByCodeTTL,
		func(u *biz.User) string { return u.Code },
	)
}

func (ro *hotspotRepo) GetActionByCode(ctx context.Context, code string) *biz.Action {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetActionByCode")
	defer span.End()

	item := &biz.Action{}
	code = strings.TrimSpace(code)
	if code == "" {
		return item
	}

	if v, ok := ro.actionByCode.Get(code); ok {
		if a, ok := v.(biz.Action); ok {
			ac := a
			return &ac
		}
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	m, err := db.Where("code = ?", code).First(ctx)
	if err != nil {
		return item
	}

	copierx.Copy(item, m)
	ro.actionByCode.Set(code, *item, hotspotActionTTL)
	return item
}

func (ro *hotspotRepo) FindUserPermissions(ctx context.Context, userID int64) ([]*biz.Action, error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "FindUserPermissions")
	defer span.End()

	if userID == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "userID")
	}
	key := strconv.FormatInt(userID, 10)

	if v, ok := ro.userPermissions.Get(key); ok {
		if list, ok := v.([]biz.Action); ok {
			return hotspotCloneActions(list), nil
		}
	}

	codes, err := ro.getUserActionCodes(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		ro.userPermissions.Set(key, []biz.Action{}, hotspotPermissionTTL)
		return []*biz.Action{}, nil
	}

	byCode := make(map[string]biz.Action, len(codes))
	missing := make([]string, 0)
	for _, c := range codes {
		if v, ok := ro.actionByCode.Get(c); ok {
			if a, ok := v.(biz.Action); ok {
				byCode[c] = a
				continue
			}
		}
		missing = append(missing, c)
	}

	if len(missing) > 0 {
		list, qErr := gorm.G[model.Action](ro.data.DB(ctx)).
			Where("code IN ?", missing).
			Find(ctx)
		if qErr != nil {
			log.WithContext(ctx).WithError(qErr).Error("load hotspot user permissions actions failed")
			return nil, qErr
		}
		for _, m := range list {
			c := strings.TrimSpace(m.Code)
			if c == "" {
				continue
			}
			var a biz.Action
			if err := copierx.Copy(&a, m); err != nil {
				log.WithContext(ctx).WithError(err).Error("copy hotspot action failed")
				continue
			}
			ro.actionByCode.Set(c, a, hotspotActionTTL)
			byCode[c] = a
		}
	}

	// Preserve code order.
	cachedList := make([]biz.Action, 0, len(codes))
	res := make([]*biz.Action, 0, len(codes))
	for _, c := range codes {
		a, ok := byCode[c]
		if !ok {
			continue
		}
		cachedList = append(cachedList, a)
		ac := a
		res = append(res, &ac)
	}

	ro.userPermissions.Set(key, cachedList, hotspotPermissionTTL)
	return res, nil
}

func (ro *hotspotRepo) CheckPermission(ctx context.Context, userID int64, resource, method string) (bool, error) {
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

	actions, err := ro.FindUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	req := permissionMatchReq{}
	if method == "" {
		req.Resource = resource
	} else {
		req.Method = method
		req.URI = resource
	}

	for _, a := range actions {
		if a == nil {
			continue
		}
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

func (ro *hotspotRepo) getUserActionCodes(ctx context.Context, userID int64) ([]string, error) {
	if userID == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "userID")
	}
	key := strconv.FormatInt(userID, 10)

	if v, ok := ro.userActionCodes.Get(key); ok {
		if codes, ok := v.([]string); ok {
			return append([]string(nil), codes...), nil
		}
	}

	// Reuse the canonical permission aggregation logic, then cache the result.
	codes, err := (permissionRepo{data: ro.data}).getUserActionCodes(ctx, userID)
	if err != nil {
		return nil, err
	}
	ro.userActionCodes.Set(key, append([]string(nil), codes...), hotspotPermissionTTL)
	return codes, nil
}

func hotspotCloneActions(list []biz.Action) []*biz.Action {
	rp := make([]*biz.Action, 0, len(list))
	for i := range list {
		a := list[i]
		rp = append(rp, &a)
	}
	return rp
}
