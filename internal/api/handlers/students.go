package handlers

import "net/http"

func StudentsHanlder(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writer.Write([]byte("THis GET Call for students"))

	case http.MethodPost:
		writer.Write([]byte("THis POST Call for students"))

	case http.MethodPut:
		writer.Write([]byte("THis PUT Call for students"))

	case http.MethodPatch:
		writer.Write([]byte("THis PATCH Call for students"))

	case http.MethodDelete:
		writer.Write([]byte("THis DELETE Call for students"))
	}
}
