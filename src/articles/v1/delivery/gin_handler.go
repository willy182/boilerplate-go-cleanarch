package delivery

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/usecase"
	"github.com/willy182/boilerplate-go-cleanarch/src/shared"
	"github.com/willy182/boilerplate-go-cleanarch/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ArticleHandler struct for http brand handling
type ArticleHandler struct {
	ArticleUseCase usecase.UseCase
}

// NewArticleHTTPHandler route handler for article
func NewArticleHTTPHandler(usecase usecase.UseCase) *ArticleHandler {
	return &ArticleHandler{ArticleUseCase: usecase}
}

// Mount function
func (h *ArticleHandler) Mount(group *gin.RouterGroup) {
	// group.POST("/article", h.GetAll)
	group.GET("/article/:id", h.GetByID)
}

// GetByID method for handling route article by ID
func (h *ArticleHandler) GetByID(c gin.Context) error {
	ctxHandler := "article_handler_get_by_id"
	idParam := c.Param("id")
	ctx := context.Context()
	multiError := shared.NewMultiError()

	if ok := shared.ValidateNumeric(idParam); !ok {
		multiError.Append("error", fmt.Errorf("id must be numeric"))
		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "validate_id")
		response := shared.NewHTTPResponse(http.StatusBadRequest, "validate id", multiError)
		return response.JSON(c.Response())
	}

	multiError.Clear()

	id, _ := strconv.Atoi(idParam)
	res := <-h.ArticleUseCase.GetByID(ctx, id)
	if res.Error != nil {
		utils.Log(log.ErrorLevel, res.Error.Error(), ctxHandler, "err_res_get_by_id")
		response := shared.NewHTTPResponse(http.StatusBadRequest, res.Error.Error(), multiError)
		return response.JSON(c.Response())
	}

	result := res.Result.(model.Article)
	meta := shared.CreateMeta(1, 1, 1)
	response := shared.NewHTTPResponse(http.StatusOK, "Article Get By ID", result, meta)
	return response.JSON(c.Response())
}

// GetAll method for handling route for get article list
// func (h *ArticleHandler) GetAll(c gin.Context) error {
// 	ctxHandler := "article_handler_get_all"
// 	ctx := context.Context()
// 	multiError := shared.NewMultiError()
// 	var params model.CategoryRequestParams

// 	err := shared.BindQueryParam(req.URL, &params)
// 	if err != nil {
// 		multiError.Append("error", err)
// 		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "bind_params")
// 		response := shared.NewHTTPResponse(http.StatusBadRequest, "bind params", multiError)
// 		return response.JSON(c.Response())
// 	}

// 	multiError.Clear()

// 	params.FilterQuery, _ = url.PathUnescape(params.FilterQuery)

// 	multiError = jsonschema.Validate("article_get_params", params)
// 	if multiError != nil && multiError.HasError() {
// 		golib.Log(golib.ErrorLevel, multiError.Error(), ctxHandler, "validate_params")
// 		tracer.Log(ctx, ctxHandler, "validate_params", params)
// 		response := golib.NewHTTPResponse(http.StatusBadRequest, "validate params", multiError)
// 		return response.JSON(c.Response())
// 	}

// 	if params.PageSize == "0" {
// 		params.PageSize = ""
// 	}

// 	if params.PageNumber == "0" {
// 		params.PageNumber = ""
// 	}

// 	var categoryParams model.CategoryParams
// 	categoryParams.FilterQuery = params.FilterQuery
// 	categoryParams.FilterOldID = params.FilterOldID
// 	categoryParams.FilterStatus = params.FilterStatus
// 	categoryParams.FilterCurrentLevel = params.FilterCurrentLevel
// 	categoryParams.FieldsCategories = fieldsCategories
// 	categoryParams.Sort = params.Sort
// 	categoryParams.PageSize = params.PageSize
// 	categoryParams.PageNumber = params.PageNumber
// 	categoryParams.Include = params.Include
// 	categoryParams.NoCache = params.NoCache

// 	res := <-h.ArticleUseCase.GetAll(ctx, categoryParams, req)
// 	if res.Error != nil {
// 		golib.Log(golib.ErrorLevel, res.Error.Error(), ctxHandler, "err_res_get_all")
// 		tracer.Log(ctx, ctxHandler, "err_res_get_all", res.Error.Error(), categoryParams)
// 		response := golib.NewHTTPResponse(http.StatusBadRequest, res.Error.Error(), multiError)
// 		return response.JSON(c.Response())
// 	}

// 	result := res.Result.([]model.Article)

// 	var (
// 		page  = 1
// 		limit = helpers.LimitDefault
// 	)

// 	if params.PageNumber != "" {
// 		page, _ = strconv.Atoi(params.PageNumber)
// 	}

// 	if params.PageSize != "" {
// 		limit, _ = strconv.Atoi(params.PageSize)
// 	}

// 	meta := shared.CreateMeta(result.Total, page, limit)
// 	response := golib.NewHTTPResponse(http.StatusOK, "Article List", result.Data, meta)
// 	return response.JSON(c.Response())
// }
