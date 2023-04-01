package repository

import (
	"article-dispatcher/internal/domain/models"
	"context"
)

// Repository high-level methods to the repository function
// Set - insert the value into the repository
// Get - retrieve data from repository
// Filter - fetch conditioned articles data from repository
type Repository interface {
	Set(ctx context.Context, article *models.Article) error
	Get(ctx context.Context, id string) (models.Article, error)
	Filter(ctx context.Context, tag string, date int) (models.TaggedArticles, error)
}
