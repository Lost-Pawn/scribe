package articles

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}
