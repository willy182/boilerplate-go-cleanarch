package repository

import (
	"context"

	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// Repository interface for article repository
type Repository interface {
	Save(ctx context.Context, param *model.GormArticle) <-chan error
	GetByID(ctx context.Context, ID int) <-chan ResultRepository
	// GetAll(ctx context.Context, param model.Article) <-chan ResultRepository
	// GetTotal(ctx context.Context, param model.Article) <-chan ResultRepository
}
