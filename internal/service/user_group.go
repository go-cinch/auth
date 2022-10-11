package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateUserGroup")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UserGroup{}
	copier.Copy(&r, req)
	err = s.userGroup.Create(ctx, r)
	return
}
