package handlers

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"scribe/articles"

	"time"
)

func isPublished(article *articles.Article, now time.Time) bool {
	return !article.Date.After(now)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	articleslist, err := articles.ListArticles()
	if err != nil {
		log.Printf("Error occurred while listing articles: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	for _, article := range articleslist {
		if !isPublished(article, now) {
			continue
		}

		fmt.Fprintf(w, "Title: %s\n", article.Title)
		fmt.Fprintf(w, "Date: %s\n", article.Date.Format("2006-01-02 15:04:05"))
		fmt.Fprintln(w, "-------------------------")
	}
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

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

	fmt.Fprintf(w, "Title: %s\n", article.Title)
	fmt.Fprintf(w, "Date: %s\n", article.Date.Format("2006-01-02"))
	fmt.Fprintln(w, "-------------------------")
	fmt.Fprintf(w, "%s\n", article.Content)
}
