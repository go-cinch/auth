package data

import (
	"context"
	"strings"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type userGroupRepo struct {
	data *Data
}

func NewUserGroupRepo(data *Data) biz.UserGroupRepo {
	return &userGroupRepo{data: data}
}

func (ro userGroupRepo) Create(ctx context.Context, item *biz.UserGroup) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	db := gorm.G[model.UserGroup](ro.data.DB(ctx))

	// Check word uniqueness.
	word := strings.TrimSpace(item.Word)
	if word != "" {
		count, err := db.Where("word = ?", word).Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("check user_group word exists failed")
			return err
		}
		if count > 0 {
			return biz.ErrDuplicateField(ctx, "word", word)
		}
	}

	var m model.UserGroup
	if err := copierx.Copy(&m, item); err != nil {
		log.WithContext(ctx).WithError(err).Error("copy user_group failed")
		return err
	}

	if m.ID == 0 {
		m.ID = ro.data.ID(ctx)
	}

	// Validate action codes if provided.
	if action := strings.TrimSpace(item.Action); action != "" {
		if err := ro.actionCodesExist(ctx, action); err != nil {
			return err
		}
	}

	if err := db.Create(ctx, &m); err != nil {
		log.WithContext(ctx).WithError(err).Error("create user_group failed")
		return err
	}

	// Create initial user relations if provided.
	if len(item.Users) > 0 {
		userIDs := make([]int64, 0, len(item.Users))
		for _, u := range item.Users {
			if u.Id == 0 {
				continue
			}
			userIDs = append(userIDs, u.Id)
		}
		if err := ro.AddUsers(ctx, m.ID, userIDs); err != nil {
			return err
		}
	}

	return nil
}

func (ro userGroupRepo) Find(ctx context.Context, condition *biz.FindUserGroup) (rp []biz.UserGroup) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	rp = make([]biz.UserGroup, 0)

	q := gorm.G[model.UserGroup](ro.data.DB(ctx)).Where("1 = 1")

	if condition.Name != nil {
		q = q.Where("name LIKE ?", "%"+*condition.Name+"%")
	}
	if condition.Word != nil {
		q = q.Where("word LIKE ?", "%"+*condition.Word+"%")
	}
	if condition.Action != nil {
		q = q.Where("action LIKE ?", "%"+*condition.Action+"%")
	}

	if !condition.Page.Disable {
		count, err := q.Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("count user_group failed")
			return rp
		}
		condition.Page.Total = count
		if count == 0 {
			return rp
		}
	}

	q = q.Order("id DESC")
	if !condition.Page.Disable {
		limit, offset := condition.Page.Limit()
		q = q.Limit(limit).Offset(offset)
	}

	list, err := q.Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find user_group failed")
		return rp
	}

	copierx.Copy(&rp, list)

	// Populate Actions and Users for each group
	for i := range rp {
		if rp[i].Action != "" {
			rp[i].Actions = ro.getActionsByCode(ctx, rp[i].Action)
		}
		rp[i].Users = ro.getUsersByGroupID(ctx, rp[i].Id)
	}

	return rp
}

func (ro userGroupRepo) Update(ctx context.Context, item *biz.UpdateUserGroup) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	db := gorm.G[model.UserGroup](ro.data.DB(ctx))

	m, err := db.Where("id = ?", item.Id).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return biz.ErrRecordNotFound(ctx)
	}
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("get user_group failed")
		return err
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		return biz.ErrDataNotChange(ctx)
	}

	// Handle relation update by replacing all group users (users is a comma-separated list).
	if a, ok := change["users"]; ok {
		if v, ok := a.(string); ok {
			if err := ro.replaceUsers(ctx, item.Id, utils.Str2Int64Arr(v)); err != nil {
				return err
			}
		}
		delete(change, "users")
	}

	// Check word uniqueness if changed.
	if v, ok := change["word"]; ok {
		if word, ok := v.(string); ok {
			word = strings.TrimSpace(word)
			if word != "" {
				count, err := db.Where("word = ? AND id != ?", word, item.Id).Count(ctx, "*")
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("check user_group word uniqueness failed")
					return err
				}
				if count > 0 {
					return biz.ErrDuplicateField(ctx, "word", word)
				}
			}
		}
	}

	// Validate action codes if changed.
	if v, ok := change["action"]; ok {
		if codes, ok := v.(string); ok {
			if err := ro.actionCodesExist(ctx, codes); err != nil {
				return err
			}
		}
	}

	if len(change) == 0 {
		return nil
	}

	err = ro.data.DB(ctx).
		Model(&model.UserGroup{}).
		Where("id = ?", item.Id).
		Updates(change).
		Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update user_group failed")
	}
	return err
}

func (ro userGroupRepo) Delete(ctx context.Context, ids ...int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	if len(ids) == 0 {
		return nil
	}

	// Remove relation rows first to avoid FK constraint issues.
	if _, err := gorm.G[model.UserUserGroupRelation](ro.data.DB(ctx)).
		Where("user_group_id IN ?", ids).
		Delete(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("delete user_user_group_relation failed")
		return err
	}

	if _, err := gorm.G[model.UserGroup](ro.data.DB(ctx)).
		Where("id IN ?", ids).
		Delete(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("delete user_group failed")
		return err
	}
	return nil
}

func (ro userGroupRepo) AddUsers(ctx context.Context, groupID int64, userIDs []int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "AddUsers")
	defer span.End()

	userIDs = uniqueInt64(userIDs)
	if groupID == 0 || len(userIDs) == 0 {
		return nil
	}

	// Ensure user group exists.
	gdb := gorm.G[model.UserGroup](ro.data.DB(ctx))
	gCount, err := gdb.Where("id = ?", groupID).Count(ctx, "*")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("check user_group exists failed")
		return err
	}
	if gCount == 0 {
		return biz.ErrRecordNotFound(ctx)
	}

	// Ensure all users exist.
	uCount, err := gorm.G[model.User](ro.data.DB(ctx)).
		Where("id IN ?", userIDs).
		Count(ctx, "*")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("check user exists failed")
		return err
	}
	if uCount != int64(len(userIDs)) {
		return biz.ErrRecordNotFound(ctx)
	}

	rels := make([]model.UserUserGroupRelation, 0, len(userIDs))
	for _, uid := range userIDs {
		rels = append(rels, model.UserUserGroupRelation{
			UserID:      uid,
			UserGroupID: groupID,
		})
	}

	rdb := gorm.G[model.UserUserGroupRelation](ro.data.DB(ctx), clause.OnConflict{DoNothing: true})
	if err := rdb.CreateInBatches(ctx, &rels, 200); err != nil {
		log.WithContext(ctx).WithError(err).Error("add users to user_group failed")
		return err
	}
	return nil
}

func (ro userGroupRepo) RemoveUsers(ctx context.Context, groupID int64, userIDs []int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "RemoveUsers")
	defer span.End()

	userIDs = uniqueInt64(userIDs)
	if groupID == 0 || len(userIDs) == 0 {
		return nil
	}

	_, err = gorm.G[model.UserUserGroupRelation](ro.data.DB(ctx)).
		Where("user_group_id = ? AND user_id IN ?", groupID, userIDs).
		Delete(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("remove users from user_group failed")
	}
	return err
}

func (ro userGroupRepo) replaceUsers(ctx context.Context, groupID int64, userIDs []int64) error {
	if groupID == 0 {
		return nil
	}

	if _, err := gorm.G[model.UserUserGroupRelation](ro.data.DB(ctx)).
		Where("user_group_id = ?", groupID).
		Delete(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("delete user_user_group_relation failed")
		return err
	}
	return ro.AddUsers(ctx, groupID, userIDs)
}

func (ro userGroupRepo) actionCodesExist(ctx context.Context, codes string) error {
	codes = strings.TrimSpace(codes)
	if codes == "" {
		return nil
	}

	arr := splitComma(codes)
	if len(arr) == 0 {
		return nil
	}

	count, err := gorm.G[model.Action](ro.data.DB(ctx)).
		Where("code IN ?", arr).
		Count(ctx, "*")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("check action code exists failed")
		return err
	}
	if count != int64(len(arr)) {
		return biz.ErrRecordNotFound(ctx)
	}
	return nil
}

func splitComma(s string) []string {
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

func uniqueInt64(ids []int64) []int64 {
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

func (ro userGroupRepo) getActionsByCode(ctx context.Context, code string) []biz.Action {
	rp := make([]biz.Action, 0)
	code = strings.TrimSpace(code)
	if code == "" {
		return rp
	}
	codes := ro.splitCodes(code)
	if len(codes) == 0 {
		return rp
	}

	list, err := gorm.G[model.Action](ro.data.DB(ctx)).
		Where("code IN ?", codes).
		Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find actions by code failed")
		return rp
	}
	copierx.Copy(&rp, list)
	return rp
}

func (ro userGroupRepo) getUsersByGroupID(ctx context.Context, groupID int64) []biz.User {
	rp := make([]biz.User, 0)
	if groupID == 0 {
		return rp
	}

	rels, err := gorm.G[model.UserUserGroupRelation](ro.data.DB(ctx)).
		Where("user_group_id = ?", groupID).
		Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find user_user_group_relation failed")
		return rp
	}
	if len(rels) == 0 {
		return rp
	}

	userIDs := make([]int64, 0, len(rels))
	for _, rel := range rels {
		if rel.UserID != 0 {
			userIDs = append(userIDs, rel.UserID)
		}
	}
	if len(userIDs) == 0 {
		return rp
	}

	users, err := gorm.G[model.User](ro.data.DB(ctx)).
		Where("id IN ?", userIDs).
		Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find users by group failed")
		return rp
	}
	copierx.Copy(&rp, users)
	return rp
}

func (userGroupRepo) splitCodes(code string) []string {
	arr := strings.Split(code, ",")
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}
