package biz

const (
	JwtMissingToken           = "jwt.token.missing"
	JwtTokenInvalid           = "jwt.token.invalid"
	JwtTokenExpired           = "jwt.token.expired"
	JwtTokenParseFail         = "jwt.token.parse.failed"
	JwtUnSupportSigningMethod = "jwt.wrong.signing.method"
	IdempotentMissingToken    = "idempotent.token.missing"
	IdempotentTokenExpired    = "idempotent.token.invalid"

	TooManyRequests    = "too.many.requests"
	DataNotChange      = "data.not.change"
	DuplicateField     = "duplicate.field"
	RecordNotFound     = "record.not.found"
	NoPermission       = "no.permission"
	IncorrectPassword  = "login.incorrect.password"
	SamePassword       = "login.same.password"
	InvalidCaptcha     = "login.invalid.captcha"
	LoginFailed        = "login.failed"
	UserLocked         = "login.user.locked"
	KeepLeastOneAction = "action.keep.least.one.action"
	DeleteYourself     = "user.delete.yourself"
)
