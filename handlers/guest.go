package handlers

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"scribe/articles"
	"html/template"
	"bytes"

	"time"
)

var homeTemplate = template.Must(template.ParseFiles("templates/home.html"))
var articleTemplate = template.Must(template.ParseFiles("templates/article.html"))

func isPublished(article *articles.Article, now time.Time) bool {
	return !article.Date.After(now)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	articleslist, err := articles.ListArticles()
	if err != nil {
		log.Printf("Error occurred while listing articles: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	var published []*articles.Article
	for _, article := range articleslist {
		if !isPublished(article, now) {
			continue
		}
		published = append(published, article)
	}
	var buf bytes.Buffer

	if err := homeTemplate.Execute(&buf, published); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := articles.GetArticleByID(id)
	if errors.Is(err, os.ErrNotExist) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		log.Printf("failed to load article %s: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	if !isPublished(article, now) {
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer
	if err := articleTemplate.Execute(&buf, article); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}
