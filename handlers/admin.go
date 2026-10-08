package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"scribe/articles"
	"time"
	"github.com/google/uuid"
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

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	dateStr := r.FormValue("date")

	data := FormData{
		Title:   title,
		Content: content,
		Date:    dateStr,
		Today:   time.Now().UTC().Format("2006-01-02"),
		Action:  "/admin/new",
		Heading: "Write a new article",
	}

	if title == "" {
		data.TitleError = "Title is required."
	}

	if len(title) > 200 {
		data.TitleError = "Maximum title length is 200 characters."
	}

	if content == "" {
		data.ContentError = "Content is required."
	}

	if dateStr == "" {
		data.DateError = "Date is required."
	}

	var articleDate time.Time

	if data.DateError == "" {
		articleDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			data.DateError = "Invalid date format. Use YYYY-MM-DD."
		}
	}

	if data.TitleError == "" && data.ContentError == "" && data.DateError == "" {
		article := &articles.Article{
			ID: 	GenerateID(),
			Title:   title,
			Content: content,
			Date:    articleDate,
		}

		log.Printf("Article ID before saving: %v", article.ID)

		err = articles.SaveArticle(article)
		if err != nil {
			log.Printf("error saving article: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data.Success = true

		data.Title = ""
		data.Content = ""
		data.Date = data.Today
	}

	var buf bytes.Buffer

	if err := newArticleFormTemplate.Execute(&buf, data); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func GenerateID() uuid.UUID {
	return uuid.New()
}