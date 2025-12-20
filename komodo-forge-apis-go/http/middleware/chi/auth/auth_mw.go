package auth

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO: implement JWT handling
		next.ServeHTTP(wtr, req)
	})
}

	// return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
	// 	ctx := req.Context()

	// 	// Authorization: Bearer <token>
	// 	if auth := req.Header.Get("Authorization"); auth != "" {
	// 		parts := strings.SplitN(auth, " ", 2)
	// 		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
	// 			logger.Error("invalid Authorization header format", req)
	// 			http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
	// 			return
	// 		}

	// 		token := strings.TrimSpace(parts[1])
	// 		ok, err := jwt.VerifyJWT(token)

	// 		if err != nil || !ok {
	// 			logger.Error("invalid bearer token: " + err.Error(), req)
	// 			http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
	// 			return
	// 		}
	// 		ctx = context.WithValue(ctx, AuthValidCtxKey, true)
	// 	}

	// 	// X-Session-Token: optional header for session tokens
	// 	if sess := req.Header.Get("X-Session-Token"); sess != "" {
	// 		ok, err := jwt.VerifyJWT(sess)
	// 		if err != nil || !ok {
	// 			logger.Error("invalid session token: " + err.Error(), req)
	// 			http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
	// 			return
	// 		}
	// 		ctx = context.WithValue(ctx, SessionValidCtxKey, true)
	// 	}

	// 	next.ServeHTTP(wtr, req.WithContext(ctx))
	// })
// }
