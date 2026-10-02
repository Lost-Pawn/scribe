package articles

import (
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
)

func SaveArticle(article Article) error {
	path := filepath.Join("data/", article.ID.String() +".json")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(article)
}

func GetArticleByID(id string) (Article, error) {
	if id == "" {
		return Article{}, fmt.Errorf("invalid ID")
	}

	path := filepath.Join("data/", id +".json")
	file, err := os.Open(path)
	if err != nil {
		return Article{}, err
	}
	defer file.Close()
	
	var article Article
	err = json.NewDecoder(file).Decode(&article)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

