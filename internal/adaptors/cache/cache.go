package cache

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/adaptors/repository"
	"article-dispatcher/internal/domain/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const LatestArticleLimit = 10

type Cache struct {
	log             logger.Logger
	lock            *sync.Mutex
	articles        map[string]models.Article
	tagDateIndexMap map[string][]string
}

func NewCache(l logger.Logger) repository.Repository {
	return &Cache{
		log:             l,
		lock:            &sync.Mutex{},
		articles:        make(map[string]models.Article),
		tagDateIndexMap: make(map[string][]string),
	}
}

// Set article data into the cache
func (c Cache) Set(ctx context.Context, article *models.Article) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	art, ok := c.articles[article.Id]
	if ok {
		c.log.Error(fmt.Sprintf("article id [%s] already exist", art.Id))
		err := fmt.Errorf("error, article id [%s] already exist", art.Id)
		return DuplicateError{err}
	}

	c.articles[article.Id] = *article
	date, err := strconv.Atoi(strings.ReplaceAll(article.Date, "-", ""))
	if err != nil {
		err := fmt.Errorf("error, invalid date format [%s] expected yyyymmdd ", article.Date)
		return InvalidDataError{err}
	}
	// update tag-date index cache
	for _, tag := range article.Tags {
		tagDate := fmt.Sprintf("%s#%d", tag, date)
		_, ok := c.tagDateIndexMap[tagDate]
		if !ok {
			c.tagDateIndexMap[tagDate] = make([]string, 0)
		}
		c.tagDateIndexMap[tagDate] = append(c.tagDateIndexMap[tagDate], article.Id)
	}

	return nil
}

// Get article data from the cache
func (c Cache) Get(ctx context.Context, id string) (models.Article, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	article, ok := c.articles[id]
	if !ok {
		err := fmt.Errorf("error, no article found with id [%s]", id)
		return article, DataNotFoundError{err}
	}

	return article, nil
}

// Filter get list of articles satisfying with the filter options
func (c Cache) Filter(ctx context.Context, tag string, date int) (models.TaggedArticles, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	// temporary maps
	articlesMap := make(map[string]int)
	tagsMap := make(map[string]int)

	taggedArticles := models.TaggedArticles{
		Tag:         tag,
		Articles:    make([]string, 0),
		RelatedTags: make([]string, 0),
	}
	tagDateKey := fmt.Sprintf("%s#%d", tag, date)
	articleIDs, ok := c.tagDateIndexMap[tagDateKey]
	if !ok {
		err := fmt.Errorf("error, no article found with tag [%s] - date [%d]", tag, date)
		return taggedArticles, DataNotFoundError{err}
	}

	for _, id := range articleIDs {
		article := c.articles[id]
		// write into maps to avoid duplications
		articlesMap[article.Id]++
		for _, t := range article.Tags {
			// skip the query tag
			if t == tag {
				continue
			}
			tagsMap[t]++
		}
	}

	// get last 10 articles added
	articleTaggedCount := len(articleIDs)
	latestArticleIDs := articleIDs
	// slice only the last required articles
	if articleTaggedCount > LatestArticleLimit {
		limit := articleTaggedCount - LatestArticleLimit
		latestArticleIDs = articleIDs[limit:]
	}

	relatedTags := getValueSlice(tagsMap)
	taggedArticles.RelatedTags = relatedTags
	taggedArticles.Articles = latestArticleIDs
	// added one since filtered tag was removed initially from map
	taggedArticles.Count = len(relatedTags) + 1

	return taggedArticles, nil
}

func getValueSlice(inputMap map[string]int) (outKeySlice []string) {
	outKeySlice = make([]string, 0)
	for key := range inputMap {
		outKeySlice = append(outKeySlice, key)
	}
	return
}
