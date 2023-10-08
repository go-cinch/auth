package biz

import (
	"context"

	"auth/api/reason"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/middleware/i18n"
)

var (
	ErrJwtMissingToken = func(ctx context.Context) error {
		return reason.ErrorUnauthorized(i18n.FromContext(ctx).T(constant.JwtMissingToken))
	}
	ErrJwtTokenInvalid = func(ctx context.Context) error {
		return reason.ErrorUnauthorized(i18n.FromContext(ctx).T(constant.JwtTokenInvalid))
	}
	ErrJwtTokenExpired = func(ctx context.Context) error {
		return reason.ErrorUnauthorized(i18n.FromContext(ctx).T(constant.JwtTokenExpired))
	}
	ErrJwtTokenParseFail = func(ctx context.Context) error {
		return reason.ErrorUnauthorized(i18n.FromContext(ctx).T(constant.JwtTokenParseFail))
	}
	ErrJwtUnSupportSigningMethod = func(ctx context.Context) error {
		return reason.ErrorUnauthorized(i18n.FromContext(ctx).T(constant.JwtUnSupportSigningMethod))
	}
	ErrIdempotentMissingToken = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.IdempotentMissingToken))
	}
	ErrIdempotentTokenExpired = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.IdempotentTokenExpired))
	}

	ErrTooManyRequests = func(ctx context.Context) error {
		return reason.ErrorTooManyRequests(i18n.FromContext(ctx).T(constant.TooManyRequests))
	}
	ErrDataNotChange = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.DataNotChange))
	}
	ErrDuplicateField = func(ctx context.Context, k, v string) error {
		return reason.ErrorIllegalParameter("%s `%s`: %s", i18n.FromContext(ctx).T(constant.DuplicateField), k, v)
	}
	ErrRecordNotFound = func(ctx context.Context) error {
		return reason.ErrorNotFound(i18n.FromContext(ctx).T(constant.RecordNotFound))
	}
	ErrNoPermission = func(ctx context.Context) error {
		return reason.ErrorForbidden(i18n.FromContext(ctx).T(constant.NoPermission))
	}

	ErrIncorrectPassword = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.IncorrectPassword))
	}
	ErrSamePassword = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.SamePassword))
	}
	ErrInvalidCaptcha = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.InvalidCaptcha))
	}
	ErrLoginFailed = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.LoginFailed))
	}
	ErrUserLocked = func(ctx context.Context) error {
		return reason.ErrorForbidden(i18n.FromContext(ctx).T(constant.UserLocked))
	}
	ErrKeepLeastOneAction = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.KeepLeastOneAction))
	}
	ErrDeleteYourself = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter(i18n.FromContext(ctx).T(constant.DeleteYourself))
	}
)
