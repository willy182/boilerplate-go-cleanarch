package delivery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/willy182/boilerplate-go-cleanarch/articles/v1/model"
	"github.com/willy182/boilerplate-go-cleanarch/articles/v1/usecase"
	"github.com/willy182/boilerplate-go-cleanarch/lib"
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
	group.POST("/article", h.Create)
	group.GET("/article/:id", h.GetByID)
	group.GET("/article", h.GetAll)
}

// Create method for handling route save article
func (h *ArticleHandler) Create(c *gin.Context) {
	ctxHandler := "article_handler_create"
	multiError := lib.NewMultiError()

	params := &model.ArticleInput{}
	if err := c.Bind(params); err != nil {
		multiError.Append("bindParam", err)
		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "bind_param")
		response := lib.NewHTTPResponse(http.StatusBadRequest, "error bind param", multiError)
		response.JSON(c.Writer)
		return
	}

	multiError.Clear()

	multiError = lib.Validate("article_create_params", params)
	if multiError != nil && multiError.HasError() {
		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "validate_params")
		response := lib.NewHTTPResponse(http.StatusBadRequest, "validate params", multiError)
		response.JSON(c.Writer)
		return
	}

	paramDB := &model.GormArticle{
		Title:       params.Title,
		Summary:     params.Summary,
		Description: params.Description,
		Image:       params.Image,
		Created:     time.Now(),
	}

	err := <-h.ArticleUseCase.Save(paramDB)
	if err != nil {
		utils.Log(log.ErrorLevel, err.Error(), ctxHandler, "err_res_save")
		response := lib.NewHTTPResponse(http.StatusBadRequest, err.Error())
		response.JSON(c.Writer)
		return
	}

	response := lib.NewHTTPResponse(http.StatusOK, "Data has been save")
	response.JSON(c.Writer)
	return
}

// GetByID method for handling route article by ID
func (h *ArticleHandler) GetByID(c *gin.Context) {
	ctxHandler := "article_handler_get_by_id"
	idParam := c.Param("id")
	multiError := lib.NewMultiError()

	if ok := lib.ValidateNumeric(idParam); !ok {
		multiError.Append("error", fmt.Errorf("id must be numeric"))
		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "validate_id")
		response := lib.NewHTTPResponse(http.StatusBadRequest, "validate id", multiError)
		response.JSON(c.Writer)
		return
	}

	multiError.Clear()

	id, _ := strconv.Atoi(idParam)
	res := <-h.ArticleUseCase.GetByID(id)
	if res.Error != nil && res.Error.Error() == lib.ErrorRecordNotFound {
		utils.Log(log.ErrorLevel, res.Error.Error(), ctxHandler, "record_not_found")
		response := lib.NewHTTPResponse(http.StatusOK, res.Error.Error())
		response.JSON(c.Writer)
		return
	} else if res.Error != nil {
		utils.Log(log.ErrorLevel, res.Error.Error(), ctxHandler, "err_res_get_by_id")
		response := lib.NewHTTPResponse(http.StatusBadRequest, res.Error.Error(), multiError)
		response.JSON(c.Writer)
		return
	}

	result := res.Result.(model.Article)
	meta := lib.CreateMeta(1, 1, 1)
	response := lib.NewHTTPResponse(http.StatusOK, "Article Get By ID", result, meta)
	response.JSON(c.Writer)
	return
}

// GetAll method for handling route for get article list
func (h *ArticleHandler) GetAll(c *gin.Context) {
	ctxHandler := "article_handler_get_all"
	multiError := lib.NewMultiError()
	req := c.Request

	var params model.QueryParamArticle

	err := lib.BindQueryParam(req.URL, &params)
	if err != nil {
		multiError.Append("bindError", err)
		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "bind_params")
		response := lib.NewHTTPResponse(http.StatusBadRequest, "bind params", multiError)
		response.JSON(c.Writer)
		return
	}

	multiError.Clear()

	params.Query, _ = url.PathUnescape(params.Query)

	multiError = lib.Validate("article_get_params", params)
	if multiError != nil && multiError.HasError() {
		utils.Log(log.ErrorLevel, multiError.Error(), ctxHandler, "validate_params")
		response := lib.NewHTTPResponse(http.StatusBadRequest, "validate params", multiError)
		response.JSON(c.Writer)
		return
	}

	if params.Limit == "0" {
		params.Limit = ""
	}

	if params.Page == "0" {
		params.Page = ""
	}

	res := <-h.ArticleUseCase.GetAll(params)
	if res.Error != nil {
		utils.Log(log.ErrorLevel, res.Error.Error(), ctxHandler, "validate_params")
		response := lib.NewHTTPResponse(http.StatusBadRequest, res.Error.Error())
		response.JSON(c.Writer)
		return
	}

	result := res.Result.(usecase.ResponseUseCase)

	if result.Total == 0 {
		msg := "Search results empty"
		utils.Log(log.ErrorLevel, msg, ctxHandler, "list_empty")
		response := lib.NewHTTPResponse(http.StatusOK, msg)
		response.JSON(c.Writer)
		return
	}

	var (
		page  = 1
		limit = lib.LimitDefault
	)

	if params.Page != "" {
		page, _ = strconv.Atoi(params.Page)
	}

	if params.Limit != "" {
		limit, _ = strconv.Atoi(params.Limit)
	}

	meta := lib.CreateMeta(result.Total, page, limit)
	response := lib.NewHTTPResponse(http.StatusOK, "Article List", result.Data, meta)
	response.JSON(c.Writer)
	return
}
