package cors

import "net/http"

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO: implement CORS handling
		next.ServeHTTP(wtr, req)
	})
}