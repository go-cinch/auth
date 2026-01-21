package data

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/bsm/redislock"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/plugins/gorm/tenant/v2"

	"auth/internal/biz"
	"auth/internal/conf"
)

// Cache coordinates Redis, redis-lock, and local cache state.
type Cache struct {
	redis   redis.UniversalClient
	locker  *redislock.Client
	disable bool
	prefix  string
	lock    string
	val     string
	refresh bool
}

var local = cache.New(30*time.Minute, 60*time.Minute)

// NewCache creates a Cache using the provided configuration and Redis client.
func NewCache(c *conf.Bootstrap, client redis.UniversalClient) biz.Cache {
	return &Cache{
		redis:   client,
		locker:  redislock.New(client),
		disable: c.Server.Nocache,
		lock:    "lock",
		val:     "val",
	}
}

func (c *Cache) Cache() redis.UniversalClient {
	return c.redis
}

func (c *Cache) WithPrefix(prefix string) biz.Cache {
	return &Cache{
		redis:   c.redis,
		locker:  c.locker,
		disable: c.disable,
		prefix:  prefix,
		lock:    c.lock,
		val:     c.val,
	}
}

func (c *Cache) WithRefresh() biz.Cache {
	return &Cache{
		redis:   c.redis,
		locker:  c.locker,
		disable: c.disable,
		prefix:  c.prefix,
		lock:    c.lock,
		val:     c.val,
		refresh: true,
	}
}

func (c *Cache) Get(
	ctx context.Context,
	action string,
	write func(context.Context) (string, error),
) (res string, err error) {
	tr := otel.Tracer("cache")
	ctx, span := tr.Start(ctx, "Get")
	defer span.End()
	if c.disable {
		return write(ctx)
	}
	key := c.getValKey(ctx, action)
	if !c.refresh {
		// 1. first get cache
		res, err = c.tryGetFromCaches(ctx, key, span)
		if err == nil {
			return res, err
		}
	}
	// 2. get lock before read db
	lock, err := c.Lock(ctx, action)
	if err != nil {
		// Soft wait and read-back once to reduce tail errors on hot keys.
		return c.softWaitReadBack(ctx, key, span)
	}
	span.AddEvent("lock_acquired")
	renewStop := make(chan struct{})
	go c.renewLockPeriodically(ctx, lock, renewStop)
	defer func() {
		close(renewStop)
		_ = lock.Release(ctx)
	}()
	if !c.refresh {
		// 3. double check cache exists (avoid concurrency step 1 ok=false)
		res, err = c.tryGetFromCaches(ctx, key, span)
		if err == nil {
			return res, err
		}
	}
	// 4. load data from db and write to cache
	if write != nil {
		start := time.Now()
		res, err = write(ctx)
		span.SetAttributes(attribute.Int64("db_load_ms", time.Since(start).Milliseconds()))
	}
	return res, err
}

// tryGetFromCaches attempts to get value from local cache first, then Redis.
func (c *Cache) tryGetFromCaches(ctx context.Context, key string, span trace.Span) (string, error) {
	// 1.1. get from local
	res, err := GetFromLocal(key)
	if err == nil {
		span.AddEvent("local_hit_prelock")
		return res, nil
	}
	span.AddEvent("local_miss")
	// 1.2. get from redis
	res, err = c.redis.Get(ctx, key).Result()
	if err == nil {
		c.cacheRedisValueToLocal(ctx, key, res)
		span.AddEvent("redis_hit_prelock")
		return res, nil
	}
	span.AddEvent("redis_miss_prelock")
	return "", err
}

// cacheRedisValueToLocal saves Redis value to local cache with TTL.
func (c *Cache) cacheRedisValueToLocal(ctx context.Context, key, val string) {
	ttl, err := c.redis.TTL(ctx, key).Result()
	if err == nil && ttl > 0 {
		sec := int(ttl.Seconds())
		if sec < 1 { // avoid 0s causing local cache to use default expiration
			sec = 1
		}
		Set2Local(key, val, sec)
	}
}

// softWaitReadBack waits briefly and retries cache reads to reduce lock contention errors.
func (c *Cache) softWaitReadBack(ctx context.Context, key string, span trace.Span) (string, error) {
	span.AddEvent("lock_acquire_failed")
	if c.refresh {
		return "", biz.ErrTooManyRequests(ctx)
	}
	// 15-30ms with jitter
	wait := time.Duration(15+rand.New(rand.NewSource(time.Now().UnixNano())).Intn(16)) * time.Millisecond
	time.Sleep(wait)
	span.AddEvent("soft_wait_read_back", trace.WithAttributes(attribute.Int("wait_ms", int(wait/time.Millisecond))))
	// try read again from caches
	res, err := GetFromLocal(key)
	if err == nil {
		span.AddEvent("soft_read_back_local_hit")
		return res, nil
	}
	res, err = c.redis.Get(ctx, key).Result()
	if err == nil {
		c.cacheRedisValueToLocal(ctx, key, res)
		span.AddEvent("soft_read_back_redis_hit")
		return res, nil
	}
	return "", biz.ErrTooManyRequests(ctx)
}

// renewLockPeriodically renews distributed lock every 10s until stopped.
func (c *Cache) renewLockPeriodically(ctx context.Context, lock *redislock.Lock, stop chan struct{}) {
	// Create an independent span for lock renewal to avoid using parent span that may have ended.
	renewTr := otel.Tracer("cache")
	parentSpanCtx := trace.SpanContextFromContext(ctx)
	_, renewSpan := renewTr.Start(context.Background(), "LockRenewal",
		trace.WithLinks(trace.Link{SpanContext: parentSpanCtx}),
	)
	defer renewSpan.End()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	attempt := 0
	for {
		select {
		case <-ticker.C:
			attempt++
			renewSpan.AddEvent("lock_refresh_start", trace.WithAttributes(
				attribute.Int("attempt", attempt),
			))
			start := time.Now()
			err := lock.Refresh(ctx, 20*time.Second, nil)
			dur := time.Since(start).Milliseconds()
			if err != nil {
				renewSpan.AddEvent("lock_refresh_end", trace.WithAttributes(
					attribute.Int("attempt", attempt),
					attribute.Int("ok", 0),
					attribute.Int64("duration_ms", dur),
					attribute.String("error", err.Error()),
				))
			} else {
				renewSpan.AddEvent("lock_refresh_end", trace.WithAttributes(
					attribute.Int("attempt", attempt),
					attribute.Int("ok", 1),
					attribute.Int64("duration_ms", dur),
				))
			}
		case <-stop:
			return
		}
	}
}

func (c *Cache) Set(ctx context.Context, action, data string, short bool) {
	// set random expiration avoid a large number of keys expire at the same time
	// normal cache: 300-599s; negative cache (short): 45-75s
	seconds := randSeconds(300, 599)
	if short {
		// negative cache TTL with jitter to avoid synchronized expirations
		seconds = randSeconds(45, 75)
	}
	c.SetWithExpiration(ctx, action, data, seconds)
}

func (c *Cache) SetWithExpiration(ctx context.Context, action, data string, seconds int64) {
	if c.disable {
		return
	}
	key := c.getValKey(ctx, action)
	// set to local cache
	Set2Local(key, data, int(seconds))
	// set to redis
	err := c.redis.Set(ctx, key, data, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action":  action,
				"seconds": seconds,
			}).
			Warn("set cache failed")
		return
	}
}

func (c *Cache) Del(ctx context.Context, action string) {
	if c.disable {
		return
	}
	key := c.getValKey(ctx, action)
	DelFromLocal(key)
	err := c.redis.Unlink(ctx, key).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action": action,
				"key":    key,
			}).
			Warn("del cache failed")
	}
}

func (c *Cache) Flush(ctx context.Context, handler func(ctx context.Context) error) (err error) {
	tr := otel.Tracer("cache")
	ctx, span := tr.Start(ctx, "Flush")
	defer span.End()

	err = handler(ctx)
	if err != nil {
		return
	}
	if c.disable {
		return
	}
	pattern := c.getPrefixKey(ctx, "*")
	if errDel := c.deleteByScan(ctx, pattern); errDel != nil {
		log.
			WithContext(ctx).
			WithError(errDel).
			WithFields(log.Fields{
				"pattern": pattern,
			}).
			Warn("flush cache failed")
	}
	return
}

func (c *Cache) FlushByPrefix(ctx context.Context, prefix ...string) (err error) {
	tr := otel.Tracer("cache")
	ctx, span := tr.Start(ctx, "FlushByPrefix")
	defer span.End()

	pattern := c.getPrefixKey(ctx, prefix...)
	if errDel := c.deleteByScan(ctx, pattern); errDel != nil {
		log.
			WithContext(ctx).
			WithError(errDel).
			WithFields(log.Fields{
				"pattern": pattern,
			}).
			Warn("flush cache by prefix failed")
	}
	return
}

func (c *Cache) deleteByScan(ctx context.Context, pattern string) error {
	switch rc := c.redis.(type) {
	case *redis.ClusterClient:
		return rc.ForEachMaster(ctx, func(ctx context.Context, shard *redis.Client) error {
			return c.scanDeleteOnUniversal(ctx, shard, pattern)
		})
	case *redis.Ring:
		return rc.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
			return c.scanDeleteOnUniversal(ctx, shard, pattern)
		})
	default:
		return c.scanDeleteOnUniversal(ctx, c.redis, pattern)
	}
}

func (c *Cache) scanDeleteOnUniversal(ctx context.Context, client redis.UniversalClient, pattern string) error {
	var cursor uint64
	deletedTotal := 0
	for {
		keys, cur, err := client.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			return err
		}
		cursor = cur
		deleted, err := c.batchDeleteKeys(ctx, client, keys)
		if err != nil {
			return err
		}
		deletedTotal += deleted
		if cursor == 0 {
			break
		}
	}
	// record per-shard deleted count
	trace.SpanFromContext(ctx).AddEvent("flush_delete_shard", trace.WithAttributes(
		attribute.Int("deleted_keys", deletedTotal),
		attribute.String("pattern", pattern),
	))
	return nil
}

// batchDeleteKeys deletes keys in batch, skipping lock keys for safety.
func (c *Cache) batchDeleteKeys(ctx context.Context, client redis.UniversalClient, keys []string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	p := client.Pipeline()
	batch := 0
	for _, k := range keys {
		// extra safety: skip lock keys even if matched
		if strings.Contains(k, "_lock_") {
			continue
		}
		DelFromLocal(k)
		p.Unlink(ctx, k)
		batch++
	}
	if batch > 0 {
		if _, err := p.Exec(ctx); err != nil {
			return 0, err
		}
	}
	return batch, nil
}

func (c *Cache) Lock(ctx context.Context, action string) (*redislock.Lock, error) {
	tr := otel.Tracer("cache")
	ctx, span := tr.Start(ctx, "Lock")
	defer span.End()
	lock, err := c.locker.Obtain(
		ctx,
		c.getLockKey(ctx, action),
		20*time.Second,
		&redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(5*time.Millisecond), 400),
		},
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return lock, err
}

// Set2Local caches the value in the in-memory store for the given TTL.
func Set2Local(key, val string, expire int) {
	local.Set(key, val, time.Duration(expire)*time.Second)
}

// GetFromLocal returns cached data from the in-memory store.
func GetFromLocal(key string) (string, error) {
	val, ok := local.Get(key)
	if !ok {
		return "", errors.New("key not found")
	}
	return val.(string), nil
}

// DelFromLocal removes the cached key from the in-memory store.
func DelFromLocal(key string) {
	local.Delete(key)
}

// randSeconds returns a random integer seconds in [minVal, maxVal] inclusive.
// It is used to add TTL jitter to avoid synchronized expirations.
func randSeconds(minVal, maxVal int64) int64 {
	if maxVal <= minVal {
		return minVal
	}
	return minVal + rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(maxVal-minVal+1)
}

func (c *Cache) getPrefixKey(ctx context.Context, arr ...string) string {
	id := tenant.FromContext(ctx)

	// Build base prefix; when additional segments provided, include value namespace (c.val).
	var prefix string
	if len(arr) > 0 {
		// append params prefix need add val
		segs := append([]string{c.prefix, c.val}, arr...)
		prefix = strings.Join(segs, "_")
	} else {
		prefix = c.prefix
	}

	if strings.TrimSpace(prefix) == "" {
		// avoid flush all key
		log.WithContext(ctx).Warn("invalid prefix")
		prefix = "prefix"
	}

	// ensure exactly one wildcard at the end
	if !strings.HasSuffix(prefix, "*") {
		prefix = prefix + "*"
	}

	if id == "" {
		return prefix
	}
	return id + "_" + prefix
}

func (c *Cache) getValKey(ctx context.Context, action string) string {
	id := tenant.FromContext(ctx)
	if id == "" {
		return strings.Join([]string{c.prefix, c.val, action}, "_")
	}
	return strings.Join([]string{id, c.prefix, c.val, action}, "_")
}

func (c *Cache) getLockKey(ctx context.Context, action string) string {
	id := tenant.FromContext(ctx)
	if id == "" {
		return strings.Join([]string{c.prefix, c.lock, action}, "_")
	}
	return strings.Join([]string{id, c.prefix, c.lock, action}, "_")
}
