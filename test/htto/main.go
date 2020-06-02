package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("in\n");
		writer.Write([]byte("hello"))
	})

	http.ListenAndServe(":9099", nil)
}
