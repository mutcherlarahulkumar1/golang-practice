package temp

import (
	"fmt"
	"log"
	"net/http"
)

func server() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "You have hit base URL of the server")

	})

	http.HandleFunc("/users", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "You have hit URL for the Users")

	})

	http.HandleFunc("/objects", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "You have hit URL for the Objects")

	})

	const serverAddr string = "localhost:3000"
	fmt.Println("Server is Listining on Port 3000")
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Error Creating the Server", err)
		return
	}
}
