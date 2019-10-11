package usecase

import (
	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result interface{}
	Error  error
}

// ResponseUseCase data structure
type ResponseUseCase struct {
	Data  interface{}
	Total int
}

// UseCase use case for category
type UseCase interface {
	Save(param *model.GormArticle) <-chan error
	GetByID(ID int) <-chan ResultUseCase
	GetAll(params model.QueryParamArticle) <-chan ResultUseCase
}
