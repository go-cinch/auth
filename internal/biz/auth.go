package biz

import (
	"context"
	"strings"

	"auth/internal/conf"
	"github.com/go-cinch/common/jwt"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest is the core authentication request payload.
// Keep it minimal (no captcha/lock/cache flows) for the standard preset.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Platform string `json:"platform"`
}

type LoginReply struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// LoginUser is the minimum user info required to authenticate and issue a JWT.
type LoginUser struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Code     string `json:"code"`
	Platform string `json:"platform"`
	Locked   bool   `json:"locked"`

	PasswordHash string `json:"-"`
}

type AuthUseCase struct {
	c     *conf.Bootstrap
	repo  AuthRepo
	tx    Transaction
	cache Cache
}

func NewAuthUseCase(c *conf.Bootstrap, repo AuthRepo, tx Transaction, cache Cache) *AuthUseCase {
	if cache != nil {
		cache = cache.WithPrefix(strings.Join([]string{c.Name, "auth"}, "_"))
	}
	return &AuthUseCase{
		c:     c,
		repo:  repo,
		tx:    tx,
		cache: cache,
	}
}

// Login verifies username/password and returns a signed JWT.
func (uc *AuthUseCase) Login(ctx context.Context, req *LoginRequest) (*LoginReply, error) {
	if req == nil {
		return nil, ErrIllegalParameter(ctx, "req")
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Platform = strings.TrimSpace(req.Platform)

	if req.Username == "" || req.Password == "" {
		return nil, ErrIllegalParameter(ctx, "username/password")
	}

	u, err := uc.repo.GetLoginUser(ctx, req.Username)
	if err != nil {
		// Avoid leaking whether the username exists.
		if err.Error() == ErrRecordNotFound(ctx).Error() {
			return nil, ErrLoginFailed(ctx)
		}
		return nil, err
	}
	if u == nil {
		return nil, ErrLoginFailed(ctx)
	}
	if u.Locked {
		return nil, ErrUserLocked(ctx)
	}
	if req.Platform != "" && u.Platform != "" && req.Platform != u.Platform {
		return nil, ErrLoginFailed(ctx)
	}
	if !comparePassword(req.Password, u.PasswordHash) {
		return nil, ErrLoginFailed(ctx)
	}

	authUser := jwt.User{
		Attrs: map[string]string{
			"code":     u.Code,
			"platform": u.Platform,
		},
	}
	token, expireTime := authUser.CreateToken(uc.c.Server.Jwt.Key, uc.c.Server.Jwt.Expires)
	return &LoginReply{
		Token:   token,
		Expires: expireTime.ToDateTimeString(),
	}, nil
}

func comparePassword(plain, hash string) bool {
	if hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
