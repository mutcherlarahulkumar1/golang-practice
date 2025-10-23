package handlers

import "net/http"

func ExecsHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writer.Write([]byte("THis GET Call for execs"))

	case http.MethodPost:
		writer.Write([]byte("THis POST Call for execs"))

	case http.MethodPut:
		writer.Write([]byte("THis PUT Call for execs"))

	case http.MethodPatch:
		writer.Write([]byte("THis PATCH Call for execs"))

	case http.MethodDelete:
		writer.Write([]byte("THis DELETE Call for execs"))
	}
}
