package service

import (
	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/pkg/jwt"
	"context"
	"github.com/go-cinch/common/utils"
	"github.com/go-cinch/common/worker"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Register")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.User{}
	copier.Copy(&r, req)
	err = s.user.Create(ctx, r)
	return
}

func (s *AuthService) Pwd(ctx context.Context, req *v1.PwdRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Pwd")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.User{}
	copier.Copy(&r, req)
	err = s.user.Pwd(ctx, r)
	return
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (rp *v1.LoginReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Login")
	defer span.End()
	rp = &v1.LoginReply{}
	r := &biz.Login{}
	copier.Copy(&r, req)
	res, err := s.user.Login(ctx, r)
	if err != nil {
		if err == biz.LoginFailed {
			s.task.Once(
				worker.WithRunUuid(uuid.NewString()),
				worker.WithRunCategory("login.failed"),
				worker.WithRunNow(true),
				worker.WithRunTimeout(10),
				worker.WithRunPayload(utils.Struct2Json(biz.LoginTime{
					Username: req.Username,
					LastLogin: carbon.DateTime{
						Carbon: carbon.Now(),
					},
				})),
			)
		} else if err == biz.UserNotFound {
			// avoid guess username
			err = biz.LoginFailed
		}
		return
	}
	copier.Copy(&rp, res)
	return
}

func (s *AuthService) Status(ctx context.Context, req *v1.StatusRequest) (rp *v1.StatusReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Status")
	defer span.End()
	rp = &v1.StatusReply{}
	res, err := s.user.Status(ctx, req.Username, true)
	if err != nil {
		return
	}
	copier.Copy(&rp, res)
	return
}

func (s *AuthService) Captcha(ctx context.Context, req *emptypb.Empty) (rp *v1.CaptchaReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Captcha")
	defer span.End()
	rp = &v1.CaptchaReply{}
	rp.Captcha = &v1.Captcha{}
	res := s.user.Captcha(ctx)
	copier.Copy(&rp.Captcha, res)
	return
}

func (s *AuthService) Permission(ctx context.Context, req *v1.PermissionRequest) (rp *v1.PermissionReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Permission")
	defer span.End()
	rp = &v1.PermissionReply{}
	r := &biz.Permission{}
	copier.Copy(&r, req)
	rp.Pass = s.permission.Check(ctx, r)
	return
}

func (s *AuthService) Info(ctx context.Context, req *emptypb.Empty) (rp *v1.InfoReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Info")
	defer span.End()
	rp = &v1.InfoReply{}
	user := jwt.FromContext(ctx)
	res, err := s.user.Info(ctx, user.Code)
	if err != nil {
		return
	}
	copier.Copy(&rp, res)
	return
}
