package service

import (
	"context"

	"auth/api/auth"
	"auth/internal/biz"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/proto/params"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateUserGroup(ctx context.Context, req *auth.CreateUserGroupRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UserGroup{}
	copierx.Copy(&r, req)
	err = s.userGroup.Create(ctx, r)
	return
}

func (s *AuthService) FindUserGroup(ctx context.Context, req *auth.FindUserGroupRequest) (rp *auth.FindUserGroupReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindUserGroup")
	defer span.End()
	rp = &auth.FindUserGroupReply{}
	rp.Page = &params.Page{}
	r := &biz.FindUserGroup{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.userGroup.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateUserGroup(ctx context.Context, req *auth.UpdateUserGroupRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateUserGroup{}
	copierx.Copy(&r, req)
	err = s.userGroup.Update(ctx, r)
	if err == nil {
		s.permission.FlushCache(ctx)
	}
	return
}

func (s *AuthService) DeleteUserGroup(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.userGroup.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	if err == nil {
		s.permission.FlushCache(ctx)
	}
	return
}
