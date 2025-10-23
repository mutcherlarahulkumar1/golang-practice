package temp

import (
	"fmt"
	"io"
	"net/http"
)

func client() {
	client := &http.Client{}

	resp, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")

	if err != nil {
		fmt.Println("Error Occured")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error Occured")
		return
	}
	fmt.Println("Response : ", string(body))
}
