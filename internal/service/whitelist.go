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

func (s *AuthService) CreateWhitelist(ctx context.Context, req *v1.CreateWhitelistRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "CreateWhitelist")
	defer span.End()

	w := &biz.Whitelist{}
	copierx.Copy(w, req)
	if err := s.whitelist.Create(ctx, w); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) FindWhitelist(ctx context.Context, req *v1.FindWhitelistRequest) (*v1.FindWhitelistReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "FindWhitelist")
	defer span.End()

	r := &biz.FindWhitelist{}
	copierx.Copy(&r.Page, req.GetPage())
	copierx.Copy(r, req)

	list, err := s.whitelist.Find(ctx, r)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindWhitelistReply{
		Page: &params.Page{},
	}
	copierx.Copy(&rp.Page, &r.Page)
	copierx.Copy(&rp.List, list)
	return rp, nil
}

func (s *AuthService) UpdateWhitelist(ctx context.Context, req *v1.UpdateWhitelistRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "UpdateWhitelist")
	defer span.End()

	w := &biz.UpdateWhitelist{}
	copierx.Copy(w, req)

	if err := s.whitelist.Update(ctx, w); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) DeleteWhitelist(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "DeleteWhitelist")
	defer span.End()

	if err := s.whitelist.Delete(ctx, utils.Str2Int64Arr(req.GetIds())...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
