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

func (s *AuthService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateRole")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Role{}
	copierx.Copy(&r, req)
	err = s.role.Create(ctx, r)
	return
}

func (s *AuthService) FindRole(ctx context.Context, req *v1.FindRoleRequest) (rp *v1.FindRoleReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindRole")
	defer span.End()
	rp = &v1.FindRoleReply{}
	rp.Page = &v1.Page{}
	r := &biz.FindRole{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.role.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateRole")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateRole{}
	copierx.Copy(&r, req)
	err = s.role.Update(ctx, r)
	return
}

func (s *AuthService) DeleteRole(ctx context.Context, req *v1.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteRole")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.role.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
