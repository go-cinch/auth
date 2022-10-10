package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateRole")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Role{}
	copier.Copy(&r, req)
	err = s.role.Create(ctx, r)
	return
}
