package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"scribe/articles"
	"time"
)

type AdminRow struct {
	Article   *articles.Article
	Scheduled bool
}

var dashboardTemplate = template.Must(template.ParseFiles("templates/dashboard.html"))

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	articlesList, err := articles.ListArticles()
	if err != nil {
		log.Printf("Error occurred while listing articles: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	rowSlice := make([]AdminRow, 0, len(articlesList))
	for _, article := range articlesList {
		isScheduled := !isPublished(article, now)
		rowSlice = append(rowSlice, AdminRow{Article: article, Scheduled: isScheduled})

	}
	var buf bytes.Buffer
	if err := dashboardTemplate.Execute(&buf, rowSlice); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}
