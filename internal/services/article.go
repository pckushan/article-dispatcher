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
	err := as.repo.Set(article)
	if err != nil {
		as.log.Error(fmt.Sprintf("article service error due to %s", err))
	}
	return err
}

func (as ArticleService) Get(ctx context.Context, id string) (models.Article, error) {
	article, err := as.repo.Get(id)
	if err != nil {
		as.log.Error(fmt.Sprintf("article service error due to %s", err))
	}
	return article, err
}
func (as ArticleService) Filter(ctx context.Context, tag, date string) (models.TaggedArticles, error) {
	taggedArticles, err := as.repo.Filter(tag, date)
	if err != nil {
		as.log.Error(fmt.Sprintf("article service error due to %s", err))
	}

	return taggedArticles, err
}
