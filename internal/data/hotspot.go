package data

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data/model"
	"auth/internal/data/query"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type hotspotRepo struct {
	c    *conf.Bootstrap
	data *Data
}

func NewHotspotRepo(c *conf.Bootstrap, data *Data) biz.HotspotRepo {
	return &hotspotRepo{
		c:    c,
		data: data,
	}
}

func (ro hotspotRepo) Refresh(ctx context.Context) (err error) {
	pipe := ro.data.redis.Pipeline()
	ro.refreshUserGroup(ctx, pipe)
	ro.refreshUserUserGroupRelation(ctx, pipe)
	ro.refreshRole(ctx, pipe)
	ro.refreshAction(ctx, pipe)
	ro.refreshWhitelist(ctx, pipe)
	ro.refreshUser(ctx, pipe)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.
			WithContext(ctx).
			Warn("refresh failed: %v", err)
	}
	return
}

func (ro hotspotRepo) GetUserByCode(ctx context.Context, code string) *biz.User {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetUserByCode")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).User
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.Code.ColumnName().String(),
		code,
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	m := make(map[string]interface{}, len(res))
	copierx.Copy(&m, res)
	if v, ok := m[utils.CamelCase(p.Locked.ColumnName().String())]; ok {
		m[utils.CamelCase(p.Locked.ColumnName().String())], _ = strconv.ParseBool(v.(string))
	}
	if v, ok := m[utils.CamelCase(p.Wrong.ColumnName().String())]; ok {
		m[utils.CamelCase(p.Wrong.ColumnName().String())], _ = strconv.ParseUint(v.(string), 10, 64)
	}
	var item biz.User
	utils.Struct2StructByJson(&item, m)
	if item.RoleId > constant.UI0 {
		item.Role = *ro.GetRoleByID(ctx, item.RoleId)
	}
	span.SetAttributes(
		attribute.String("code", code),
		attribute.String("id", strconv.FormatUint(item.Id, 10)),
	)
	return &item
}

func (ro hotspotRepo) GetUserByUsername(ctx context.Context, username string) *biz.User {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetUserByUsername")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).User
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.Username.ColumnName().String(),
		username,
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	m := make(map[string]interface{}, len(res))
	copierx.Copy(&m, res)
	if v, ok := m[utils.CamelCase(p.Locked.ColumnName().String())]; ok {
		m[utils.CamelCase(p.Locked.ColumnName().String())], _ = strconv.ParseBool(v.(string))
	}
	if v, ok := m[utils.CamelCase(p.Wrong.ColumnName().String())]; ok {
		m[utils.CamelCase(p.Wrong.ColumnName().String())], _ = strconv.ParseUint(v.(string), 10, 64)
	}
	var item biz.User
	utils.Struct2StructByJson(&item, m)
	span.SetAttributes(
		attribute.String("id", strconv.FormatUint(item.Id, 10)),
		attribute.String("username", item.Username),
	)
	return &item
}

func (ro hotspotRepo) GetRoleByID(ctx context.Context, id uint64) *biz.Role {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetRoleByID")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).Role
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.ID.ColumnName().String(),
		strconv.FormatUint(id, 10),
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	var item biz.Role
	utils.Struct2StructByJson(&item, res)
	span.SetAttributes(
		attribute.String("id", strconv.FormatUint(item.Id, 10)),
		attribute.String("word", item.Word),
	)
	return &item
}

func (ro hotspotRepo) GetActionByWord(ctx context.Context, word string) *biz.Action {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetActionByWord")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).Action
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.Word.ColumnName().String(),
		word,
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	var item biz.Action
	utils.Struct2StructByJson(&item, res)
	span.SetAttributes(
		attribute.String("id", strconv.FormatUint(item.Id, 10)),
		attribute.String("word", item.Word),
	)
	return &item
}

func (ro hotspotRepo) GetActionByCode(ctx context.Context, code string) *biz.Action {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetActionByCode")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).Action
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.Code.ColumnName().String(),
		code,
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	var item biz.Action
	utils.Struct2StructByJson(&item, res)
	span.SetAttributes(
		attribute.String("id", strconv.FormatUint(item.Id, 10)),
		attribute.String("code", item.Code),
	)
	return &item
}

func (ro hotspotRepo) FindActionByCode(ctx context.Context, codes ...string) (list []biz.Action) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "FindActionByCode")
	defer span.End()
	for _, code := range codes {
		list = append(list, *ro.GetActionByCode(ctx, code))
	}
	return
}

func (ro hotspotRepo) FindUserGroupByUserCode(ctx context.Context, code string) (list []biz.UserGroup) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "FindUserGroupByUserCode")
	defer span.End()
	list = make([]biz.UserGroup, 0)
	user := ro.GetUserByCode(ctx, code)
	return ro.FindUserGroupByUserID(ctx, user.Id)
}

func (ro hotspotRepo) FindUserGroupByUserID(ctx context.Context, id uint64) (list []biz.UserGroup) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "FindUserGroupIDByUserID")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).UserUserGroupRelation
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.UserID.ColumnName().String(),
		strconv.FormatUint(id, 10),
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	groupKeys := lo.Keys(res)
	groupIDs := lo.Map(groupKeys, func(item string, index int) uint64 {
		num, _ := strconv.ParseUint(item, 10, 64)
		return num
	})
	for _, groupID := range groupIDs {
		list = append(list, *ro.GetUserGroupByID(ctx, groupID))
	}
	span.SetAttributes(
		attribute.String("userID", strconv.FormatUint(id, 10)),
		attribute.StringSlice("groups", groupKeys),
	)
	return
}

func (ro hotspotRepo) GetUserGroupByID(ctx context.Context, id uint64) *biz.UserGroup {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetUserGroupByID")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).UserGroup
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.ID.ColumnName().String(),
		strconv.FormatUint(id, 10),
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	var item biz.UserGroup
	utils.Struct2StructByJson(&item, res)
	span.SetAttributes(
		attribute.String("id", strconv.FormatUint(id, 10)),
		attribute.String("word", item.Word),
	)
	return &item
}

func (ro hotspotRepo) FindWhitelistResourceByCategory(ctx context.Context, category uint32) []string {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "FindUserGroupIDByUserID")
	defer span.End()
	rds := ro.data.redis
	p := query.Use(ro.data.DB(ctx)).Whitelist
	key := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		p.TableName(),
		p.Category.ColumnName().String(),
		strconv.FormatUint(uint64(category), 10),
	}, ".")
	res, _ := rds.HGetAll(ctx, key).Result()
	resourceKeys := lo.Keys(res)
	resources := lo.Map(resourceKeys, func(item string, index int) string {
		return item
	})
	span.SetAttributes(
		attribute.String("category", strconv.FormatUint(uint64(category), 10)),
		attribute.StringSlice("resources", resourceKeys),
	)
	return resources
}

func (ro hotspotRepo) refreshUserGroup(ctx context.Context, pipe redis.Pipeliner) {
	p := query.Use(ro.data.DB(ctx)).UserGroup
	db := p.WithContext(ctx)
	name := p.TableName()
	list, _ := db.Find()
	matchKey := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		name,
		"*",
	}, ".")
	matches := scanKeys(ctx, ro.data.redis, matchKey)
	for _, key := range matches {
		pipe.Del(ctx, key)
	}
	for _, item := range list {
		idKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.ID.ColumnName().String(),
			strconv.FormatUint(item.ID, 10),
		}, ".")
		pipe.Del(ctx, idKey)
		pipe.HSet(
			ctx, idKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
		)
		pipe.Expire(ctx, idKey, ro.randomExpire())
		wordKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Word.ColumnName().String(),
			item.Word,
		}, ".")
		pipe.Del(ctx, wordKey)
		pipe.HSet(
			ctx, wordKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
		)
		pipe.Expire(ctx, wordKey, ro.randomExpire())
	}
	return
}

func (ro hotspotRepo) refreshUserUserGroupRelation(ctx context.Context, pipe redis.Pipeliner) {
	p := query.Use(ro.data.DB(ctx)).UserUserGroupRelation
	db := p.WithContext(ctx)
	name := p.TableName()
	list, _ := db.Find()
	matchKey := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		name,
		"*",
	}, ".")
	matches := scanKeys(ctx, ro.data.redis, matchKey)
	for _, key := range matches {
		pipe.Del(ctx, key)
	}
	// group id by user id
	groupIDMap := lo.MapValues(
		lo.GroupBy(list, func(item *model.UserUserGroupRelation) uint64 {
			return item.UserID
		}),
		func(items []*model.UserUserGroupRelation, _ uint64) []uint64 {
			return lo.Map(items, func(item *model.UserUserGroupRelation, _ int) uint64 {
				return item.UserGroupID
			})
		})
	for userID, groupIDs := range groupIDMap {
		userIDKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.UserID.ColumnName().String(),
			strconv.FormatUint(userID, 10),
		}, ".")
		pipe.Del(ctx, userIDKey)
		for _, id := range groupIDs {
			pipe.HSet(
				ctx, userIDKey,
				strconv.FormatUint(id, 10), "",
			)
		}
		pipe.Expire(ctx, userIDKey, ro.randomExpire())
	}
	return
}

func (ro hotspotRepo) refreshRole(ctx context.Context, pipe redis.Pipeliner) {
	p := query.Use(ro.data.DB(ctx)).Role
	db := p.WithContext(ctx)
	name := p.TableName()
	list, _ := db.Find()
	matchKey := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		name,
		"*",
	}, ".")
	matches := scanKeys(ctx, ro.data.redis, matchKey)
	for _, key := range matches {
		pipe.Del(ctx, key)
	}
	for _, item := range list {
		idKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.ID.ColumnName().String(),
			strconv.FormatUint(item.ID, 10),
		}, ".")
		pipe.Del(ctx, idKey)
		pipe.HSet(
			ctx, idKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
		)
		pipe.Expire(ctx, idKey, ro.randomExpire())
		wordKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Word.ColumnName().String(),
			item.Word,
		}, ".")
		pipe.Del(ctx, wordKey)
		pipe.HSet(
			ctx, wordKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
		)
		pipe.Expire(ctx, wordKey, ro.randomExpire())
	}
	return
}

func (ro hotspotRepo) refreshAction(ctx context.Context, pipe redis.Pipeliner) {
	p := query.Use(ro.data.DB(ctx)).Action
	db := p.WithContext(ctx)
	name := p.TableName()
	list, _ := db.Find()
	matchKey := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		name,
		"*",
	}, ".")
	matches := scanKeys(ctx, ro.data.redis, matchKey)
	for _, key := range matches {
		pipe.Del(ctx, key)
	}
	for _, item := range list {
		idKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.ID.ColumnName().String(),
			strconv.FormatUint(item.ID, 10),
		}, ".")
		pipe.Del(ctx, idKey)
		pipe.HSet(
			ctx, idKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Code.ColumnName().String()), item.Code,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Resource.ColumnName().String()), item.Resource,
			utils.CamelCase(p.Menu.ColumnName().String()), item.Menu,
			utils.CamelCase(p.Btn.ColumnName().String()), item.Btn,
		)
		pipe.Expire(ctx, idKey, ro.randomExpire())
		codeKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Code.ColumnName().String(),
			item.Code,
		}, ".")
		pipe.Del(ctx, codeKey)
		pipe.HSet(
			ctx, codeKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Code.ColumnName().String()), item.Code,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Resource.ColumnName().String()), item.Resource,
			utils.CamelCase(p.Menu.ColumnName().String()), item.Menu,
			utils.CamelCase(p.Btn.ColumnName().String()), item.Btn,
		)
		pipe.Expire(ctx, codeKey, ro.randomExpire())
		wordKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Word.ColumnName().String(),
			item.Word,
		}, ".")
		pipe.Del(ctx, wordKey)
		pipe.HSet(
			ctx, wordKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Name.ColumnName().String()), item.Name,
			utils.CamelCase(p.Code.ColumnName().String()), item.Code,
			utils.CamelCase(p.Word.ColumnName().String()), item.Word,
			utils.CamelCase(p.Resource.ColumnName().String()), item.Resource,
			utils.CamelCase(p.Menu.ColumnName().String()), item.Menu,
			utils.CamelCase(p.Btn.ColumnName().String()), item.Btn,
		)
		pipe.Expire(ctx, wordKey, ro.randomExpire())
	}
	return
}

func (ro hotspotRepo) refreshWhitelist(ctx context.Context, pipe redis.Pipeliner) {
	p := query.Use(ro.data.DB(ctx)).Whitelist
	db := p.WithContext(ctx)
	name := p.TableName()
	list, _ := db.Find()
	matchKey := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		name,
		"*",
	}, ".")
	matches := scanKeys(ctx, ro.data.redis, matchKey)
	for _, key := range matches {
		pipe.Del(ctx, key)
	}
	for _, item := range list {
		idKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.ID.ColumnName().String(),
			strconv.FormatUint(item.ID, 10),
		}, ".")
		pipe.Del(ctx, idKey)
		pipe.HSet(
			ctx, idKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Category.ColumnName().String()), strconv.FormatUint(uint64(item.Category), 10),
			utils.CamelCase(p.Resource.ColumnName().String()), item.Resource,
		)
		pipe.Expire(ctx, idKey, ro.randomExpire())
	}

	// resource by category
	resourceMap := lo.MapValues(
		lo.GroupBy(list, func(item *model.Whitelist) uint32 {
			return item.Category
		}),
		func(items []*model.Whitelist, _ uint32) []string {
			return lo.Map(items, func(item *model.Whitelist, _ int) string {
				return item.Resource
			})
		})
	for category, resources := range resourceMap {
		categoryKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Category.ColumnName().String(),
			strconv.FormatUint(uint64(category), 10),
		}, ".")
		pipe.Del(ctx, categoryKey)
		for _, resource := range resources {
			pipe.HSet(
				ctx, categoryKey,
				resource, "",
			)
		}
		pipe.Expire(ctx, categoryKey, ro.randomExpire())
	}
	return
}

func (ro hotspotRepo) refreshUser(ctx context.Context, pipe redis.Pipeliner) {
	p := query.Use(ro.data.DB(ctx)).User
	db := p.WithContext(ctx)
	name := p.TableName()
	list, _ := db.Find()
	matchKey := strings.Join([]string{
		ro.c.Name,
		ro.c.Hotspot.Name,
		name,
		"*",
	}, ".")
	matches := scanKeys(ctx, ro.data.redis, matchKey)
	for _, key := range matches {
		pipe.Del(ctx, key)
	}
	for _, item := range list {
		idKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.ID.ColumnName().String(),
			strconv.FormatUint(item.ID, 10),
		}, ".")
		pipe.Del(ctx, idKey)
		pipe.HSet(
			ctx, idKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Code.ColumnName().String()), item.Code,
			utils.CamelCase(p.RoleID.ColumnName().String()), strconv.FormatUint(item.RoleID, 10),
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
			utils.CamelCase(p.Username.ColumnName().String()), item.Username,
			utils.CamelCase(p.Password.ColumnName().String()), item.Password,
			utils.CamelCase(p.Platform.ColumnName().String()), item.Platform,
			utils.CamelCase(p.Wrong.ColumnName().String()), strconv.FormatUint(item.Wrong, 10),
			utils.CamelCase(p.Locked.ColumnName().String()), item.Locked,
		)
		pipe.Expire(ctx, idKey, ro.randomExpire())
		codeKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Code.ColumnName().String(),
			item.Code,
		}, ".")
		pipe.Del(ctx, codeKey)
		pipe.HSet(
			ctx, codeKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Code.ColumnName().String()), item.Code,
			utils.CamelCase(p.RoleID.ColumnName().String()), strconv.FormatUint(item.RoleID, 10),
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
			utils.CamelCase(p.Username.ColumnName().String()), item.Username,
			utils.CamelCase(p.Password.ColumnName().String()), item.Password,
			utils.CamelCase(p.Platform.ColumnName().String()), item.Platform,
			utils.CamelCase(p.Wrong.ColumnName().String()), strconv.FormatUint(item.Wrong, 10),
			utils.CamelCase(p.Locked.ColumnName().String()), item.Locked,
		)
		pipe.Expire(ctx, codeKey, ro.randomExpire())
		usernameKey := strings.Join([]string{
			ro.c.Name,
			ro.c.Hotspot.Name,
			name,
			p.Username.ColumnName().String(),
			item.Username,
		}, ".")
		pipe.Del(ctx, usernameKey)
		pipe.HSet(
			ctx, usernameKey,
			utils.CamelCase(p.ID.ColumnName().String()), item.ID,
			utils.CamelCase(p.Code.ColumnName().String()), item.Code,
			utils.CamelCase(p.RoleID.ColumnName().String()), strconv.FormatUint(item.RoleID, 10),
			utils.CamelCase(p.Action.ColumnName().String()), item.Action,
			utils.CamelCase(p.Username.ColumnName().String()), item.Username,
			utils.CamelCase(p.Password.ColumnName().String()), item.Password,
			utils.CamelCase(p.Platform.ColumnName().String()), item.Platform,
			utils.CamelCase(p.Wrong.ColumnName().String()), strconv.FormatUint(item.Wrong, 10),
			utils.CamelCase(p.Locked.ColumnName().String()), item.Locked,
		)
		pipe.Expire(ctx, usernameKey, ro.randomExpire())
	}
	return
}

func scanKeys(ctx context.Context, rds redis.UniversalClient, pattern string) (list []string) {
	var cursor uint64
	for {
		keys, cursorNew, err := rds.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return
		}

		list = append(list, keys...)
		cursor = cursorNew

		if cursor == 0 {
			break
		}
	}
	return
}

func (ro hotspotRepo) randomExpire() time.Duration {
	expire := ro.c.Hotspot.Expire.AsDuration()
	seconds := rand.New(rand.NewSource(time.Now().Unix())).Int63n(3600)
	return expire + time.Duration(seconds)*time.Second
}
