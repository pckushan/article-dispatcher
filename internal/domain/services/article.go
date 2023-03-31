package services

import (
	"article-dispatcher/internal/domain/models"
	"context"
)

// ArticleService create the article in the system
type ArticleService interface {
	Create(ctx context.Context, article models.Article) error
	Get(ctx context.Context, id string) (models.Article, error)
	Filter(ctx context.Context, tag, date string) (models.TaggedArticles, error)
}
