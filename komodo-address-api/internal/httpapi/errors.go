package httpapi

import "net/http"

func ErrorObj(msg string) map[string]string {
	return map[string]string{"error": msg}
}

// Predefined error responses for common HTTP status codes
func Error500() (int, map[string]string) {
	return http.StatusInternalServerError, ErrorObj("internal server error")
}

func Error404() (int, map[string]string) {
	return http.StatusNotFound, ErrorObj("resource not found")
}

func Error403() (int, map[string]string) {
	return http.StatusForbidden, ErrorObj("access forbidden")
}

func Error401() (int, map[string]string) {
	return http.StatusUnauthorized, ErrorObj("unauthorized")
}

func Error400() (int, map[string]string) {
	return http.StatusBadRequest, ErrorObj("bad request")
}

func Error422(msg string) (int, map[string]string) {
	return http.StatusUnprocessableEntity, ErrorObj(msg)
}
