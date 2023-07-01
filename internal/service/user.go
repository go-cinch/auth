package service

import (
	"context"

	"auth/api/auth"
	"auth/internal/biz"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/proto/params"
	"github.com/go-cinch/common/utils"
	"github.com/golang-module/carbon/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) GetUserByCode(ctx context.Context, req *auth.GetUserByCodeRequest) (rp *auth.GetUserByCodeReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetUserByCode")
	defer span.End()
	rp = &auth.GetUserByCodeReply{}
	user, err := s.user.GetUserByCode(ctx, req.Code)
	if err != nil {
		return
	}
	rp.User = &auth.User{}
	copierx.Copy(&rp.User, user)
	return
}

func (s *AuthService) FindUser(ctx context.Context, req *auth.FindUserRequest) (rp *auth.FindUserReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindUser")
	defer span.End()
	rp = &auth.FindUserReply{}
	rp.Page = &params.Page{}
	r := &biz.FindUser{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res := s.user.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateUser(ctx context.Context, req *auth.UpdateUserRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "UpdateUser")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.UpdateUser{}
	copierx.Copy(&r, req)
	if req.LockExpireTime != nil {
		lockExpire := carbon.Parse(*req.LockExpireTime).Timestamp()
		r.LockExpire = &lockExpire
	}
	err = s.user.Update(ctx, r)
	if err == nil {
		s.permission.FlushCache(ctx)
	}
	return
}

func (s *AuthService) DeleteUser(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteUser")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.user.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	if err == nil {
		s.permission.FlushCache(ctx)
	}
	return
}
