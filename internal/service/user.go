package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *AuthService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateUser")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateUser{}
	copier.Copy(&r, req)
	err = s.user.Update(ctx, r)
	return
}

func (s *AuthService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteUser")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.user.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
