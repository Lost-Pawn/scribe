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

type FormData struct {
    Title        string
    Content      string
    Date         string
    Today        string

    Error        string
    TitleError   string
    ContentError string
    DateError    string

    Success   bool
    Submitting bool
    CSRFToken string

    Action  string
    Heading string
}

var dashboardTemplate = template.Must(template.ParseFiles("templates/dashboard.html"))
var newArticleFormTemplate = template.Must(template.ParseFiles("templates/form.html"))


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

func NewArticleForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	today := time.Now().UTC().Format("2006-01-02")

	data := FormData{
		Title:        "",
		Content:      "",
		Date:         today,
		Today:        today,
		Error:        "",
		TitleError:   "",
		ContentError: "",
		DateError:    "",
		Success:      false,
		Submitting:   false,
		CSRFToken:    "",
		Action:       "/admin/new",
		Heading:      "Write a new article",
	}

	var buf bytes.Buffer
	if err := newArticleFormTemplate.Execute(&buf, data); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}