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
	"sync/atomic"
)

var articleIDCount int64

const latestArticleLimit = 10

type cache struct {
	log          logger.Logger
	lock         *sync.RWMutex
	articles     map[string]models.Article
	tagDateIndex map[string][]string
}

func NewCache(l logger.Logger) repository.Repository {
	return &cache{
		log:          l,
		lock:         &sync.RWMutex{},
		articles:     make(map[string]models.Article),
		tagDateIndex: make(map[string][]string),
	}
}

// Set article data into the cache
func (c cache) Set(_ context.Context, article *models.Article) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	id := fmt.Sprintf("%d", atomic.AddInt64(&articleIDCount, 1))
	article.Id = id

	c.articles[id] = *article
	date, err := strconv.Atoi(strings.ReplaceAll(article.Date, "-", ""))
	if err != nil {
		err := fmt.Errorf("error, invalid date format [%s] expected yyyymmdd ", article.Date)
		return InvalidDataError{err}
	}
	// update tag-date index cache
	for _, tag := range article.Tags {
		tagDate := fmt.Sprintf("%s#%d", tag, date)
		_, ok := c.tagDateIndex[tagDate]
		if !ok {
			c.tagDateIndex[tagDate] = make([]string, 0)
		}
		c.tagDateIndex[tagDate] = append(c.tagDateIndex[tagDate], article.Id)
	}

	return nil
}

// Get article data from the cache
func (c cache) Get(_ context.Context, id string) (models.Article, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	article, ok := c.articles[id]
	if !ok {
		err := fmt.Errorf("error, no article found with id [%s]", id)
		return article, DataNotFoundError{err}
	}

	return article, nil
}

// Filter get list of articles satisfying with the filter options
func (c cache) Filter(_ context.Context, tag string, date int) (models.TaggedArticles, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	// temporary tags map to get the counts and avoid deduplication
	tagsMap := make(map[string]struct{})

	// initialize maps
	taggedArticles := models.TaggedArticles{
		Articles:    make([]string, 0),
		RelatedTags: make([]string, 0),
	}
	// restructure tagdate key
	// example `tagName#20230330`
	tagDateKey := fmt.Sprintf("%s#%d", tag, date)
	articleIDs, ok := c.tagDateIndex[tagDateKey]
	if !ok {
		err := fmt.Errorf("error, no article found with tag [%s] - date [%d]", tag, date)
		return taggedArticles, DataNotFoundError{err}
	}

	// get the articles one by one from the articleIDs returned from the indexMap and add the tags into the temporary
	// map defined earlier
	for _, id := range articleIDs {
		article := c.articles[id]
		// write into maps to avoid duplications
		for _, t := range article.Tags {
			// skip the query tag
			if t == tag {
				continue
			}
			tagsMap[t] = struct{}{}
		}
	}

	// get last 10 articles added
	articleTaggedCount := len(articleIDs)
	latestArticleIDs := articleIDs
	// slice only the last required articles
	if articleTaggedCount > latestArticleLimit {
		limit := articleTaggedCount - latestArticleLimit
		latestArticleIDs = articleIDs[limit:]
	}

	// returns only the values of the map into a string slice
	relatedTags := getValueSlice(tagsMap)
	taggedArticles.RelatedTags = relatedTags
	taggedArticles.Articles = latestArticleIDs
	// added one since filtered tag was removed initially from map
	taggedArticles.Count = len(relatedTags) + 1
	taggedArticles.Tag = tag

	return taggedArticles, nil
}

// getValueSlice - get keys of a map as a slice of strings
func getValueSlice(inputMap map[string]struct{}) (outKeySlice []string) {
	outKeySlice = make([]string, 0)
	for key := range inputMap {
		outKeySlice = append(outKeySlice, key)
	}
	return
}
