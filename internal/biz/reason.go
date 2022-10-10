package biz

import (
	v1 "auth/api/auth/v1"
)

var (
	// TooManyRequests is too many requests in a short time
	TooManyRequests    = v1.ErrorTooManyRequests(v1.ErrorReason_TOO_MANY_REQUESTS.String())
	DuplicateUsername  = v1.ErrorIllegalParameter("duplicate username")
	UserNotFound       = v1.ErrorNotFound("user not found")
	IncorrectPassword  = v1.ErrorIllegalParameter("incorrect password")
	SamePassword       = v1.ErrorIllegalParameter("same password")
	InvalidCaptcha     = v1.ErrorIllegalParameter("invalid captcha")
	UserLocked         = v1.ErrorForbidden("user is locked")
	LoginFailed        = v1.ErrorIllegalParameter("incorrect username or password")
	DuplicateActionKey = v1.ErrorIllegalParameter("duplicate key")
	DuplicateRoleKey   = v1.ErrorIllegalParameter("duplicate key")
)
