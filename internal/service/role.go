package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/copierx"
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
