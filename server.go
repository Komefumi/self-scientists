package main

import (
	"fmt"
	"net/http"
	"self-scientists/config"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		fmt.Fprintf(w, "Welcome to new server very much!")
	})

	http.HandleFunc("/users", registrationHandler)
	http.ListenAndServe(":"+config.ServerPortString, nil)
}
