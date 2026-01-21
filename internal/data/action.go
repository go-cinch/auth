package data

import (
	"context"
	"errors"
	"strings"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type actionRepo struct {
	data *Data
}

func NewActionRepo(data *Data) biz.ActionRepo {
	return &actionRepo{
		data: data,
	}
}

func (ro actionRepo) Create(ctx context.Context, item *biz.Action) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	db := gorm.G[model.Action](ro.data.DB(ctx))

	// Normalize and check if word exists (word is nullable).
	if item.Word != nil {
		word := strings.TrimSpace(*item.Word)
		item.Word = &word
		count, err := db.Where("word = ?", word).Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("check action word exists failed")
			return err
		}
		if count > 0 {
			return biz.ErrDuplicateField(ctx, "word", word)
		}
	}

	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}

	var m model.Action
	copierx.Copy(&m, item)
	m.ID = item.ID
	if m.Word != nil {
		word := strings.TrimSpace(*m.Word)
		m.Word = &word
	}

	// Always generate code from ID for uniqueness/consistency.
	code := id.NewCode(uint64(m.ID))
	m.Code = code

	// Default resource is "*" when empty.
	if m.Resource == nil || strings.TrimSpace(*m.Resource) == "" {
		res := "*"
		m.Resource = &res
	}

	err = db.Create(ctx, &m)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("create action failed")
		return err
	}

	// Best-effort: propagate generated code back to caller.
	item.Code = &code
	return nil
}

func (ro actionRepo) Update(ctx context.Context, item *biz.Action) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	if item.Word != nil {
		word := strings.TrimSpace(*item.Word)
		item.Word = &word
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))

	m, err := db.Where("id = ?", item.ID).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return biz.ErrRecordNotFound(ctx)
	}
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("get action failed")
		return err
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	delete(change, "id")
	delete(change, "code") // code is derived from id; don't allow updating it here

	if len(change) == 0 {
		return biz.ErrDataNotChange(ctx)
	}

	// Check word uniqueness if word is being updated.
	if item.Word != nil {
		word := strings.TrimSpace(*item.Word)
		oldWord := ""
		if m.Word != nil {
			oldWord = strings.TrimSpace(*m.Word)
		}
		if word != oldWord {
			count, err := db.Where("word = ? AND id <> ?", word, item.ID).Count(ctx, "*")
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("check action word uniqueness failed")
				return err
			}
			if count > 0 {
				return biz.ErrDuplicateField(ctx, "word", word)
			}
		}
	}

	err = ro.data.DB(ctx).
		Model(&model.Action{}).
		Where("id = ?", item.ID).
		Updates(change).
		Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("update action failed")
	}
	return err
}

func (ro actionRepo) Delete(ctx context.Context, ids []int64) (err error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Delete")
	defer span.End()

	if len(ids) == 0 {
		return nil
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	_, err = db.Where("id IN ?", ids).Delete(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("delete action failed")
	}
	return err
}

func (ro actionRepo) Find(ctx context.Context, condition *biz.FindAction) (rp []biz.Action) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "Find")
	defer span.End()

	rp = make([]biz.Action, 0)
	db := gorm.G[model.Action](ro.data.DB(ctx))
	q := db.Where("1 = 1")

	if condition.Code != nil && strings.TrimSpace(*condition.Code) != "" {
		q = q.Where("code LIKE ?", "%"+strings.TrimSpace(*condition.Code)+"%")
	}
	if condition.Name != nil && strings.TrimSpace(*condition.Name) != "" {
		q = q.Where("name LIKE ?", "%"+strings.TrimSpace(*condition.Name)+"%")
	}
	if condition.Word != nil && strings.TrimSpace(*condition.Word) != "" {
		q = q.Where("word LIKE ?", "%"+strings.TrimSpace(*condition.Word)+"%")
	}
	if condition.Resource != nil && strings.TrimSpace(*condition.Resource) != "" {
		q = q.Where("resource LIKE ?", "%"+strings.TrimSpace(*condition.Resource)+"%")
	}

	if !condition.Page.Disable {
		count, err := q.Count(ctx, "*")
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("count action failed")
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
		log.WithContext(ctx).WithError(err).Error("find action failed")
		return rp
	}

	copierx.Copy(&rp, list)
	return rp
}
