package handlers

import (
	"fmt"
	"log"
	"net/http"
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
