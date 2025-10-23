package middlewares

import (
	"fmt"
	"net/http"
)

// BASIC Middleware skeleton
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Mocking Security headers")
		next.ServeHTTP(writer, request)
	})
}
