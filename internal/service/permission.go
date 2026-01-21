package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/jwt"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "auth/api/auth"
	"auth/internal/biz"
)

func (s *AuthService) Permission(ctx context.Context, req *v1.PermissionRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Permission")
	defer span.End()

	u := jwt.FromServerContext(ctx)
	info := s.user.Info(ctx, u.Attrs["code"])
	if info == nil || info.Id == 0 {
		return nil, biz.ErrNoPermission(ctx)
	}

	check := &biz.CheckPermission{
		UserID: info.Id,
	}
	copierx.Copy(check, req)

	ok, err := s.permission.Check(ctx, check)
	if err != nil || !ok {
		return nil, biz.ErrNoPermission(ctx)
	}

	// Append minimal user attrs for auth_request integrations.
	jwt.AppendToReplyHeader(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     info.Code,
			"platform": info.Platform,
		},
	})

	return &emptypb.Empty{}, nil
}

func (s *AuthService) Info(ctx context.Context, _ *emptypb.Empty) (*v1.InfoReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Info")
	defer span.End()

	u := jwt.FromServerContext(ctx)
	info := s.user.Info(ctx, u.Attrs["code"])
	if info == nil || info.Id == 0 {
		return nil, biz.ErrNoPermission(ctx)
	}

	rp := &v1.InfoReply{
		Permission: &v1.Permission{},
	}

	p, _ := s.permission.GetByUserID(ctx, info.Id)
	copierx.Copy(&rp.Permission, p)
	copierx.Copy(&rp, info)
	return rp, nil
}
