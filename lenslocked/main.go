package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", nil)
}
