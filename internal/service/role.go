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

func (s *AuthService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "CreateRole")
	defer span.End()

	r := &biz.Role{}
	copierx.Copy(r, req)
	if err := s.role.Create(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) FindRole(ctx context.Context, req *v1.FindRoleRequest) (*v1.FindRoleReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "FindRole")
	defer span.End()

	r := &biz.FindRole{}
	copierx.Copy(&r.Page, req.GetPage())
	copierx.Copy(r, req)

	list, err := s.role.Find(ctx, r)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindRoleReply{
		Page: &params.Page{},
	}
	copierx.Copy(&rp.Page, &r.Page)
	copierx.Copy(&rp.List, list)
	return rp, nil
}

func (s *AuthService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "UpdateRole")
	defer span.End()

	r := &biz.UpdateRole{}
	copierx.Copy(r, req)
	if err := s.role.Update(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) DeleteRole(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "DeleteRole")
	defer span.End()

	if err := s.role.Delete(ctx, utils.Str2Int64Arr(req.GetIds())...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
