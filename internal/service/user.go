package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"context"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"github.com/golang-module/carbon/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) FindUser(ctx context.Context, req *v1.FindUserRequest) (rp *v1.FindUserReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "FindUser")
	defer span.End()
	rp = &v1.FindUserReply{}
	rp.Page = &v1.Page{}
	r := &biz.FindUser{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res := s.user.Find(ctx, r)
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

func (s *AuthService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (rp *emptypb.Empty, err error) {
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
	return
}

func (s *AuthService) DeleteUser(ctx context.Context, req *v1.IdsRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "DeleteUser")
	defer span.End()
	rp = &emptypb.Empty{}
	err = s.user.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
