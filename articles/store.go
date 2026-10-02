package articles

import (
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
	"github.com/google/uuid"
)

func SaveArticle(article *Article) error {
	if article == nil || article.ID == uuid.Nil {
		return fmt.Errorf("invalid article: article is nil or ID is not set")
	}
	
	_, err := os.Stat("data")
	if os.IsNotExist(err) {
		err = os.Mkdir("data", 0755)
		if err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	bytes, err := json.MarshalIndent(article, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	err = os.WriteFile(filepath.Join("data/", article.ID.String()+".json"), bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write article to file: %w", err)
	}

	return nil
}
