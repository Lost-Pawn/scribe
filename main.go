package main

import (
	"fmt"
	"net/http"
)

func Server() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Welcome to the Home Page!")
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	Server()
}
