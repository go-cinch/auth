package middleware

import (
	v1 "auth/api/auth/v1"
)

var (
	MissingToken           = v1.ErrorUnauthorized("token is missing")
	TokenInvalid           = v1.ErrorUnauthorized("token is invalid")
	TokenExpired           = v1.ErrorUnauthorized("token has expired")
	TokenParseFail         = v1.ErrorUnauthorized("fail to parse token")
	WrongContext           = v1.ErrorUnauthorized("wrong context for middleware")
	UnSupportSigningMethod = v1.ErrorUnauthorized("wrong signing method")
	MissingIdempotentToken = v1.ErrorIllegalParameter("idempotent token is missing")
	IdempotentTokenExpired = v1.ErrorIllegalParameter("idempotent token has expired")
)
