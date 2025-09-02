package middleware

import (
	"context"
	"komodo-auth-api/internal/httpapi/utils"
	"net/http"
	"strings"
)

type ctxKey string
const isUIKey ctxKey = "isUI"

func RequestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		if !utils.IsValidAPIPath(req) {
			http.Error(wtr, "Invalid API version", http.StatusBadRequest)
			return
		}

		isUI := utils.IsUIRequest(req)
		ctx := context.WithValue(req.Context(), isUIKey, isUI)
		req = req.WithContext(ctx)

		// TODO delete when ready
		if true {
			next.ServeHTTP(wtr, req)
			return
		}

		if ((isUI && !uiHeadersValid(wtr, req)) || !apiHeadersValid(wtr, req) || !sharedHeadersValid(wtr, req)) {
			http.Error(wtr, "Invalid request headers", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(wtr, req)
	})
}

func apiHeadersValid(wtr http.ResponseWriter, req *http.Request) bool {
	// ----- Mandatory -----
	if !utils.IsValidBearer(req.Header.Get("Authorization")) {
		http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	if !utils.IsValidContentAcceptType(req.Header.Get("Accept")) {
		http.Error(wtr, "Accept header must be application/json", http.StatusBadRequest)
		return false
	}
	if !utils.IsValidContentAcceptType(req.Header.Get("Content-Type")) {
		http.Error(wtr, "Content-Type header must be application/json", http.StatusBadRequest)
		return false
	}
	if !utils.IsValidContentLength(req.Header.Get("Content-Length")) {
		http.Error(wtr, "Content-Length header is required", http.StatusBadRequest)
		return false
	}
	return true
}

func uiHeadersValid(wtr http.ResponseWriter, req *http.Request) bool {
	// ----- Mandatory -----
	session := req.Header.Get("X-Session-Token")
	csrf := req.Header.Get("X-CSRF")

	// TODO - change to checks for User and/or Cart

	if session == "" || !utils.IsValidSession(session) {
		http.Error(wtr, "Invalid session", http.StatusUnauthorized)
		return false
	}
	if !utils.IsValidCSRF(csrf, session) {
		http.Error(wtr, "CSRF validation failed", http.StatusBadRequest)
		return false
	}
	if !utils.IsValidContentLength(req.Header.Get("Content-Length")) {
		http.Error(wtr, "Content-Length header is required", http.StatusBadRequest)
		return false
	}
	// ----- Optional -----
	cookie := req.Header.Get("Cookie")
	if cookie != "" && !utils.IsValidCookie(cookie) {
		http.Error(wtr, "Cookie header is invalid", http.StatusBadRequest)
		return false
	}
	return true
}

func sharedHeadersValid(wtr http.ResponseWriter, req *http.Request) bool {
	userAgent := req.Header.Get("User-Agent")
	if userAgent != "" && !utils.IsValidUserAgent(userAgent) {
		http.Error(wtr, "User-Agent header is invalid", http.StatusBadRequest)
		return false
	}
	referer := req.Header.Get("Referer")
	if referer != "" && !strings.HasPrefix(referer, "http") {
		http.Error(wtr, "Referer header must be a valid URL", http.StatusBadRequest)
		return false
	}
	cacheControl := req.Header.Get("Cache-Control")
	if cacheControl != "" && !utils.IsValidCacheControl(cacheControl) {
		http.Error(wtr, "Cache-Control header is invalid", http.StatusBadRequest)
		return false
	}
	return true
}