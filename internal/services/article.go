package services

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/adaptors/repository"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/domain/services"
	"context"
	"fmt"
)

type ArticleService struct {
	log  logger.Logger
	repo repository.Repository
}

func NewArticleService(l logger.Logger, repo repository.Repository) services.ArticleService {
	return &ArticleService{
		log:  l,
		repo: repo,
	}
}

func (as ArticleService) Create(ctx context.Context, article models.Article) error {
	err := as.repo.Set(ctx, article)
	if err != nil {
		as.log.Error(fmt.Sprintf("article service, create article error due to %s", err))
	}
	return err
}

func (as ArticleService) Get(ctx context.Context, id string) (models.Article, error) {
	article, err := as.repo.Get(ctx, id)
	if err != nil {
		as.log.Error(fmt.Sprintf("article service, get article error due to %s", err))
	}
	return article, err
}
func (as ArticleService) Filter(ctx context.Context, tag string, date int) (models.TaggedArticles, error) {
	taggedArticles, err := as.repo.Filter(ctx, tag, date)
	if err != nil {
		as.log.Error(fmt.Sprintf("article service, filter articles error due to %s", err))
	}

	return taggedArticles, err
}
