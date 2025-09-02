package errors

import (
	"net/http"
)

type Error struct {
	Code    int
	Message string
}

var (
	ErrBadRequest     = &Error{Code: http.StatusBadRequest, Message: "bad request"}
	ErrUnauthorized   = &Error{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden      = &Error{Code: http.StatusForbidden, Message: "access forbidden"}
	ErrNotFound       = &Error{Code: http.StatusNotFound, Message: "resource not found"}
	ErrMethodNotAllowed = &Error{Code: http.StatusMethodNotAllowed, Message: "method not allowed"}
	ErrUnprocessableEntity = &Error{Code: http.StatusUnprocessableEntity, Message: "unprocessable entity"}
	ErrTooManyRequests = &Error{Code: http.StatusTooManyRequests, Message: "too many requests"}
	ErrInternalServer = &Error{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrBadGateway      = &Error{Code: http.StatusBadGateway, Message: "bad gateway"}
)

func (e *Error) Error() string {
  return e.Message
}

// Predefined error responses for common HTTP status codes
func Error400(msg string) *Error {
	return &Error{Code: http.StatusBadRequest, Message: msg}
}

func Error401(msg string) *Error {
	return &Error{Code: http.StatusUnauthorized, Message: msg}
}

func Error403(msg string) *Error {
	return &Error{Code: http.StatusForbidden, Message: msg}
}

func Error404(msg string) *Error {
	return &Error{Code: http.StatusNotFound, Message: msg}
}

func Error405(msg string) *Error {
	return &Error{Code: http.StatusMethodNotAllowed, Message: msg}
}

func Error422(msg string) *Error {
  return &Error{Code: http.StatusUnprocessableEntity, Message: msg}
}

func Error429(msg string) *Error {
	return &Error{Code: http.StatusTooManyRequests, Message: msg}
}

func Error500(msg string) *Error {
	return &Error{Code: http.StatusInternalServerError, Message: msg}
}

func Error502(msg string) *Error {
	return &Error{Code: http.StatusBadGateway, Message: msg}
}
