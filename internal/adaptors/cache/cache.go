package cache

import (
	"article-dispatcher/internal/domain/adaptors/repository"
	"article-dispatcher/internal/domain/models"
	"fmt"
	"sync"
)

type Cache struct {
	lock            *sync.Mutex
	articles        map[string]models.Article
	tagDateIndexMap map[string][]string
}

func NewCache() repository.Repository {
	return &Cache{
		lock:            &sync.Mutex{},
		articles:        make(map[string]models.Article),
		tagDateIndexMap: make(map[string][]string),
	}
}

// Set article data into the cache
func (c Cache) Set(article models.Article) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	art, ok := c.articles[article.Id]
	if ok {
		return fmt.Errorf("error, article id [%s] already exist", art.Id)
	}

	c.articles[article.Id] = article

	// update tag-date index cache
	for _, tag := range article.Tags {
		tagDate := fmt.Sprintf("%s#%s", tag, article.Date)
		_, ok := c.tagDateIndexMap[tagDate]
		if !ok {
			c.tagDateIndexMap[tagDate] = make([]string, 0)
		}
		c.tagDateIndexMap[tagDate] = append(c.tagDateIndexMap[tagDate], article.Id)
	}

	return nil
}

// Get article data from the cache
func (c Cache) Get(id string) (models.Article, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	article, ok := c.articles[id]
	if !ok {
		return article, fmt.Errorf("error, nor article found with id [%s]", id)
	}

	return article, nil
}

// Filter get list of articles satisfying with the filter options
func (c Cache) Filter(tag, date string) (models.Articles, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	articles := make(models.Articles, 0)
	tagDateKey := fmt.Sprintf("%s#%s", tag, date)
	articleIDs, ok := c.tagDateIndexMap[tagDateKey]
	if !ok {
		return articles, fmt.Errorf("error, no article found with tag [%s] - date [%s]", tag, date)
	}

	for _, id := range articleIDs {
		articles = append(articles, c.articles[id])
	}

	return articles, nil
}
