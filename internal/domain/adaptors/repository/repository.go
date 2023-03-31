package repository

import "article-dispatcher/internal/domain/models"

// Repository high-level methods to the repository function
// Set - insert the value into the repository
// Get - retrieve data from repository
// Filter - fetch conditioned articles data from repository
type Repository interface {
	Set(article models.Article) error
	Get(id string) (models.Article, error)
	Filter(tag, date string) (models.Articles, error)
}
