package main

import (
	"net/http"
	"scribe/handlers"
)

func Server() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /article/{id}", handlers.ArticleHandler)
	mux.HandleFunc("GET /admin", handlers.AdminDashboard)  // unprotected route for now, will add authentication later
	mux.HandleFunc("GET /admin/new", handlers.NewArticleForm) // unprotected route for now, will add authentication later
	mux.HandleFunc("POST /admin/new", handlers.CreateArticle) // unprotected route for now, will add authentication later
	mux.HandleFunc("GET /admin/edit/{id}", handlers.GetArticleForm) // unprotected route for now, will add authentication later
	mux.HandleFunc("POST /admin/edit/{id}", handlers.EditArticle) // unprotected route for now, will add authentication later
	mux.HandleFunc("POST /admin/delete/{id}", handlers.DeleteArticle) // unprotected route for now, will add authentication later
	
	err := http.ListenAndServe("localhost:8001", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	Server()
}
