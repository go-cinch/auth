package jwt

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
)

const (
	ClaimAuthorityId = "authorityId"
	ClaimExpires     = "exp"
)

type User struct {
	Code string
}

func FromContext(ctx context.Context) (u *User) {
	u = new(User)
	if claims, ok := jwt.FromContext(ctx); ok {
		u.Code = claims.(jwtV4.MapClaims)[ClaimAuthorityId].(string)
	}
	return
}

func (u *User) CreateToken(key, duration string) (token string, expires carbon.Carbon) {
	expires = carbon.Now().AddDuration(duration)
	claims := jwtV4.NewWithClaims(
		jwtV4.SigningMethodHS512,
		jwtV4.MapClaims{
			ClaimAuthorityId: u.Code,
			ClaimExpires:     expires.Timestamp(),
		},
	)
	token, _ = claims.SignedString([]byte(key))
	return
}
