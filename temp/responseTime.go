package temp

import (
	"net/http"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//start := time.Now()
		next.ServeHTTP(writer, request)
	})
}
