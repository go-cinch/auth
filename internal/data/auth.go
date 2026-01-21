package data

import (
	"context"
	"strings"

	"github.com/go-cinch/common/log"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"auth/internal/biz"
	"auth/internal/data/model"
)

type authRepo struct {
	data *Data
}

func NewAuthRepo(data *Data) biz.AuthRepo {
	return &authRepo{data: data}
}

func (r authRepo) GetLoginUser(ctx context.Context, username string) (*biz.LoginUser, error) {
	tr := otel.Tracer("data")
	ctx, span := tr.Start(ctx, "GetLoginUser")
	defer span.End()

	username = strings.TrimSpace(username)
	if username == "" {
		return nil, biz.ErrIllegalParameter(ctx, "username")
	}

	db := gorm.G[model.User](r.data.DB(ctx))
	m, err := db.Where("username = ?", username).First(ctx)
	if err == gorm.ErrRecordNotFound || m.ID == 0 {
		return nil, biz.ErrRecordNotFound(ctx)
	}
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("get login user failed")
		return nil, err
	}

	u := &biz.LoginUser{
		ID:     m.ID,
		Code:   m.Code,
		Locked: m.Locked != nil && *m.Locked != 0,
	}
	if m.Username != nil {
		u.Username = *m.Username
	}
	if m.Platform != nil {
		u.Platform = *m.Platform
	}
	if m.Password != nil {
		u.PasswordHash = *m.Password
	}
	return u, nil
}
