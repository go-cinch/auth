package biz

import (
	"context"

	"auth/api/reason"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/middleware/i18n"
)

var (
	ErrJwtMissingToken = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.JwtMissingToken, reason.ErrorUnauthorized)
	}
	ErrJwtTokenInvalid = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.JwtTokenInvalid, reason.ErrorUnauthorized)
	}
	ErrJwtTokenExpired = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.JwtTokenExpired, reason.ErrorUnauthorized)
	}
	ErrJwtTokenParseFail = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.JwtTokenParseFail, reason.ErrorUnauthorized)
	}
	ErrJwtUnSupportSigningMethod = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.JwtUnSupportSigningMethod, reason.ErrorUnauthorized)
	}
	ErrIdempotentMissingToken = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.IdempotentMissingToken, reason.ErrorIllegalParameter)
	}
	ErrIdempotentTokenExpired = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.IdempotentTokenExpired, reason.ErrorIllegalParameter)
	}

	ErrTooManyRequests = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.TooManyRequests, reason.ErrorTooManyRequests)
	}
	ErrDataNotChange = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.DataNotChange, reason.ErrorIllegalParameter, args...)
	}
	ErrDuplicateField = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.DuplicateField, reason.ErrorIllegalParameter, args...)
	}
	ErrRecordNotFound = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, constant.RecordNotFound, reason.ErrorNotFound, args...)
	}
	ErrNoPermission = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.NoPermission, reason.ErrorForbidden)
	}

	ErrIncorrectPassword = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.IncorrectPassword, reason.ErrorIllegalParameter)
	}
	ErrSamePassword = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.SamePassword, reason.ErrorIllegalParameter)
	}
	ErrInvalidCaptcha = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.InvalidCaptcha, reason.ErrorIllegalParameter)
	}
	ErrLoginFailed = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.LoginFailed, reason.ErrorIllegalParameter)
	}
	ErrUserLocked = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.UserLocked, reason.ErrorForbidden)
	}
	ErrKeepLeastOneAction = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.KeepLeastOneAction, reason.ErrorIllegalParameter)
	}
	ErrDeleteYourself = func(ctx context.Context) error {
		return i18n.NewError(ctx, constant.DeleteYourself, reason.ErrorIllegalParameter)
	}
)
