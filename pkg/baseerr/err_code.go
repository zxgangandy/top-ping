package baseerr

import (
	"fmt"
	"net/http"
)

var (
	// 预定义错误
	// Common errors
	Success               = NewError(0, "Success")
	ErrInternalServer     = NewError(10001, "Internal server error")
	ErrBind               = NewError(10002, "Bind request error")
	ErrInvalidParam       = NewError(10003, "Invalid params")
	ErrSignParam          = NewError(10004, "Invalid sign")
	ErrValidation         = NewError(10005, "Validation failed")
	ErrDatabase           = NewError(10006, "Database error")
	ErrToken              = NewError(10007, "Gen token error")
	ErrInvalidToken       = NewError(10108, "Invalid token")
	ErrTokenTimeout       = NewError(10109, "Token timeout")
	ErrTooManyRequests    = NewError(10110, "Too many request")
	ErrInvalidTransaction = NewError(10111, "Invalid transaction")
	ErrEncrypt            = NewError(10112, "Encrypting the user password error")
	ErrLimitExceed        = NewError(10113, "Beyond limit")
	ErrServiceUnavailable = NewError(10114, "Service Unavailable")
)

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]struct{}{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code %d is exsit, please change one", code))
	}
	codes[code] = struct{}{}
	return &Error{code: code, msg: msg}
}

func (e Error) Error() string {
	return fmt.Sprintf("code：%d, msg:：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}

	return &newError
}

// StatusCode trans err code to http status code
func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ErrInternalServer.Code():
		return http.StatusInternalServerError
	case ErrInvalidParam.Code():
		return http.StatusBadRequest
	case ErrToken.Code():
		fallthrough
	case ErrInvalidToken.Code():
		fallthrough
	case ErrTokenTimeout.Code():
		return http.StatusUnauthorized
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	case ErrServiceUnavailable.Code():
		return http.StatusServiceUnavailable
	}

	return http.StatusOK
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

// DecodeErr 对错误进行解码，返回错误code和错误提示
func DecodeErr(err error) (int, string) {
	if err == nil {
		return Success.code, Success.msg
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Error:
		return typed.code, typed.msg
	default:
	}

	return ErrInternalServer.Code(), err.Error()
}
