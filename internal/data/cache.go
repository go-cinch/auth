package data

import (
	"auth/internal/biz"
	"context"
	"fmt"
	"github.com/go-cinch/common/log"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

// Cache .
type Cache struct {
	redis  redis.UniversalClient
	prefix string
	lock   string
	val    string
}

func (c *Cache) Cache() redis.UniversalClient {
	return c.redis
}

func (c *Cache) WithPrefix(prefix string) biz.Cache {
	return &Cache{
		redis:  c.redis,
		prefix: prefix,
		lock:   fmt.Sprintf("%s_%s", prefix, c.lock),
		val:    fmt.Sprintf("%s_%s", prefix, c.val),
	}
}

func (c *Cache) Get(ctx context.Context, action string, write func(context.Context) (string, bool)) (res string, ok bool, lock bool, db bool) {
	var err error
	// 1. first get cache
	res, err = c.redis.Get(ctx, fmt.Sprintf("%s_%s", c.val, action)).Result()
	if err == nil {
		// cache exists
		ok = true
		return
	}
	// 2. get lock before read db
	lock = c.Lock(ctx, action)
	if !lock {
		return
	}
	defer c.Unlock(ctx, action)
	// 3. double check cache exists(avoid concurrency step 1 ok=false)
	res, err = c.redis.Get(ctx, fmt.Sprintf("%s_%s", c.val, action)).Result()
	if err == nil {
		// cache exists
		ok = true
		return
	}
	if write != nil {
		res, ok = write(ctx)
		db = true
	}
	return
}

func (c *Cache) Set(ctx context.Context, action, data string, short bool) {
	// set random expiration avoid a large number of keys expire at the same time
	seconds := rand.New(rand.NewSource(time.Now().Unix())).Int63n(300) + 300
	if short {
		// if record not found, set a short expiration
		seconds = 60
	}
	c.SetWithExpiration(ctx, action, data, seconds)
}

func (c *Cache) SetWithExpiration(ctx context.Context, action, data string, seconds int64) {
	// set random expiration avoid a large number of keys expire at the same time
	err := c.redis.Set(ctx, fmt.Sprintf("%s_%s", c.val, action), data, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action":  action,
				"seconds": seconds,
			}).
			Warn("set cache failed")
	}
}

func (c *Cache) Del(ctx context.Context, action string) {
	err := c.redis.Del(ctx, fmt.Sprintf("%s_%s", c.val, action)).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action": action,
			}).
			Warn("del cache failed")
	}
}

func (c *Cache) Flush(ctx context.Context, handler func(ctx context.Context) error) (err error) {
	err = handler(ctx)
	if err != nil {
		return
	}
	action := c.prefix + "*"
	arr := c.redis.Keys(ctx, action).Val()
	p := c.redis.Pipeline()
	for _, item := range arr {
		if item == c.lock {
			continue
		}
		p.Del(ctx, item)
	}
	_, pErr := p.Exec(ctx)
	if pErr != nil {
		log.
			WithContext(ctx).
			WithError(pErr).
			WithFields(log.Fields{
				"action": action,
			}).
			Warn("flush cache failed")
	}
	return
}

func (c *Cache) Lock(ctx context.Context, action string) (ok bool) {
	ok, _ = c.redis.SetNX(ctx, fmt.Sprintf("%s_%s", c.lock, action), 1, time.Minute).Result()
	retry := 0
	for retry < 20 && !ok {
		time.Sleep(25 * time.Millisecond)
		ok, _ = c.redis.SetNX(ctx, fmt.Sprintf("%s_%s", c.lock, action), 1, time.Minute).Result()
		retry++
	}
	return
}

func (c *Cache) Unlock(ctx context.Context, action string) {
	err := c.redis.Del(ctx, fmt.Sprintf("%s_%s", c.lock, action)).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action": action,
			}).
			Warn("unlock cache failed")
	}
}

// NewCache .
func NewCache(client redis.UniversalClient) biz.Cache {
	return &Cache{
		redis:  client,
		prefix: "",
		lock:   "lock",
		val:    "val",
	}
}
