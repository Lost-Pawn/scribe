package articles

import (
	"time"
	"github.com/google/uuid"
)

type Article struct {
	ID      uuid.UUID 	`json:"id"`
	Title   string 		`json:"title"`
	Content string 		`json:"content"`
	Date   	time.Time 	`json:"date"`
}

