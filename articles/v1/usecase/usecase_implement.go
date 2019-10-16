package usecase

import (
	"fmt"

	"github.com/willy182/boilerplate-go-cleanarch/articles/v1/model"
	"github.com/willy182/boilerplate-go-cleanarch/articles/v1/repository"
	"github.com/willy182/boilerplate-go-cleanarch/utils"

	log "github.com/sirupsen/logrus"
)

type articleUseCase struct {
	articleRepo repository.Repository
}

// NewArticleUseCase use case handler for article
func NewArticleUseCase(repo repository.Repository) UseCase {
	return &articleUseCase{
		articleRepo: repo,
	}
}

// Save use case handler for get article
func (u *articleUseCase) Save(param *model.GormArticle) <-chan error {
	ctxUsecase := "category_usecase_get_by_id"
	output := make(chan error)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxUsecase, "recover_usecase_save")
				output <- fmt.Errorf(message)
			}
			close(output)
		}()

		err := <-u.articleRepo.Save(param)
		if err != nil {
			utils.Log(log.ErrorLevel, err.Error(), ctxUsecase, "res_repo_save")
			output <- err
			return
		}

		output <- nil
	}()

	return output
}

// GetByID use case handler for get article by ID
func (u *articleUseCase) GetByID(ID int) <-chan ResultUseCase {
	ctxUsecase := "category_usecase_get_by_id"
	output := make(chan ResultUseCase)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxUsecase, "recover_usecase_get_by_id")
				output <- ResultUseCase{Error: fmt.Errorf(message)}
			}
			close(output)
		}()

		res := <-u.articleRepo.GetByID(ID)
		if res.Error != nil {
			utils.Log(log.ErrorLevel, res.Error.Error(), ctxUsecase, "res_repo_get_by_id")
			output <- ResultUseCase{Error: res.Error}
			return
		}

		response := res.Result.(model.Article)

		output <- ResultUseCase{Result: response}
	}()

	return output
}

// GetAll use case handler for find all get article
func (u *articleUseCase) GetAll(params model.QueryParamArticle) <-chan ResultUseCase {
	ctxUsecase := "category_usecase_get_all"
	output := make(chan ResultUseCase)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxUsecase, "recover_usecase_get_all")
				output <- ResultUseCase{Error: fmt.Errorf(message)}
			}
			close(output)
		}()

		var response ResponseUseCase

		resChan := u.articleRepo.GetAll(params)
		resTotalChan := u.articleRepo.GetTotal(params)

		res := <-resChan
		resTotal := <-resTotalChan

		if res.Error != nil {
			utils.Log(log.ErrorLevel, res.Error.Error(), ctxUsecase, "res_repo_get_all")
			output <- ResultUseCase{Error: res.Error}
			return
		}

		data := res.Result.([]model.Article)

		if resTotal.Error != nil {
			utils.Log(log.ErrorLevel, resTotal.Error.Error(), ctxUsecase, "res_repo_get_total")
			output <- ResultUseCase{Error: resTotal.Error}
			return
		}

		total := resTotal.Result.(int)

		response.Data = data
		response.Total = total
		output <- ResultUseCase{Result: response}
	}()

	return output
}
