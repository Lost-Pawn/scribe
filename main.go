package main

import (
	"log"
	"net/http"
	"os"
	"scribe/handlers"
)

func Server() {
	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")

	if adminUser == "" || adminPass == "" {
		log.Fatal("ADMIN_USER and ADMIN_PASS environment variables must be set")
	}

	auth := func(handler http.HandlerFunc) http.HandlerFunc {
		return handlers.RequireAuth(adminUser, adminPass, handler)
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /article/{id}", handlers.ArticleHandler)
	mux.HandleFunc("GET /admin", auth(handlers.AdminDashboard))
	mux.HandleFunc("GET /admin/new", auth(handlers.NewArticleForm))
	mux.HandleFunc("POST /admin/new", auth(handlers.CreateArticle))
	mux.HandleFunc("GET /admin/edit/{id}", auth(handlers.GetArticleForm))
	mux.HandleFunc("POST /admin/edit/{id}", auth(handlers.EditArticle))
	mux.HandleFunc("POST /admin/delete/{id}", auth(handlers.DeleteArticle))

	err := http.ListenAndServe("localhost:8001", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	Server()
}
