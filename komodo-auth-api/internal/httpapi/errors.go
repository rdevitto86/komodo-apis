package httpapi

import (
	"net/http"
)

type Error struct {
	StatusCode    int
	ErrorCode     string
	Message       string
}

var (
	ErrBadRequest     = &Error{StatusCode: http.StatusBadRequest, ErrorCode: "bad_request", Message: "bad request"}
	ErrUnauthorized   = &Error{StatusCode: http.StatusUnauthorized, ErrorCode: "unauthorized", Message: "unauthorized"}
	ErrForbidden      = &Error{StatusCode: http.StatusForbidden, ErrorCode: "access_forbidden", Message: "access forbidden"}
	ErrNotFound       = &Error{StatusCode: http.StatusNotFound, ErrorCode: "resource_not_found", Message: "resource not found"}
	ErrMethodNotAllowed = &Error{StatusCode: http.StatusMethodNotAllowed, ErrorCode: "method_not_allowed", Message: "method not allowed"}
	ErrUnprocessableEntity = &Error{StatusCode: http.StatusUnprocessableEntity, ErrorCode: "unprocessable_entity", Message: "unprocessable entity"}
	ErrTooManyRequests = &Error{StatusCode: http.StatusTooManyRequests, ErrorCode: "too_many_requests", Message: "too many requests"}
	ErrInternalServer = &Error{StatusCode: http.StatusInternalServerError, ErrorCode: "internal_server_error", Message: "internal server error"}
	ErrBadGateway      = &Error{StatusCode: http.StatusBadGateway, ErrorCode: "bad_gateway", Message: "bad gateway"}
)

func ToJSON(err *Error) (int, map[string]string) {
  return err.StatusCode, map[string]string{"error": err.Message}
}

func (e *Error) Error() string {
  return e.Message
}

// Predefined error responses for common HTTP status codes
func Error400(msg string) *Error {
	return &Error{StatusCode: http.StatusBadRequest, ErrorCode: "bad_request", Message: msg}
}

func Error401(msg string) *Error {
	return &Error{StatusCode: http.StatusUnauthorized, ErrorCode: "unauthorized", Message: msg}
}

func Error403(msg string) *Error {
	return &Error{StatusCode: http.StatusForbidden, ErrorCode: "access_forbidden", Message: msg}
}

func Error404(msg string) *Error {
	return &Error{StatusCode: http.StatusNotFound, ErrorCode: "resource_not_found", Message: msg}
}

func Error405(msg string) *Error {
	return &Error{StatusCode: http.StatusMethodNotAllowed, ErrorCode: "method_not_allowed", Message: msg}
}

func Error422(msg string) *Error {
  return &Error{StatusCode: http.StatusUnprocessableEntity, ErrorCode: "unprocessable_entity", Message: msg}
}

func Error429(msg string) *Error {
	return &Error{StatusCode: http.StatusTooManyRequests, ErrorCode: "too_many_requests", Message: msg}
}

func Error500(msg string) *Error {
	return &Error{StatusCode: http.StatusInternalServerError, ErrorCode: "internal_server_error", Message: msg}
}

func Error502(msg string) *Error {
	return &Error{StatusCode: http.StatusBadGateway, ErrorCode: "bad_gateway", Message: msg}
}