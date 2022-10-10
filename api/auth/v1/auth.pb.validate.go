// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: auth/v1/auth.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Captcha with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Captcha) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Captcha with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CaptchaMultiError, or nil if none found.
func (m *Captcha) ValidateAll() error {
	return m.validate(true)
}

func (m *Captcha) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Img

	if len(errors) > 0 {
		return CaptchaMultiError(errors)
	}

	return nil
}

// CaptchaMultiError is an error wrapping multiple validation errors returned
// by Captcha.ValidateAll() if the designated constraints aren't met.
type CaptchaMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CaptchaMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CaptchaMultiError) AllErrors() []error { return m }

// CaptchaValidationError is the validation error returned by Captcha.Validate
// if the designated constraints aren't met.
type CaptchaValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CaptchaValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CaptchaValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CaptchaValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CaptchaValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CaptchaValidationError) ErrorName() string { return "CaptchaValidationError" }

// Error satisfies the builtin error interface
func (e CaptchaValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCaptcha.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CaptchaValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CaptchaValidationError{}

// Validate checks the field values on RegisterRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RegisterRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterRequestMultiError, or nil if none found.
func (m *RegisterRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 5 || l > 50 {
		err := RegisterRequestValidationError{
			field:  "Username",
			reason: "value length must be between 5 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 6 || l > 50 {
		err := RegisterRequestValidationError{
			field:  "Password",
			reason: "value length must be between 6 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return RegisterRequestMultiError(errors)
	}

	return nil
}

// RegisterRequestMultiError is an error wrapping multiple validation errors
// returned by RegisterRequest.ValidateAll() if the designated constraints
// aren't met.
type RegisterRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterRequestMultiError) AllErrors() []error { return m }

// RegisterRequestValidationError is the validation error returned by
// RegisterRequest.Validate if the designated constraints aren't met.
type RegisterRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterRequestValidationError) ErrorName() string { return "RegisterRequestValidationError" }

// Error satisfies the builtin error interface
func (e RegisterRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterRequestValidationError{}

// Validate checks the field values on PwdRequest with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PwdRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PwdRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PwdRequestMultiError, or
// nil if none found.
func (m *PwdRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PwdRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Username

	// no validation rules for OldPassword

	if l := utf8.RuneCountInString(m.GetNewPassword()); l < 6 || l > 50 {
		err := PwdRequestValidationError{
			field:  "NewPassword",
			reason: "value length must be between 6 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return PwdRequestMultiError(errors)
	}

	return nil
}

// PwdRequestMultiError is an error wrapping multiple validation errors
// returned by PwdRequest.ValidateAll() if the designated constraints aren't met.
type PwdRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PwdRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PwdRequestMultiError) AllErrors() []error { return m }

// PwdRequestValidationError is the validation error returned by
// PwdRequest.Validate if the designated constraints aren't met.
type PwdRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PwdRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PwdRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PwdRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PwdRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PwdRequestValidationError) ErrorName() string { return "PwdRequestValidationError" }

// Error satisfies the builtin error interface
func (e PwdRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPwdRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PwdRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PwdRequestValidationError{}

// Validate checks the field values on LoginRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginRequestMultiError, or
// nil if none found.
func (m *LoginRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 5 || l > 50 {
		err := LoginRequestValidationError{
			field:  "Username",
			reason: "value length must be between 5 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Password

	if m.CaptchaId != nil {
		// no validation rules for CaptchaId
	}

	if m.CaptchaAnswer != nil {
		// no validation rules for CaptchaAnswer
	}

	if len(errors) > 0 {
		return LoginRequestMultiError(errors)
	}

	return nil
}

// LoginRequestMultiError is an error wrapping multiple validation errors
// returned by LoginRequest.ValidateAll() if the designated constraints aren't met.
type LoginRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginRequestMultiError) AllErrors() []error { return m }

// LoginRequestValidationError is the validation error returned by
// LoginRequest.Validate if the designated constraints aren't met.
type LoginRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginRequestValidationError) ErrorName() string { return "LoginRequestValidationError" }

// Error satisfies the builtin error interface
func (e LoginRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginRequestValidationError{}

// Validate checks the field values on LoginReply with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginReplyMultiError, or
// nil if none found.
func (m *LoginReply) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	// no validation rules for Expires

	if len(errors) > 0 {
		return LoginReplyMultiError(errors)
	}

	return nil
}

// LoginReplyMultiError is an error wrapping multiple validation errors
// returned by LoginReply.ValidateAll() if the designated constraints aren't met.
type LoginReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginReplyMultiError) AllErrors() []error { return m }

// LoginReplyValidationError is the validation error returned by
// LoginReply.Validate if the designated constraints aren't met.
type LoginReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginReplyValidationError) ErrorName() string { return "LoginReplyValidationError" }

// Error satisfies the builtin error interface
func (e LoginReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginReplyValidationError{}

// Validate checks the field values on StatusRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *StatusRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StatusRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in StatusRequestMultiError, or
// nil if none found.
func (m *StatusRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StatusRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 5 || l > 50 {
		err := StatusRequestValidationError{
			field:  "Username",
			reason: "value length must be between 5 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return StatusRequestMultiError(errors)
	}

	return nil
}

// StatusRequestMultiError is an error wrapping multiple validation errors
// returned by StatusRequest.ValidateAll() if the designated constraints
// aren't met.
type StatusRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StatusRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StatusRequestMultiError) AllErrors() []error { return m }

// StatusRequestValidationError is the validation error returned by
// StatusRequest.Validate if the designated constraints aren't met.
type StatusRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StatusRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StatusRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StatusRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StatusRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StatusRequestValidationError) ErrorName() string { return "StatusRequestValidationError" }

// Error satisfies the builtin error interface
func (e StatusRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStatusRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StatusRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StatusRequestValidationError{}

// Validate checks the field values on StatusReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *StatusReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StatusReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in StatusReplyMultiError, or
// nil if none found.
func (m *StatusReply) ValidateAll() error {
	return m.validate(true)
}

func (m *StatusReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCaptcha()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, StatusReplyValidationError{
					field:  "Captcha",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, StatusReplyValidationError{
					field:  "Captcha",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCaptcha()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return StatusReplyValidationError{
				field:  "Captcha",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Locked

	// no validation rules for LockExpire

	if len(errors) > 0 {
		return StatusReplyMultiError(errors)
	}

	return nil
}

// StatusReplyMultiError is an error wrapping multiple validation errors
// returned by StatusReply.ValidateAll() if the designated constraints aren't met.
type StatusReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StatusReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StatusReplyMultiError) AllErrors() []error { return m }

// StatusReplyValidationError is the validation error returned by
// StatusReply.Validate if the designated constraints aren't met.
type StatusReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StatusReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StatusReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StatusReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StatusReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StatusReplyValidationError) ErrorName() string { return "StatusReplyValidationError" }

// Error satisfies the builtin error interface
func (e StatusReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStatusReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StatusReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StatusReplyValidationError{}

// Validate checks the field values on CaptchaReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CaptchaReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CaptchaReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CaptchaReplyMultiError, or
// nil if none found.
func (m *CaptchaReply) ValidateAll() error {
	return m.validate(true)
}

func (m *CaptchaReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCaptcha()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CaptchaReplyValidationError{
					field:  "Captcha",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CaptchaReplyValidationError{
					field:  "Captcha",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCaptcha()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CaptchaReplyValidationError{
				field:  "Captcha",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CaptchaReplyMultiError(errors)
	}

	return nil
}

// CaptchaReplyMultiError is an error wrapping multiple validation errors
// returned by CaptchaReply.ValidateAll() if the designated constraints aren't met.
type CaptchaReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CaptchaReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CaptchaReplyMultiError) AllErrors() []error { return m }

// CaptchaReplyValidationError is the validation error returned by
// CaptchaReply.Validate if the designated constraints aren't met.
type CaptchaReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CaptchaReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CaptchaReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CaptchaReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CaptchaReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CaptchaReplyValidationError) ErrorName() string { return "CaptchaReplyValidationError" }

// Error satisfies the builtin error interface
func (e CaptchaReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCaptchaReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CaptchaReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CaptchaReplyValidationError{}

// Validate checks the field values on RefreshRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *RefreshRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RefreshRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in RefreshRequestMultiError,
// or nil if none found.
func (m *RefreshRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RefreshRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if len(errors) > 0 {
		return RefreshRequestMultiError(errors)
	}

	return nil
}

// RefreshRequestMultiError is an error wrapping multiple validation errors
// returned by RefreshRequest.ValidateAll() if the designated constraints
// aren't met.
type RefreshRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RefreshRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RefreshRequestMultiError) AllErrors() []error { return m }

// RefreshRequestValidationError is the validation error returned by
// RefreshRequest.Validate if the designated constraints aren't met.
type RefreshRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RefreshRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RefreshRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RefreshRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RefreshRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RefreshRequestValidationError) ErrorName() string { return "RefreshRequestValidationError" }

// Error satisfies the builtin error interface
func (e RefreshRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRefreshRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RefreshRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RefreshRequestValidationError{}

// Validate checks the field values on CreateActionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateActionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateActionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateActionRequestMultiError, or nil if none found.
func (m *CreateActionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateActionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) > 50 {
		err := CreateActionRequestValidationError{
			field:  "Name",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKey()) > 50 {
		err := CreateActionRequestValidationError{
			field:  "Key",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Path

	if len(errors) > 0 {
		return CreateActionRequestMultiError(errors)
	}

	return nil
}

// CreateActionRequestMultiError is an error wrapping multiple validation
// errors returned by CreateActionRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateActionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateActionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateActionRequestMultiError) AllErrors() []error { return m }

// CreateActionRequestValidationError is the validation error returned by
// CreateActionRequest.Validate if the designated constraints aren't met.
type CreateActionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateActionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateActionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateActionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateActionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateActionRequestValidationError) ErrorName() string {
	return "CreateActionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateActionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateActionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateActionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateActionRequestValidationError{}

// Validate checks the field values on CreateRoleRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CreateRoleRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateRoleRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateRoleRequestMultiError, or nil if none found.
func (m *CreateRoleRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateRoleRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) > 50 {
		err := CreateRoleRequestValidationError{
			field:  "Name",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKey()) > 50 {
		err := CreateRoleRequestValidationError{
			field:  "Key",
			reason: "value length must be at most 50 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateRoleRequestMultiError(errors)
	}

	return nil
}

// CreateRoleRequestMultiError is an error wrapping multiple validation errors
// returned by CreateRoleRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateRoleRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateRoleRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateRoleRequestMultiError) AllErrors() []error { return m }

// CreateRoleRequestValidationError is the validation error returned by
// CreateRoleRequest.Validate if the designated constraints aren't met.
type CreateRoleRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRoleRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRoleRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRoleRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRoleRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRoleRequestValidationError) ErrorName() string {
	return "CreateRoleRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateRoleRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRoleRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRoleRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRoleRequestValidationError{}