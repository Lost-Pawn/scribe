package articles

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

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

func GetArticleByID(id uuid.UUID) (*Article, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid article ID: ID is not set")
	}

	path := filepath.Join("data/", id.String()+".json")
	bytes, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("article not found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to read article file: %w", err)
	}

	var article Article
	err = json.Unmarshal(bytes, &article)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal article: %w", err)
	}

	return &article, nil
}

func ListArticles() ([]*Article, error) {
	path := filepath.Join("data/")
	files, err := os.ReadDir(path)
	if errors.Is(err, os.ErrNotExist) {
		return []*Article{}, nil
	} else if err != nil {
		return []*Article{}, fmt.Errorf("failed to read data directory: %w", err)
	}

	var articles []*Article
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		idStr := file.Name()
		if filepath.Ext(idStr) == ".json" {
			idStr = idStr[:len(idStr)-len(".json")]
		} else {
			continue
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("skipping article %s: %v", file.Name(), err)
			continue
		}

		article, err := GetArticleByID(id)
		if err != nil {
			log.Printf("skipping article %s: %v", file.Name(), err)
			continue
		}

		articles = append(articles, article)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})

	return articles, nil
}

func DeleteArticle(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid article ID: ID is not set")
	}
	path := filepath.Join("data/" + id.String() + ".json")

	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("article not found: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	return nil
}
