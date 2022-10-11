package data

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/id"
	"github.com/jinzhu/copier"
	"strings"
)

type actionRepo struct {
	data *Data
}

// Action is database fields map
type Action struct {
	Id       uint64 `json:"id,string"` // auto increment id
	Code     string `json:"code"`      // unique code
	Name     string `json:"name"`      // name
	Key      string `json:"key"`       // keyword, must be unique, used as frontend display
	Resource string `json:"resource"`  // resource array, split by break line str, example: GET,/user+\n+POST,/role+\n+GET,/action
}

func NewActionRepo(data *Data) biz.ActionRepo {
	return &actionRepo{
		data: data,
	}
}

func (ro actionRepo) Create(ctx context.Context, item *biz.Action) (err error) {
	var m Action
	db := ro.data.DB(ctx)
	db.
		Where("`key` = ?", item.Key).
		First(&m)
	if m.Id > constant.UI0 {
		err = biz.DuplicateActionKey
		return
	}
	copier.Copy(&m, item)
	m.Id = ro.data.Id(ctx)
	m.Code = id.NewCode(m.Id)
	if m.Resource == "" {
		m.Resource = "*"
	}
	err = db.Create(&m).Error
	return
}

func (ro actionRepo) CodeExists(ctx context.Context, code string) (err error) {
	var m Action
	db := ro.data.DB(ctx)
	arr := strings.Split(code, ",")
	for _, item := range arr {
		db.
			Where("code = ?", item).
			First(&m)
		ok := m.Id > constant.UI1
		if !ok {
			err = v1.ErrorIllegalParameter("%s: %s", biz.ActionNotFound.Message, item)
			return
		}
	}
	return
}
