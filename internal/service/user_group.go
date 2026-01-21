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

func (s *AuthService) CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "CreateUserGroup")
	defer span.End()

	g := &biz.UserGroup{}
	copierx.Copy(g, req)
	if len(req.GetUsers()) > 0 {
		g.Users = make([]biz.User, 0, len(req.GetUsers()))
		for _, id := range req.GetUsers() {
			g.Users = append(g.Users, biz.User{Id: id})
		}
	}

	if err := s.userGroup.Create(ctx, g); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) FindUserGroup(ctx context.Context, req *v1.FindUserGroupRequest) (*v1.FindUserGroupReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "FindUserGroup")
	defer span.End()

	r := &biz.FindUserGroup{}
	copierx.Copy(&r.Page, req.GetPage())
	copierx.Copy(r, req)

	list, err := s.userGroup.Find(ctx, r)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindUserGroupReply{
		Page: &params.Page{},
	}
	copierx.Copy(&rp.Page, &r.Page)
	copierx.Copy(&rp.List, list)
	return rp, nil
}

func (s *AuthService) UpdateUserGroup(ctx context.Context, req *v1.UpdateUserGroupRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "UpdateUserGroup")
	defer span.End()

	g := &biz.UpdateUserGroup{}
	copierx.Copy(g, req)
	if err := s.userGroup.Update(ctx, g); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) DeleteUserGroup(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "DeleteUserGroup")
	defer span.End()

	if err := s.userGroup.Delete(ctx, utils.Str2Int64Arr(req.GetIds())...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
