package main

import (
	"fmt"
	"golang/internal/api/handlers"
	"golang/internal/api/middlewares"
	"golang/sqlconnect"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "THis is the Base Route of the API End Point")
	fmt.Println("Hello Base Route")
}

func anotherRootHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("This is another way of passing into response"))
	fmt.Println("Hello Base Route")
}

func main() {
	godotenv.Load()
	_, err := sqlconnect.ConnectToDB()
	if err != nil {
		fmt.Println("Error Connecting ")
		return
	}
	port := ":3000"
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/anotherWay", anotherRootHandler)

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)
	mux.HandleFunc("PUT /teachers/", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)

	mux.HandleFunc("/students/", handlers.StudentsHanlder)

	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	fmt.Println("The Server is running on Port : ", port)

	err = http.ListenAndServe(port, middlewares.Cors(middlewares.SecurityHeaders(mux)))
	if err != nil {
		log.Fatalf("Error Starting Server on Port : ", port, err)
	}
}
