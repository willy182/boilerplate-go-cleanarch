package usecase

import (
	"context"
	"fmt"

	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/repository"
	"github.com/willy182/boilerplate-go-cleanarch/utils"

	log "github.com/sirupsen/logrus"
)

type articleUseCase struct {
	articleRepo repository.Repository
}

// NewArticleUseCase use case handler for category
func NewArticleUseCase(repo repository.Repository) UseCase {
	return &articleUseCase{
		articleRepo: repo,
	}
}

// Save use case handler for get category by ID
func (u *articleUseCase) Save(ctx context.Context, param *model.GormArticle) <-chan error {
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

		err := <-u.articleRepo.Save(ctx, param)
		if err != nil {
			utils.Log(log.ErrorLevel, err.Error(), ctxUsecase, "res_repo_save")
			output <- err
			return
		}

		output <- nil
	}()

	return output
}

// GetByID use case handler for get category by ID
func (u *articleUseCase) GetByID(ctx context.Context, ID int) <-chan ResultUseCase {
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

		res := <-u.articleRepo.GetByID(ctx, ID)
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
