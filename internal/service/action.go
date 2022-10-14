package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateAction(ctx context.Context, req *v1.CreateActionRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateAction")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Action{}
	copierx.Copy(&r, req)
	err = s.action.Create(ctx, r)
	return
}

func (s *AuthService) FindAction(ctx context.Context, req *v1.FindActionRequest) (rp *v1.FindActionReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindAction")
	defer span.End()
	rp = &v1.FindActionReply{}
	rp.Page = &v1.Page{}
	r := &biz.FindAction{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.action.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}
