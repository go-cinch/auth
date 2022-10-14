package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/page"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
)

func (s *AuthService) FindUser(ctx context.Context, req *v1.FindUserRequest) (rp *v1.FindUserReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindUser")
	defer span.End()
	rp = &v1.FindUserReply{}
	rp.Page = &v1.Page{}
	r := &biz.FindUser{}
	r.Page = page.Page{}
	copier.Copy(&r, req)
	copier.Copy(&r.Page, req.Page)
	res, err := s.user.Find(ctx, r)
	copier.Copy(&rp.Page, r.Page)
	copier.Copy(&rp.List, res)
	return
}
