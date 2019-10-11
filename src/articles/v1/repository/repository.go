package repository

import (
	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// Repository interface for article repository
type Repository interface {
	Save(param *model.GormArticle) <-chan error
	GetByID(ID int) <-chan ResultRepository
	GetTotal(param model.QueryParamArticle) <-chan ResultRepository
	GetAll(param model.QueryParamArticle) <-chan ResultRepository
}
