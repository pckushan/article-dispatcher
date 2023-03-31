package services

import (
	"article-dispatcher/internal/domain/models"
	"context"
)

// ArticleService create the article in the system
type ArticleService interface {
	Create(ctx context.Context, article models.Article) error
}
