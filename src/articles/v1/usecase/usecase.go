package usecase

import (
	"context"

	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result interface{}
	Error  error
}

// UseCase use case for category
type UseCase interface {
	Save(ctx context.Context, param *model.GormArticle) <-chan error
	GetByID(ctx context.Context, ID int) <-chan ResultUseCase
	// GetAll(ctx context.Context, params model.CategoryParams, req *http.Request) <-chan ResultUseCase
}
