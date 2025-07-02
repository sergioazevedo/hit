package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting server on port 8082")
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello world\n")
	})
	http.ListenAndServe(":8082", nil)
}
