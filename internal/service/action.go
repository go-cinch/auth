package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateAction(ctx context.Context, req *v1.CreateActionRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateAction")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Action{}
	copier.Copy(&r, req)
	err = s.action.Create(ctx, r)
	return
}
