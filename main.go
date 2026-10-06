package main

import (
	"net/http"
	"scribe/handlers"
)

func Server() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.HomeHandler)


	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	Server()
}
