package service

import (
	"context"
	"strings"
	"time"

	"auth/api/auth"
	"auth/internal/biz"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/utils"
	"github.com/go-cinch/common/worker"
	"github.com/golang-module/carbon/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Register")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.User{}
	copierx.Copy(&r, req)
	if !s.user.VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaAnswer) {
		err = biz.ErrInvalidCaptcha(ctx)
		return
	}
	err = s.user.Create(ctx, r)
	s.flushCache(ctx)
	return
}

func (s *AuthService) Pwd(ctx context.Context, req *auth.PwdRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Pwd")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.User{}
	copierx.Copy(&r, req)
	err = s.user.Pwd(ctx, r)
	s.flushCache(ctx)
	return
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (rp *auth.LoginReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Login")
	defer span.End()
	rp = &auth.LoginReply{}
	r := &biz.Login{}
	copierx.Copy(&r, req)
	res, err := s.user.Login(ctx, r)
	if err != nil {
		loginFailedErr := biz.ErrLoginFailed(ctx)
		loginFailed := err.Error() == loginFailedErr.Error()
		notFound := err.Error() == biz.ErrRecordNotFound(ctx).Error()
		invalidCaptcha := err.Error() == biz.ErrInvalidCaptcha(ctx).Error()
		if invalidCaptcha {
			return 
		}
		if notFound {
			// avoid guess username
			err = loginFailedErr
			return
		}
		if loginFailed {
			_ = s.task.Once(
				worker.WithRunCtx(ctx),
				worker.WithRunUUID(strings.Join([]string{s.c.Task.Group.LoginFailed, req.Username}, ".")),
				worker.WithRunGroup(s.c.Task.Group.LoginFailed),
				worker.WithRunNow(true),
				worker.WithRunTimeout(10),
				worker.WithRunReplace(true),
				worker.WithRunPayload(utils.Struct2Json(biz.LoginTime{
					Username: req.Username,
					LastLogin: carbon.DateTime{
						Carbon: carbon.Now(),
					},
					Wrong: res.Wrong,
				})),
			)
			// need refresh hotspot
			s.flushCache(ctx)
			return
		}
		return
	}
	copierx.Copy(&rp, res)
	_ = s.task.Once(
		worker.WithRunCtx(ctx),
		worker.WithRunUUID(strings.Join([]string{s.c.Task.Group.LoginLast, req.Username}, ".")),
		worker.WithRunGroup(s.c.Task.Group.LoginLast),
		worker.WithRunIn(time.Duration(10)*time.Second),
		worker.WithRunTimeout(3),
		worker.WithRunReplace(true),
		worker.WithRunPayload(utils.Struct2Json(biz.LoginTime{
			Username: req.Username,
		})),
	)
	s.flushCache(ctx)
	return
}

func (*AuthService) Logout(ctx context.Context, _ *emptypb.Empty) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Logout")
	defer span.End()
	rp = &emptypb.Empty{}
	return
}

func (s *AuthService) Status(ctx context.Context, req *auth.StatusRequest) (rp *auth.StatusReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Status")
	defer span.End()
	rp = &auth.StatusReply{}
	res, err := s.user.Status(ctx, req.Username, true)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

func (s *AuthService) Captcha(ctx context.Context, _ *emptypb.Empty) (rp *auth.CaptchaReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Captcha")
	defer span.End()
	rp = &auth.CaptchaReply{}
	rp.Captcha = &auth.Captcha{}
	res := s.user.Captcha(ctx)
	copierx.Copy(&rp.Captcha, res)
	return
}

func (s *AuthService) Permission(ctx context.Context, req *auth.PermissionRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Permission")
	defer span.End()
	rp = &emptypb.Empty{}
	user := jwt.FromServerContext(ctx)
	r := &biz.CheckPermission{
		UserCode: user.Attrs["code"],
	}
	if req.Resource != nil {
		r.Resource = *req.Resource
	}
	if req.Method != nil {
		r.Method = *req.Method
	}
	if req.Uri != nil {
		r.URI = *req.Uri
	}
	pass := s.permission.Check(ctx, r)
	if !pass {
		err = biz.ErrNoPermission(ctx)
		return
	}
	info := s.user.Info(ctx, user.Attrs["code"])
	jwt.AppendToReplyHeader(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     info.Code,
			"platform": info.Platform,
		},
	})
	return
}

func (s *AuthService) Info(ctx context.Context, _ *emptypb.Empty) (rp *auth.InfoReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Info")
	defer span.End()
	rp = &auth.InfoReply{}
	rp.Permission = &auth.Permission{}
	user := jwt.FromServerContext(ctx)
	res := s.user.Info(ctx, user.Attrs["code"])
	permission := s.permission.GetByUserCode(ctx, user.Attrs["code"])
	copierx.Copy(&rp.Permission, permission)
	copierx.Copy(&rp, res)
	return
}
