package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UserGroup{}
	copierx.Copy(&r, req)
	err = s.userGroup.Create(ctx, r)
	return
}

func (s *AuthService) FindUserGroup(ctx context.Context, req *v1.FindUserGroupRequest) (rp *v1.FindUserGroupReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindUserGroup")
	defer span.End()
	rp = &v1.FindUserGroupReply{}
	rp.Page = &v1.Page{}
	r := &biz.FindUserGroup{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.userGroup.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateUserGroup(ctx context.Context, req *v1.UpdateUserGroupRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateUserGroup{}
	copierx.Copy(&r, req)
	err = s.userGroup.Update(ctx, r)
	return
}

func (s *AuthService) DeleteUserGroup(ctx context.Context, req *v1.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.userGroup.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
