package server

import (
	"fmt"
	"net/http"
	"self-scientists/config"
)

func SetupServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		fmt.Fprintf(w, "Welcome to new server very much!")
	})

	http.HandleFunc("/users", registrationHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/test", tokenCheckMiddleware(testHandler))
	http.HandleFunc("/threads", tokenCheckMiddleware(threadsHandler))
	http.HandleFunc("/threads/count", tokenCheckMiddleware(testHandler))
	http.ListenAndServe(":"+config.ServerPortString, nil)
}
