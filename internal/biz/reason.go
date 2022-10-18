package biz

import (
	v1 "auth/api/auth/v1"
)

var (
	IllegalParameter   = v1.ErrorIllegalParameter
	TooManyRequests    = v1.ErrorTooManyRequests("too many requests, please try again later")
	DataNotChange      = v1.ErrorIllegalParameter("data has not changed")
	DuplicateField     = v1.ErrorIllegalParameter("duplicate field")
	NotFound           = v1.ErrorIllegalParameter("not found")
	IncorrectPassword  = v1.ErrorIllegalParameter("incorrect password")
	SamePassword       = v1.ErrorIllegalParameter("same password")
	InvalidCaptcha     = v1.ErrorIllegalParameter("invalid captcha")
	LoginFailed        = v1.ErrorIllegalParameter("incorrect username or password")
	KeepLeastOntAction = v1.ErrorIllegalParameter("keep at least one action")
	DeleteYourself     = v1.ErrorIllegalParameter("you cannot delete yourself")
	UserLocked         = v1.ErrorForbidden("user is locked")
	NoPermission       = v1.ErrorForbidden("no permission to access this resource")
)
