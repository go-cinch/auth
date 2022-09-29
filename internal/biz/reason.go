package biz

import (
	v1 "auth/api/auth/v1"
)

var (
	// TooManyRequests is too many requests in a short time
	TooManyRequests   = v1.ErrorTooManyRequests(v1.ErrorReason_TOO_MANY_REQUESTS.String())
	DuplicateUsername = v1.ErrorDuplicateUsername(v1.ErrorReason_DUPLICATE_USERNAME.String())
	UserNotFound      = v1.ErrorUserNotFound(v1.ErrorReason_USER_NOT_FOUND.String())
	IncorrectPassword = v1.ErrorIncorrectPassword(v1.ErrorReason_INCORRECT_PASSWORD.String())
	SamePassword      = v1.ErrorSamePassword(v1.ErrorReason_SAME_PASSWORD.String())
	InvalidCaptcha    = v1.ErrorInvalidCaptcha(v1.ErrorReason_INVALID_CAPTCHA.String())
	UserLocked        = v1.ErrorUserLocked(v1.ErrorReason_USER_LOCKED.String())
	LoginFailed       = v1.ErrorLoginFailed(v1.ErrorReason_LOGIN_FAILED.String())
)
