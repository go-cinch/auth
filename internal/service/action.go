package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/copierx"
	params "github.com/go-cinch/common/proto/params"
	"github.com/go-cinch/common/utils"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "auth/api/auth"
	"auth/internal/biz"
)

func (s *AuthService) CreateAction(ctx context.Context, req *v1.CreateActionRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "CreateAction")
	defer span.End()

	r := &biz.Action{}
	copierx.Copy(r, req)
	if err := s.action.Create(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) FindAction(ctx context.Context, req *v1.FindActionRequest) (*v1.FindActionReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "FindAction")
	defer span.End()

	r := &biz.FindAction{}
	copierx.Copy(&r.Page, req.GetPage())
	copierx.Copy(r, req)

	list, err := s.action.Find(ctx, r)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindActionReply{
		Page: &params.Page{},
	}
	copierx.Copy(&rp.Page, &r.Page)
	copierx.Copy(&rp.List, list)
	return rp, nil
}

func (s *AuthService) UpdateAction(ctx context.Context, req *v1.UpdateActionRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "UpdateAction")
	defer span.End()

	r := &biz.Action{}
	copierx.Copy(r, req)
	if err := s.action.Update(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) DeleteAction(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "DeleteAction")
	defer span.End()

	if err := s.action.Delete(ctx, utils.Str2Int64Arr(req.GetIds())...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
