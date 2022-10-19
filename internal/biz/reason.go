package biz

import "auth/api/reason"

var (
	IllegalParameter   = reason.ErrorIllegalParameter
	NotFound           = reason.ErrorNotFound
	TooManyRequests    = reason.ErrorTooManyRequests("too many requests, please try again later")
	DataNotChange      = reason.ErrorIllegalParameter("data has not changed")
	DuplicateField     = reason.ErrorIllegalParameter("duplicate field")
	RecordNotFound     = reason.ErrorNotFound("not found")
	IncorrectPassword  = reason.ErrorIllegalParameter("incorrect password")
	SamePassword       = reason.ErrorIllegalParameter("same password")
	InvalidCaptcha     = reason.ErrorIllegalParameter("invalid captcha")
	LoginFailed        = reason.ErrorIllegalParameter("incorrect username or password")
	KeepLeastOntAction = reason.ErrorIllegalParameter("keep at least one action")
	DeleteYourself     = reason.ErrorIllegalParameter("you cannot delete yourself")
	UserLocked         = reason.ErrorForbidden("user is locked")
	NoPermission       = reason.ErrorForbidden("no permission to access this resource")
)
