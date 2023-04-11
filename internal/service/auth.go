package service

import (
	"auth/api/auth"
	"auth/api/reason"
	"auth/internal/biz"
	"context"
	"errors"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/utils"
	"github.com/go-cinch/common/worker"
	"github.com/golang-module/carbon/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"time"
)

func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (rp *emptypb.Empty, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Register")
	defer span.End()
	rp = &emptypb.Empty{}
	r := &biz.User{}
	copierx.Copy(&r, req)
	if !s.user.VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaAnswer) {
		err = reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.InvalidCaptcha))
		return
	}
	err = s.user.Create(ctx, r)
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
	loginFailedErr := reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(biz.LoginFailed))
	loginFailed := errors.Is(err, loginFailedErr)
	notFound := errors.Is(err, reason.ErrorNotFound(i18n.FromContext(ctx).T(biz.RecordNotFound)))
	ctx, _ = context.WithTimeout(ctx, time.Second)
	if err != nil {
		if loginFailed {
			s.task.Once(
				worker.WithRunCtx(ctx),
				worker.WithRunUuid(strings.Join([]string{"login.failed", req.Username}, ".")),
				worker.WithRunGroup("login.failed"),
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
		} else if notFound {
			// avoid guess username
			err = loginFailedErr
		}
		return
	}
	copierx.Copy(&rp, res)
	s.task.Once(
		worker.WithRunCtx(ctx),
		worker.WithRunUuid(strings.Join([]string{"login.last", req.Username}, ".")),
		worker.WithRunGroup("login.last"),
		worker.WithRunIn(time.Duration(10)*time.Second),
		worker.WithRunTimeout(10),
		worker.WithRunReplace(true),
		worker.WithRunPayload(utils.Struct2Json(biz.LoginTime{
			Username: req.Username,
		})),
	)
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

func (s *AuthService) Captcha(ctx context.Context, req *emptypb.Empty) (rp *auth.CaptchaReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Captcha")
	defer span.End()
	rp = &auth.CaptchaReply{}
	rp.Captcha = &auth.Captcha{}
	res := s.user.Captcha(ctx)
	copierx.Copy(&rp.Captcha, res)
	return
}

func (s *AuthService) Permission(ctx context.Context, req *auth.PermissionRequest) (rp *auth.PermissionReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Permission")
	defer span.End()
	rp = &auth.PermissionReply{}
	user := jwt.FromServerContext(ctx)
	r := &biz.CheckPermission{
		UserCode: user.Code,
		Resource: req.Resource,
	}
	rp.Pass = s.permission.Check(ctx, r)
	info, err := s.user.Info(ctx, user.Code)
	if err != nil {
		return
	}
	jwt.AppendToReplyHeader(ctx, jwt.User{
		Code:     info.Code,
		Platform: info.Platform,
	})
	return
}

func (s *AuthService) Info(ctx context.Context, req *emptypb.Empty) (rp *auth.InfoReply, err error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Info")
	defer span.End()
	rp = &auth.InfoReply{}
	rp.Permission = &auth.Permission{}
	user := jwt.FromServerContext(ctx)
	res, err := s.user.Info(ctx, user.Code)
	if err != nil {
		return
	}
	permission, err := s.permission.GetByUserCode(ctx, user.Code)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Permission, permission)
	copierx.Copy(&rp, res)
	return
}
