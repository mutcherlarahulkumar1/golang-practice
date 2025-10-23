package middlewares

import (
	"net/http"
)

var allowedOrigin = []string{
	"http://localhost:3000",
	"http://localhost:3001",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if !isOriginAllowed(origin) {
			http.Error(writer, "Request Blocked By CORS", http.StatusForbidden)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func isOriginAllowed(origin string) bool {
	for _, val := range allowedOrigin {
		if val == origin {
			return true
		}
	}
	return false
}
