package biz

import (
	v1 "auth/api/auth/v1"
)

var (
	// TooManyRequests is too many requests in a short time
	TooManyRequests        = v1.ErrorTooManyRequests(v1.ErrorReason_TOO_MANY_REQUESTS.String())
	DataNotChange          = v1.ErrorIllegalParameter("data has not changed")
	DuplicateUsername      = v1.ErrorIllegalParameter("duplicate username")
	UserNotFound           = v1.ErrorNotFound("user not found")
	IncorrectPassword      = v1.ErrorIllegalParameter("incorrect password")
	SamePassword           = v1.ErrorIllegalParameter("same password")
	InvalidCaptcha         = v1.ErrorIllegalParameter("invalid captcha")
	UserLocked             = v1.ErrorForbidden("user is locked")
	LoginFailed            = v1.ErrorIllegalParameter("incorrect username or password")
	DuplicateActionWord    = v1.ErrorIllegalParameter("duplicate word")
	ActionNotFound         = v1.ErrorNotFound("action not found")
	KeepLeastOntAction     = v1.ErrorIllegalParameter("keep at least one action")
	DuplicateRoleWord      = v1.ErrorIllegalParameter("duplicate word")
	RoleNotFound           = v1.ErrorNotFound("role not found")
	DuplicateUserGroupWord = v1.ErrorIllegalParameter("duplicate word")
	NoPermission           = v1.ErrorForbidden("no permission to access this resource")
	DeleteYourself         = v1.ErrorIllegalParameter("you cannot delete yourself")
	UserGroupNotFound      = v1.ErrorNotFound("user group not found")
)
