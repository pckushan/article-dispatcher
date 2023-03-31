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

func (ar ArticleService) Create(ctx context.Context, article models.Article) error {
	err := ar.repo.Set(article)
	if err != nil {
		ar.log.Error(fmt.Sprintf("article service error due to %s", err))
	}
	return err
}
