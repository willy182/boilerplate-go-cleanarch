package main

import (
	"github.com/willy182/boilerplate-go-cleanarch/config"

	articleV1HTTP "github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/delivery"
	articleRepo "github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/repository"
	articleUseCase "github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/usecase"
)

// HSIService main service structure
type HSIService struct {
	Config  *config.Config
	Article struct {
		Usecase articleUseCase.UseCase
		Handler struct {
			V1 *articleV1HTTP.ArticleHandler
		}
	}
}

// InitHSIService function for initializing service
func InitHSIService(conf *config.Config) *HSIService {
	article := articleRepo.NewPostgresArticleRepository(conf.PostgresDB.Read, conf.PostgresDB.Write)
	articleUC := articleUseCase.NewArticleUseCase(article)
	articleV1Handler := articleV1HTTP.NewArticleHTTPHandler(articleUC)

	hsi := new(HSIService)
	hsi.Config = conf
	hsi.Article.Usecase = articleUC
	hsi.Article.Handler.V1 = articleV1Handler

	return hsi
}
