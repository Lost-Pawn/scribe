# Scribe
Roadmap.sh Project: https://roadmap.sh/projects/personal-blog

A lightweight, personal blogging platform built with Go. Create, edit, delete, and publish articles with scheduled publishing support. Features admin authentication and file-based storage.

## Features

- **Read Articles**: Browse published articles on the home page
- **Create Articles**: Write new articles with optional scheduled publishing
- **Edit Articles**: Update existing articles' title, content, and publication date
- **Delete Articles**: Remove articles from the platform
- **Scheduled Publishing**: Articles can be scheduled for future publication dates and won't appear until their scheduled time
- **Admin Authentication**: Basic auth protection for administrative operations
- **Responsive Design**: Static CSS styling included
- **File-Based Storage**: Articles stored as JSON files with UUID-based naming

## Project Structure

```
.
├── main.go                 # Server setup, routing, and entry point
├── go.mod / go.sum        # Go module dependencies
├── LICENSE                # Project license
├── README.md              # This file
├── articles/              # Article model and data storage
│   ├── article.go         # Article struct definition
│   └── store.go           # File I/O operations for articles
├── handlers/              # HTTP request handlers
│   ├── guest.go          # Public endpoints (HomeHandler, ArticleHandler)
│   └── admin.go          # Admin endpoints with authentication
├── templates/            # HTML templates for rendering
│   ├── home.html         # Homepage listing articles
│   ├── article.html      # Individual article view
│   ├── dashboard.html    # Admin dashboard
│   └── form.html         # Create/edit article form
├── static/               # Static CSS files
│   ├── home.css
│   └── article.css
├── public/               # Public assets directory
└── data/                 # Article data storage (JSON files)
    └── {uuid}.json       # Individual article files
```

## Tech Stack

- **Language**: Go 1.27.1
- **UUID Generation**: github.com/google/uuid v1.6.0
- **Server**: Built-in Go `net/http` package
- **Templating**: Go `html/template`
- **Storage**: File-based JSON (no database required)
- **Frontend**: HTML + CSS

## Getting Started

### Prerequisites

- Go 1.27.1 or later

### Installation

1. Clone the repository:
```bash
git clone https://github.com/lost-pawn/scribe.git
cd scribe
```

2. Install dependencies:
```bash
go mod download
```

## Running the Server

The server requires two environment variables for admin authentication:

```bash
export ADMIN_USER=your_username
export ADMIN_PASS=your_password
go run main.go
```

The server will launch on `http://localhost:8080`

## API Routes

### Public Routes

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/` | Home page with published articles |
| GET | `/article/{id}` | View individual article by UUID |
| GET/POST | `/static/*` | Static CSS files |
| GET/POST | `/public/*` | Public assets |

### Admin Routes (Requires Basic Auth)

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/admin` | Admin dashboard with all articles |
| GET | `/admin/new` | New article form |
| POST | `/admin/new` | Create new article |
| GET | `/admin/edit/{id}` | Edit article form |
| POST | `/admin/edit/{id}` | Update article |
| POST | `/admin/delete/{id}` | Delete article |

## Article Data Format

Articles are stored as JSON files in the `data/` directory. Each file is named with the article's UUID.

Example: `data/550e8400-e29b-41d4-a716-446655440001.json`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "title": "Article Title",
  "content": "Article content goes here...",
  "date": "2026-10-09T16:45:00Z"
}
```

### Article Fields

- **id**: UUID unique identifier for the article
- **title**: Article title (max 200 characters)
- **content**: Article body text (required)
- **date**: ISO8601 timestamp - publication date and time in UTC

## Publishing & Scheduling

- Articles with a **past or current date** are **immediately published** and visible on the home page
- Articles with a **future date** are **scheduled** and hidden until that date arrives
- Scheduled articles appear in the admin dashboard marked as "Scheduled"
- No separate "publish" action needed—publication is automatic based on the date

## Admin Authentication

Admin endpoints are protected with HTTP Basic Authentication:
- Set `ADMIN_USER` and `ADMIN_PASS` environment variables
- Credentials are required to access `/admin` and related routes
- Authentication is enforced by the `RequireAuth` middleware

## Server Configuration

The server listens on `localhost:8080` with the following timeouts:
- Read Header Timeout: 5 seconds
- Read Timeout: 10 seconds
- Write Timeout: 10 seconds
- Idle Timeout: 60 seconds

## License

See LICENSE file for details.
