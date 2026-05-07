package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
	"strings"
)

type ErrCode string

const (
	ErrInternal     ErrCode = "internal_server_error"
	ErrInvalid      ErrCode = "invalid"
	ErrNotFound     ErrCode = "not_found"
	ErrUnauthorized ErrCode = "unauthorized"
	ErrForbidden    ErrCode = "forbidden"
)

var defaultMessages = map[ErrCode]string{
	ErrInternal:     "予期しないエラーが発生しました",
	ErrInvalid:      "入力内容に誤りがあります",
	ErrNotFound:     "リソースが見つかりません",
	ErrUnauthorized: "認証が必要です",
	ErrForbidden:    "アクセス権限がありません",
}

type Error struct {
	code    ErrCode
	message string
	err     error
	stack   []uintptr
}

func New(errs ...error) *Error             { return newErr(ErrInternal, errs...) }
func NewInvalid(errs ...error) *Error      { return newErr(ErrInvalid, errs...) }
func NewNotFound(errs ...error) *Error     { return newErr(ErrNotFound, errs...) }
func NewUnauthorized(errs ...error) *Error { return newErr(ErrUnauthorized, errs...) }
func NewForbidden(errs ...error) *Error    { return newErr(ErrForbidden, errs...) }

func newErr(code ErrCode, errs ...error) *Error {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}
	pcs := make([]uintptr, 32)
	// skip: runtime.Callers, newErr, New/NewXxx の3フレームをスキップ
	n := runtime.Callers(3, pcs)
	return &Error{
		code:  code,
		err:   err,
		stack: pcs[:n],
	}
}

func (e *Error) WithMessage(msg string) *Error {
	e.message = msg
	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%s message=%s", e.code, e.Message())
}

func (e *Error) Code() ErrCode { return e.code }

func (e *Error) Message() string {
	if e.message != "" {
		return e.message
	}
	return defaultMessages[e.code]
}

func (e *Error) Unwrap() error { return e.err }

func (e *Error) StackTrace() string {
	frames := runtime.CallersFrames(e.stack)
	var b strings.Builder
	for {
		f, more := frames.Next()
		fmt.Fprintf(&b, "%s:%d %s\n", f.File, f.Line, f.Function)
		if !more {
			break
		}
	}
	return b.String()
}

func IsNotFound(err error) bool     { return isCode(err, ErrNotFound) }
func IsInvalid(err error) bool      { return isCode(err, ErrInvalid) }
func IsInternal(err error) bool     { return isCode(err, ErrInternal) }
func IsUnauthorized(err error) bool { return isCode(err, ErrUnauthorized) }
func IsForbidden(err error) bool    { return isCode(err, ErrForbidden) }

func isCode(err error, code ErrCode) bool {
	var e *Error
	if stderrors.As(err, &e) {
		return e.code == code
	}
	return false
}

func Is(err, target error) bool { return stderrors.Is(err, target) }

func As(err error, target **Error) bool {
	return stderrors.As(err, target)
}
