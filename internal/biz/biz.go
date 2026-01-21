package biz

import (
	"context"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewAuthUseCase,
	NewUserUseCase,
	NewPermissionUseCase,
	NewRoleUseCase,
	NewActionUseCase,
	NewUserGroupUseCase,
	NewWhitelistUseCase,
)

type Transaction interface {
	Tx(ctx context.Context, handler func(context.Context) error) error
}

type Cache interface {
	// Cache is get redis instance
	Cache() redis.UniversalClient
	// WithPrefix will add cache key prefix
	WithPrefix(prefix string) Cache
	// WithRefresh get data from db skip cache and refresh cache
	WithRefresh() Cache
	// Get is get cache data by key from redis, do write handler if cache is empty
	Get(ctx context.Context, action string, write func(context.Context) (string, error)) (string, error)
	// Set is set data to redis
	Set(ctx context.Context, action, data string, short bool)
	// Del delete key
	Del(ctx context.Context, action string)
	// SetWithExpiration is set data to redis with custom expiration
	SetWithExpiration(ctx context.Context, action, data string, seconds int64)
	// Flush is clean association cache if handler err=nil
	Flush(ctx context.Context, handler func(context.Context) error) error
	// FlushByPrefix clean cache by prefix, without prefix equals flush all by default cache prefix
	FlushByPrefix(ctx context.Context, prefix ...string) (err error)
}

type HealthRepo interface {
	PingDB(ctx context.Context) error
	PingRedis(ctx context.Context) error
}

// HotspotRepo provides cached lookup helpers used by the auth domain.
//
// It is intended for "hot" read paths (login/status/permission checks) where we
// want to avoid repeated DB reads. Implementations should keep data fresh via
// periodic/manual Refresh calls.
type HotspotRepo interface {
	// Refresh clears and/or rebuilds the in-memory hotspot caches.
	Refresh(ctx context.Context) error

	// User lookups.
	GetUserByCode(ctx context.Context, code string) *User
	GetUserByUsername(ctx context.Context, username string) *User

	// Permission/Action helpers (optional callers).
	GetActionByCode(ctx context.Context, code string) *Action
	FindUserPermissions(ctx context.Context, userID int64) ([]*Action, error)
	CheckPermission(ctx context.Context, userID int64, resource, method string) (bool, error)
}

type AuthRepo interface {
	// GetLoginUser returns a user by username with password hash and JWT attrs.
	GetLoginUser(ctx context.Context, username string) (*LoginUser, error)
}

type PermissionRepo interface {
	// FindUserPermissions returns all actions a user can access.
	FindUserPermissions(ctx context.Context, userID int64) ([]*Action, error)
	// CheckPermission checks if a user is allowed to access the resource.
	// If method is empty, resource is treated as a gRPC operation (exact match).
	// If method is not empty, resource is treated as an HTTP URI and matched against action rules.
	CheckPermission(ctx context.Context, userID int64, resource, method string) (bool, error)
}

type RoleRepo interface {
	Create(ctx context.Context, item *Role) error
	Find(ctx context.Context, condition *FindRole) []Role
	Update(ctx context.Context, item *UpdateRole) error
	Delete(ctx context.Context, ids ...int64) error
}

type WhitelistRepo interface {
	Create(ctx context.Context, item *Whitelist) error
	Update(ctx context.Context, item *UpdateWhitelist) error
	Delete(ctx context.Context, ids ...int64) error
	Find(ctx context.Context, condition *FindWhitelist) []Whitelist

	// Match checks whether the given resource matches any whitelist rules in the given category.
	//
	// The resource format supports:
	// - gRPC: "package.Service/Method" (exact match)
	// - HTTP: "METHOD|/uri/path" or "METHOD|/uri/path|grpcResource"
	Match(ctx context.Context, category int16, resource string) (bool, error)
}

type ActionRepo interface {
	Create(ctx context.Context, item *Action) error
	Update(ctx context.Context, item *Action) error
	Delete(ctx context.Context, ids []int64) error

	Find(ctx context.Context, condition *FindAction) []Action
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	Find(ctx context.Context, condition *FindUser) []User
	Create(ctx context.Context, item *User) error
	Update(ctx context.Context, item *UpdateUser) error
	Delete(ctx context.Context, ids ...int64) error
	LastLogin(ctx context.Context, username string) error
	WrongPwd(ctx context.Context, req *LoginTime) error
	UpdatePassword(ctx context.Context, item *User) error
	IdExists(ctx context.Context, id int64) error
}

type UserGroupRepo interface {
	Create(ctx context.Context, item *UserGroup) error
	Find(ctx context.Context, condition *FindUserGroup) []UserGroup
	Update(ctx context.Context, item *UpdateUserGroup) error
	Delete(ctx context.Context, ids ...int64) error
}
