package service

import (
	"context"
	"errors"
	"strings"

	"go.opentelemetry.io/otel"

	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/jwt"
	"github.com/go-cinch/common/log"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "auth/api/auth"
	"auth/internal/biz"
	"auth/internal/conf"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthService)

// AuthService implements the Auth gRPC/HTTP service.
type AuthService struct {
	v1.UnimplementedAuthServer

	c *conf.Bootstrap

	uc         *biz.AuthUseCase
	user       *biz.UserUseCase
	role       *biz.RoleUseCase
	permission *biz.PermissionUseCase
	action     *biz.ActionUseCase
	userGroup  *biz.UserGroupUseCase
	whitelist  *biz.WhitelistUseCase
	hotspot    biz.HotspotRepo
	health     biz.HealthRepo
}

// NewAuthService creates a new service instance.
func NewAuthService(
	c *conf.Bootstrap,
	uc *biz.AuthUseCase,
	user *biz.UserUseCase,
	role *biz.RoleUseCase,
	permission *biz.PermissionUseCase,
	action *biz.ActionUseCase,
	userGroup *biz.UserGroupUseCase,
	whitelist *biz.WhitelistUseCase,
	hotspot biz.HotspotRepo,
	health biz.HealthRepo,
) *AuthService {
	return &AuthService{
		c:          c,
		uc:         uc,
		user:       user,
		role:       role,
		permission: permission,
		action:     action,
		userGroup:  userGroup,
		whitelist:  whitelist,
		hotspot:    hotspot,
		health:     health,
	}
}

func (s *AuthService) flushCache(ctx context.Context) {
	s.user.FlushCache(ctx)
	if err := s.hotspot.Refresh(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Warn("refresh hotspot failed")
	}
}

func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Register")
	defer span.End()

	r := &biz.User{}
	copierx.Copy(&r, req)
	// Skip captcha verification if configured.
	if !s.c.Server.SkipCaptcha {
		if !s.user.VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaAnswer) {
			return nil, biz.ErrInvalidCaptcha(ctx)
		}
	}

	if err := s.user.Create(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) Pwd(ctx context.Context, req *v1.PwdRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Pwd")
	defer span.End()

	r := &biz.User{}
	copierx.Copy(&r, req)
	if err := s.user.Pwd(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Login")
	defer span.End()
	r := &biz.Login{}
	copierx.Copy(&r, req)

	res, err := s.user.Login(ctx, r)
	if err != nil {
		loginFailedErr := biz.ErrLoginFailed(ctx)
		loginFailed := err.Error() == loginFailedErr.Error()
		notFound := err.Error() == biz.ErrRecordNotFound(ctx).Error()
		invalidCaptcha := err.Error() == biz.ErrInvalidCaptcha(ctx).Error()
		if invalidCaptcha {
			return nil, err
		}
		if notFound {
			// Avoid username probing.
			return nil, loginFailedErr
		}
		if loginFailed {
			// Record wrong password attempts; include a timestamp to avoid races with a later successful login.
			loginTime := biz.LoginTime{
				Username: req.Username,
				LastLogin: carbon.DateTime{
					Carbon: carbon.Now(),
				},
				Wrong: res.Wrong,
			}
			s.user.WrongPwd(ctx, &loginTime)
			s.flushCache(ctx)
		}
		return nil, err
	}

	// Successful login: update last_login and reset wrong/locked flags.
	s.user.LastLogin(ctx, req.Username)
	s.flushCache(ctx)

	rp := &v1.LoginReply{}
	copierx.Copy(&rp, res)
	return rp, nil
}

func (s *AuthService) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	_, span := tr.Start(ctx, "Logout")
	defer span.End()

	return &emptypb.Empty{}, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.LoginReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Refresh")
	defer span.End()

	tokenStr := strings.TrimSpace(req.GetToken())
	if tokenStr == "" {
		return nil, biz.ErrJwtMissingToken(ctx)
	}

	info, err := parseToken(ctx, s.c.Server.Jwt.Key, tokenStr)
	if err != nil {
		return nil, err
	}

	claims, ok := info.Claims.(jwtV4.MapClaims)
	if !ok {
		return nil, biz.ErrJwtTokenParseFail(ctx)
	}

	// Extract attrs from claims (defaults to empty strings when missing).
	code, _ := claims["code"].(string)
	platform, _ := claims["platform"].(string)
	authUser := jwt.User{
		Attrs: map[string]string{
			"code":     code,
			"platform": platform,
		},
	}
	token, expireTime := authUser.CreateToken(s.c.Server.Jwt.Key, s.c.Server.Jwt.Expires)
	return &v1.LoginReply{
		Token:   token,
		Expires: expireTime.ToDateTimeString(),
	}, nil
}
func (s *AuthService) Captcha(ctx context.Context, _ *emptypb.Empty) (*v1.CaptchaReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Captcha")
	defer span.End()

	res := s.user.Captcha(ctx)
	rp := &v1.CaptchaReply{
		Captcha: &v1.Captcha{},
	}
	copierx.Copy(&rp.Captcha, &res)
	return rp, nil
}

func (s *AuthService) Status(ctx context.Context, req *v1.StatusRequest) (*v1.StatusReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Status")
	defer span.End()

	res, err := s.user.Status(ctx, req.GetUsername(), true)
	if err != nil {
		return nil, err
	}

	rp := &v1.StatusReply{
		Locked: res.Locked,
	}
	rp.LockExpire = res.LockExpire
	if res.NeedCaptcha {
		rp.Captcha = &v1.Captcha{}
		copierx.Copy(&rp.Captcha, &res.Captcha)
	}
	return rp, nil
}

func parseToken(ctx context.Context, key, jwtToken string) (info *jwtV4.Token, err error) {
	info, err = jwtV4.Parse(jwtToken, func(_ *jwtV4.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		var ve *jwtV4.ValidationError
		if !errors.As(err, &ve) {
			return nil, err
		}
		switch {
		case ve.Errors&jwtV4.ValidationErrorMalformed != 0:
			return nil, biz.ErrJwtTokenInvalid(ctx)
		case ve.Errors&(jwtV4.ValidationErrorExpired|jwtV4.ValidationErrorNotValidYet) != 0:
			return nil, biz.ErrJwtTokenExpired(ctx)
		default:
			return nil, biz.ErrJwtTokenParseFail(ctx)
		}
	}
	if !info.Valid {
		return nil, biz.ErrJwtTokenParseFail(ctx)
	}
	if info.Method != jwtV4.SigningMethodHS512 {
		return nil, biz.ErrJwtUnSupportSigningMethod(ctx)
	}
	return info, nil
}
