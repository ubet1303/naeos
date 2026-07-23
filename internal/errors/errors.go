package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type ErrorCode string

const (
	ErrParse                 ErrorCode = "PARSE_ERROR"
	ErrValidation            ErrorCode = "VALIDATION_ERROR"
	ErrCloud                 ErrorCode = "CLOUD_ERROR"
	ErrPlugin                ErrorCode = "PLUGIN_ERROR"
	ErrAuth                  ErrorCode = "AUTH_ERROR"
	ErrPipeline              ErrorCode = "PIPELINE_ERROR"
	ErrConfig                ErrorCode = "CONFIG_ERROR"
	ErrDatabase              ErrorCode = "DATABASE_ERROR"
	ErrNetwork               ErrorCode = "NETWORK_ERROR"
	ErrInternal              ErrorCode = "INTERNAL_ERROR"
	ErrNotFound              ErrorCode = "NOT_FOUND"
	ErrConflict              ErrorCode = "CONFLICT"
	ErrTimeout               ErrorCode = "TIMEOUT"
	ErrRateLimit             ErrorCode = "RATE_LIMIT"
	ErrPermDenied            ErrorCode = "PERMISSION_DENIED"
	ErrCanceled              ErrorCode = "CANCELED"
	ErrFeatureNotImplemented ErrorCode = "NOT_IMPLEMENTED"
)

type StackFrame struct {
	File string
	Line int
	Func string
}

type NaeosError struct {
	Code      ErrorCode
	Message   string
	Inner     error
	stack     []StackFrame
	context   map[string]any
	retryable bool
}

func (e *NaeosError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Inner)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *NaeosError) Unwrap() error {
	return e.Inner
}

func (e *NaeosError) Is(target error) bool {
	if t, ok := target.(*NaeosError); ok {
		return e.Code == t.Code
	}
	return false
}

func (e *NaeosError) Stack() []StackFrame {
	return e.stack
}

func (e *NaeosError) Context() map[string]any {
	return e.context
}

func (e *NaeosError) IsRetryable() bool {
	return e.retryable
}

func (e *NaeosError) WithContext(key string, value any) *NaeosError {
	if e.context == nil {
		e.context = make(map[string]any)
	}
	e.context[key] = value
	return e
}

func (e *NaeosError) WithRetry() *NaeosError {
	e.retryable = true
	return e
}

func captureStack(skip int) []StackFrame {
	var frames []StackFrame
	pcs := make([]uintptr, 10)
	n := runtime.Callers(skip, pcs)
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pcs[i])
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pcs[i])
		frames = append(frames, StackFrame{
			File: file,
			Line: line,
			Func: fn.Name(),
		})
	}
	return frames
}

func New(code ErrorCode, msg string) *NaeosError {
	return &NaeosError{
		Code:    code,
		Message: msg,
		stack:   captureStack(3),
	}
}

func Wrap(code ErrorCode, msg string, err error) *NaeosError {
	return &NaeosError{
		Code:    code,
		Message: msg,
		Inner:   err,
		stack:   captureStack(3),
	}
}

func Wrapf(err error, code ErrorCode, format string, args ...any) *NaeosError {
	return &NaeosError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Inner:   err,
		stack:   captureStack(3),
	}
}

func Is(err error, code ErrorCode) bool {
	for err != nil {
		var ne *NaeosError
		if As(err, &ne) && ne.Code == code {
			return true
		}
		uw, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = uw.Unwrap()
	}
	return false
}

func IsAny(err error, codes ...ErrorCode) bool {
	for _, code := range codes {
		if Is(err, code) {
			return true
		}
	}
	return false
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func IsRetryable(err error) bool {
	var ne *NaeosError
	if As(err, &ne) && ne.retryable {
		return true
	}
	if u, ok := err.(interface{ Unwrap() error }); ok {
		return IsRetryable(u.Unwrap())
	}
	return false
}

func CodeOf(err error) ErrorCode {
	var ne *NaeosError
	if As(err, &ne) {
		return ne.Code
	}
	return ""
}

func ContextOf(err error) map[string]any {
	var ne *NaeosError
	if As(err, &ne) {
		return ne.context
	}
	return nil
}

type ErrorGroup struct {
	errors []error
}

func Group(errs ...error) *ErrorGroup {
	var valid []error
	for _, e := range errs {
		if e != nil {
			valid = append(valid, e)
		}
	}
	return &ErrorGroup{errors: valid}
}

func (g *ErrorGroup) Error() string {
	if len(g.errors) == 0 {
		return ""
	}
	if len(g.errors) == 1 {
		return g.errors[0].Error()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d errors:", len(g.errors))
	for _, e := range g.errors {
		fmt.Fprintf(&b, "\n  - %v", e)
	}
	return b.String()
}

func (g *ErrorGroup) Errors() []error {
	return g.errors
}

func (g *ErrorGroup) Unwrap() []error {
	return g.errors
}

func (g *ErrorGroup) Is(target error) bool {
	for _, e := range g.errors {
		if errors.Is(e, target) {
			return true
		}
		if sub, ok := e.(*ErrorGroup); ok && sub.Is(target) {
			return true
		}
	}
	return false
}

func (g *ErrorGroup) HasCode(code ErrorCode) bool {
	for _, e := range g.errors {
		if Is(e, code) {
			return true
		}
	}
	return false
}

func (g *ErrorGroup) Codes() []ErrorCode {
	seen := make(map[ErrorCode]bool)
	var codes []ErrorCode
	for _, e := range g.errors {
		c := CodeOf(e)
		if c != "" && !seen[c] {
			seen[c] = true
			codes = append(codes, c)
		}
	}
	return codes
}

func (g *ErrorGroup) Len() int {
	count := 0
	for _, e := range g.errors {
		var sub *ErrorGroup
		if errors.As(e, &sub) {
			count += sub.Len()
		} else {
			count++
		}
	}
	return count
}

func (g *ErrorGroup) Empty() bool {
	return len(g.errors) == 0
}

var (
	ErrNotConnected     = New(ErrNetwork, "unable to connect to the remote host; check network connectivity and host configuration")
	ErrInvalidSpec      = New(ErrValidation, "invalid specification; verify the input file format and required fields")
	ErrPluginNotFound   = New(ErrPlugin, "plugin not found; ensure the plugin is installed and the name is correct")
	ErrDeployFailed     = New(ErrCloud, "cloud deployment failed; check provider status and deployment logs for details")
	ErrTimedOut         = New(ErrTimeout, "operation timed out; increase the timeout value or check service availability").WithRetry()
	ErrRateLimited      = New(ErrRateLimit, "rate limit exceeded; wait before retrying or reduce request frequency").WithRetry()
	ErrUnauthorized     = New(ErrAuth, "authentication failed; verify credentials or re-authenticate")
	ErrForbidden        = New(ErrPermDenied, "access denied; the required permission is not granted for this resource")
	ErrAlreadyExists    = New(ErrConflict, "resource already exists; use a different name or delete the existing resource")
	ErrInvalidConfig    = New(ErrConfig, "invalid configuration; check the config file for missing or incorrect values")
	ErrInternalFailed   = New(ErrInternal, "internal error; please report this issue with the stack trace")
	ErrNotImplemented   = New(ErrFeatureNotImplemented, "feature not implemented; check the documentation for supported features")
	ErrDependencyCycle  = New(ErrPipeline, "dependency cycle detected; review module dependencies and remove circular references")
	ErrNotFoundSentinel = New(ErrNotFound, "resource not found; verify the resource name and ensure it exists")
)
