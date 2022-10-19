package middleware

import "auth/api/reason"

var (
	MissingToken           = reason.ErrorUnauthorized("token is missing")
	TokenInvalid           = reason.ErrorUnauthorized("token is invalid")
	TokenExpired           = reason.ErrorUnauthorized("token has expired")
	TokenParseFail         = reason.ErrorUnauthorized("fail to parse token")
	WrongContext           = reason.ErrorUnauthorized("wrong context for middleware")
	UnSupportSigningMethod = reason.ErrorUnauthorized("wrong signing method")
	MissingIdempotentToken = reason.ErrorIllegalParameter("idempotent token is missing")
	IdempotentTokenExpired = reason.ErrorIllegalParameter("idempotent token has expired")
)
