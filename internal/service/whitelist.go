package service

import (
	"context"

	"auth/api/auth"
	"auth/internal/biz"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/proto/params"
	"github.com/go-cinch/common/utils"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) CreateWhitelist(ctx context.Context, req *auth.CreateWhitelistRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "CreateWhitelist")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.Whitelist{}
	copierx.Copy(&r, req)
	err = s.whitelist.Create(ctx, r)
	s.flushUserAndPermissionCache(ctx)
	return
}

func (s *AuthService) HasWhitelist(ctx context.Context, req *auth.HasWhitelistRequest) (rp *auth.HasWhitelistReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "HasWhitelist")
	defer span.End()
	rp = &auth.HasWhitelistReply{}
	r := &biz.HasWhitelist{}
	copierx.Copy(&r, req)
	res, err := s.whitelist.Has(ctx, r)
	if err != nil {
		return
	}
	rp.Ok = res
	return
}

func (s *AuthService) FindWhitelist(ctx context.Context, req *auth.FindWhitelistRequest) (rp *auth.FindWhitelistReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindWhitelist")
	defer span.End()
	rp = &auth.FindWhitelistReply{}
	rp.Page = &params.Page{}
	r := &biz.FindWhitelist{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res := s.whitelist.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateWhitelist(ctx context.Context, req *auth.UpdateWhitelistRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateWhitelist")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateWhitelist{}
	copierx.Copy(&r, req)
	err = s.whitelist.Update(ctx, r)
	s.flushUserAndPermissionCache(ctx)
	return
}

func (s *AuthService) DeleteWhitelist(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteWhitelist")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.whitelist.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	s.flushUserAndPermissionCache(ctx)
	return
}
