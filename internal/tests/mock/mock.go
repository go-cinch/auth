package mock

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/go-cinch/common/mock"
	gormTenant "github.com/go-cinch/common/plugins/gorm/tenant/v2"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/redis/go-redis/v9"

	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/service"
)

func AuthService() *service.AuthService {
	c, _, cache := Data()
	tx := Transaction()

	authUC := biz.NewAuthUseCase(c, AuthRepo(), tx, cache)
	userUC := biz.NewUserUseCase(c, UserRepo(), HotspotRepo(), tx, cache)
	roleUC := biz.NewRoleUseCase(c, RoleRepo(), tx, cache)
	permissionUC := biz.NewPermissionUseCase(c, PermissionRepo())
	actionUC := biz.NewActionUseCase(c, ActionRepo(), tx, cache)
	userGroupUC := biz.NewUserGroupUseCase(c, UserGroupRepo(), tx, cache)
	whitelistUC := biz.NewWhitelistUseCase(c, WhitelistRepo(), tx, cache)

	return service.NewAuthService(
		c,
		authUC,
		userUC,
		roleUC,
		permissionUC,
		actionUC,
		userGroupUC,
		whitelistUC,
		HotspotRepo(),
		HealthRepo(),
	)
}

func Transaction() biz.Transaction {
	_, d, _ := Data()
	return data.NewTransaction(d)
}

func AuthRepo() biz.AuthRepo {
	_, d, _ := Data()
	return data.NewAuthRepo(d)
}

func UserRepo() biz.UserRepo {
	_, d, _ := Data()
	return data.NewUserRepo(d)
}

func RoleRepo() biz.RoleRepo {
	_, d, _ := Data()
	return data.NewRoleRepo(d)
}

func PermissionRepo() biz.PermissionRepo {
	_, d, _ := Data()
	return data.NewPermissionRepo(d)
}
func ActionRepo() biz.ActionRepo {
	_, d, _ := Data()
	return data.NewActionRepo(d)
}
func UserGroupRepo() biz.UserGroupRepo {
	_, d, _ := Data()
	return data.NewUserGroupRepo(d)
}
func WhitelistRepo() biz.WhitelistRepo {
	_, d, _ := Data()
	return data.NewWhitelistRepo(d)
}
func HotspotRepo() biz.HotspotRepo {
	_, d, _ := Data()
	return data.NewHotspotRepo(d)
}
func HealthRepo() biz.HealthRepo {
	_, d, _ := Data()
	return data.NewHealthRepo(d, onceRedis)
}

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string        { return http.Header(hc).Get(key) }
func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }
func (hc headerCarrier) Add(key string, value string) { http.Header(hc).Add(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice value associated with the passed key.
func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}

func newUserHeader(k, v string) *headerCarrier {
	header := &headerCarrier{}
	header.Set(k, v)
	return header
}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	reqHeader transport.Header
}

func (tr *Transport) Kind() transport.Kind            { return tr.kind }
func (tr *Transport) Endpoint() string                { return tr.endpoint }
func (tr *Transport) Operation() string               { return tr.operation }
func (tr *Transport) RequestHeader() transport.Header { return tr.reqHeader }
func (*Transport) ReplyHeader() transport.Header      { return nil }

// NewContextWithUserId creates a context containing tenant/user id for unit tests.
func NewContextWithUserId(ctx context.Context, u string) context.Context {
	tr := &Transport{
		// Align with tenant v2 middleware: it reads tenant id from "x-tenant-id".
		reqHeader: newUserHeader("x-tenant-id", u),
	}
	ctx = transport.NewServerContext(ctx, tr)
	// In unit tests we call service methods directly (no middleware chain), so set tenant id on context.
	return gormTenant.NewContext(ctx, u)
}

var (
	onceC     *conf.Bootstrap
	onceData  *data.Data
	onceCache biz.Cache
	onceRedis redis.UniversalClient
	once      sync.Once
)

func Data() (c *conf.Bootstrap, dataData *data.Data, cache biz.Cache) {
	debug.SetGCPercent(-1)
	once.Do(func() {
		onceC = DBAndRedis()
		onceData, _, _ = data.NewData(onceC)
		universalClient, err := data.NewRedis(onceC)
		if err != nil {
			panic(err)
		}
		onceRedis = universalClient
		onceCache = data.NewCache(onceC, universalClient)
	})
	return onceC, onceData, onceCache
}

// DBAndRedis returns a Bootstrap config for local integration/unit tests.
func DBAndRedis() *conf.Bootstrap {
	// Use real database connection; use in-memory Redis.
	host1 := "localhost"
	dbName := "auth"
	port1 := 5432
	dbDriver := "postgres"
	dbDsn := fmt.Sprintf("host=%s user=root password=password dbname=%s port=%d sslmode=disable TimeZone=UTC", host1, dbName, port1)

	host2, port2, err := mock.NewRedis()
	if err != nil {
		panic(err)
	}

	return &conf.Bootstrap{
		Server: &conf.Server{
			MachineId: "123",
		},
		Log: &conf.Log{
			Level:   "debug",
			JSON:    false,
			ShowSQL: true,
		},
		Db: &conf.DB{
			Driver:  dbDriver,
			Dsn:     dbDsn,
			Migrate: false,
		},
		Redis: &conf.Redis{
			Dsn: fmt.Sprintf("redis://%s", net.JoinHostPort(host2, fmt.Sprintf("%d", port2))),
		},
		Tracer: &conf.Tracer{
			Enable: true,
			Otlp:   &conf.Tracer_Otlp{},
			Stdout: &conf.Tracer_Stdout{},
		},
	}
}
