package data

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/log"
	glog "github.com/go-cinch/common/plugins/gorm/log"
	"github.com/go-cinch/common/plugins/gorm/tenant/v2"
	"github.com/go-cinch/common/utils"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/db"
)

// Data wraps all data sources used by the service.
type Data struct {
	Tenant    *tenant.Tenant
	sonyflake *id.Sonyflake
}

// NewData initializes the configured database connection via tenant v2.
func NewData(c *conf.Bootstrap) (*Data, func(), error) {
	gormTenant, err := NewDB(c)
	if err != nil {
		return nil, nil, err
	}

	sonyflake, err := NewSonyflake(c)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.Info("closing database connections")
	}

	return &Data{
		Tenant:    gormTenant,
		sonyflake: sonyflake,
	}, cleanup, nil
}

// NewDB initializes tenant-aware database connection using tenant v2 package.
// Supports both MySQL and PostgreSQL via internal driver detection.
func NewDB(c *conf.Bootstrap) (*tenant.Tenant, error) {
	dbConf := c.Db
	driver := strings.ToLower(strings.TrimSpace(dbConf.Driver))
	dsn := strings.TrimSpace(dbConf.Dsn)
	if driver == "" {
		err := errors.New("db driver is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}
	if dsn == "" {
		err := errors.New("db DSN is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}

	level := log.NewLevel(c.Log.Level)
	// Force to warn level when show sql is false.
	if level > log.WarnLevel && !c.Log.ShowSQL {
		level = log.WarnLevel
	}

	ops := []func(*tenant.Options){
		tenant.WithDriver(driver),
		tenant.WithDSN("", dsn), // Empty string for default tenant
		tenant.WithSQLFile(db.SQLFiles),
		tenant.WithSQLRoot(db.SQLRoot),
		tenant.WithSkipMigrate(!dbConf.Migrate), // Skip migration if Migrate is false
		tenant.WithConfig(&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			QueryFields: true,
			Logger: glog.New(
				glog.WithColorful(false),
				glog.WithSlow(200),
				glog.WithLevel(level),
			),
		}),
		tenant.WithMaxIdle(10),
		tenant.WithMaxOpen(100),
	}

	gormTenant, err := tenant.New(ops...)
	if err != nil {
		log.WithError(err).Error("create tenant failed")
		return nil, err
	}

	// Always call Migrate() to initialize the database connection.
	// WithSkipMigrate controls whether SQL migrations are actually executed.
	if err := gormTenant.Migrate(); err != nil {
		log.WithError(err).Error("migrate tenant failed")
		return nil, err
	}

	log.Info("initialize db success, driver: %s", driver)
	return gormTenant, nil
}

type contextTxKey struct{}

// Tx is transaction wrapper.
func (d *Data) Tx(ctx context.Context, handler func(ctx context.Context) error) error {
	return d.Tenant.DB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return handler(ctx)
	})
}

// DB returns a tenant-aware GORM DB instance from context.
// If a transaction is present in the context, it returns the transaction DB.
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.Tenant.DB(ctx)
}

// NewTransaction creates a new Transaction from Data.
func NewTransaction(d *Data) biz.Transaction {
	return d
}

// ID generates a unique distributed ID using Sonyflake.
func (d *Data) ID(ctx context.Context) int64 {
	return int64(d.sonyflake.ID(ctx))
}

// NewSonyflake initializes the Sonyflake ID generator.
func NewSonyflake(c *conf.Bootstrap) (*id.Sonyflake, error) {
	machineID, _ := strconv.ParseUint(c.Server.MachineId, 10, 16)
	sf := id.NewSonyflake(
		id.WithSonyflakeMachineID(uint16(machineID)),
		id.WithSonyflakeStartTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	if sf.Error != nil {
		log.WithError(sf.Error).Error("initialize sonyflake failed")
		return nil, errors.New("initialize sonyflake failed")
	}
	log.
		WithField("machine.id", machineID).
		Info("initialize sonyflake success")
	return sf, nil
}

// NewRedis initializes Redis client from config.
func NewRedis(c *conf.Bootstrap) (redis.UniversalClient, error) {
	return newRedis(c)
}

// newRedis is the shared Redis initialization logic.
func newRedis(c *conf.Bootstrap) (client redis.UniversalClient, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var u *url.URL
	u, err = url.Parse(c.Redis.Dsn)
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "***")
	}
	showDsn, _ := url.PathUnescape(u.String())
	client, err = utils.ParseRedisURI(c.Redis.Dsn)
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	log.
		WithField("redis.dsn", showDsn).
		Info("initialize redis success")
	return
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewSonyflake,
	NewRedis,
	NewTracer,
	NewCache,
	NewTransaction,
	NewHealthRepo,
	NewAuthRepo,
	NewUserRepo,
	NewRoleRepo,
	NewPermissionRepo,
	NewActionRepo,
	NewUserGroupRepo,
	NewWhitelistRepo,
	NewHotspotRepo,
)
