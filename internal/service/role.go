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

func (s *AuthService) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateRole")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Role{}
	copierx.Copy(&r, req)
	err = s.role.Create(ctx, r)
	return
}

func (s *AuthService) FindRole(ctx context.Context, req *auth.FindRoleRequest) (rp *auth.FindRoleReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindRole")
	defer span.End()
	rp = &auth.FindRoleReply{}
	rp.Page = &params.Page{}
	r := &biz.FindRole{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res := s.role.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateRole(ctx context.Context, req *auth.UpdateRoleRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateRole")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateRole{}
	copierx.Copy(&r, req)
	err = s.role.Update(ctx, r)
	if err == nil {
		s.permission.FlushCache(ctx)
	}
	return
}

func (s *AuthService) DeleteRole(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteRole")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.role.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	if err == nil {
		s.permission.FlushCache(ctx)
	}
	return
}
