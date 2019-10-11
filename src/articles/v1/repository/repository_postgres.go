package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/willy182/boilerplate-go-cleanarch/src/articles/v1/model"
	"github.com/willy182/boilerplate-go-cleanarch/src/shared"
	"github.com/willy182/boilerplate-go-cleanarch/utils"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const tableName = "articles"

// postgresArticleRepo struct
type postgresArticleRepo struct {
	read  *gorm.DB
	write *gorm.DB
}

// NewPostgresArticleRepository article repository postgres handler
func NewPostgresArticleRepository(read, write *gorm.DB) Repository {
	// postgresConfig.InitDB()
	return &postgresArticleRepo{
		read:  read,
		write: write,
	}
}

// Save function, for save article object into database
func (r *postgresArticleRepo) Save(param *model.GormArticle) <-chan error {
	ctxRepo := "ArticleRepositorySave"

	output := make(chan error)

	go func() {
		// begin
		tx := r.write.Begin()

		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxRepo, "recover_repository_save")
				tx.Rollback()
				output <- fmt.Errorf(message)
			}
			close(output)
		}()

		if err := tx.Error; err != nil {
			utils.Log(log.ErrorLevel, err.Error(), ctxRepo, "tx_error")
			output <- err
		}

		// Select ID
		var id int
		row := r.read.Table(tableName).Where("id = ?", param.ID).Select("id").Row()
		row.Scan(&id)

		var errStmt error

		// force checking for auto increment number to insert or update
		if id > 0 {
			errStmt = tx.Table(tableName).Where("id = ?", param.ID).Updates(&param).Error
		} else {
			errStmt = tx.Table(tableName).Save(&param).Error
		}

		if errStmt != nil {
			utils.Log(log.ErrorLevel, errStmt.Error(), ctxRepo, "save_or_update_article")
			tx.Rollback()
			output <- errStmt
			return
		}

		tx.Commit()

		output <- nil
	}()

	return output
}

// GetByID function, for find article by its primary ID
func (r *postgresArticleRepo) GetByID(id int) <-chan ResultRepository {
	ctxRepo := "ArticleRepositoryGetByID"

	output := make(chan ResultRepository)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxRepo, "recover_repository_get_by_id")
				output <- ResultRepository{Error: fmt.Errorf(message)}
			}
			close(output)
		}()

		var (
			article   model.Article
			desc, img sql.NullString
			modified  pq.NullTime
		)

		row := r.read.Table(tableName).Where("id = ?", id).Select("id, title, summary, description, image, created, modified").Row()
		row.Scan(&article.ID, &article.Title, &article.Summary, &desc, &img, &article.Created, &modified)

		if article.ID == 0 {
			utils.Log(log.ErrorLevel, shared.ErrorRecordNotFound, ctxRepo, "recover_repository_get_by_id")
			output <- ResultRepository{Error: fmt.Errorf(shared.ErrorRecordNotFound)}
			return
		}

		if desc.Valid {
			article.Description = desc.String
		}

		if img.Valid {
			article.Image = img.String
		}

		if modified.Valid {
			article.Modified = modified.Time.Format(time.RFC3339)
		}

		output <- ResultRepository{Result: article}
	}()

	return output
}

// GetTotal function, for find all article
func (r *postgresArticleRepo) GetTotal(param model.QueryParamArticle) <-chan ResultRepository {
	ctxRepo := "ArticleRepositoryGetTotal"

	output := make(chan ResultRepository)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxRepo, "recover_repository_get_total")
				output <- ResultRepository{Error: fmt.Errorf(message)}
			}
			close(output)
		}()

		var total int

		q := r.read.Table(tableName)
		if param.Query != "" {
			q = q.Where("title LIKE ?", fmt.Sprintf("%%%s%%", param.Query))
		}

		err := q.Count(&total).Error
		if err != nil {
			utils.Log(log.ErrorLevel, err.Error(), ctxRepo, "query_err_rows")
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: total}
	}()

	return output
}

// GetAll function, for find all article
func (r *postgresArticleRepo) GetAll(param model.QueryParamArticle) <-chan ResultRepository {
	ctxRepo := "ArticleRepositoryGetAll"

	output := make(chan ResultRepository)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %v", r)
				utils.Log(log.ErrorLevel, message, ctxRepo, "recover_repository_get_total")
				output <- ResultRepository{Error: fmt.Errorf(message)}
			}
			close(output)
		}()

		var (
			articles []model.Article
			page     int
		)

		limit := shared.LimitDefault
		orderBy := "id"
		sortBy := "desc"

		if param.Limit != "" {
			limit, _ = strconv.Atoi(param.Limit)
		}

		if param.Page != "" {
			page, _ = strconv.Atoi(param.Page)
		}

		if param.OrderBy != "" {
			orderBy = param.OrderBy
		}

		if param.SortBy != "" {
			sortBy = param.SortBy
		}

		q := r.read.Table(tableName)
		if param.Query != "" {
			q = q.Where("title LIKE ?", "%"+param.Query+"%")
		}

		rows, err := q.Order(fmt.Sprintf("%s %s", orderBy, sortBy)).Limit(limit).Offset(page).Rows()
		if err != nil {
			utils.Log(log.ErrorLevel, err.Error(), ctxRepo, "query_err_rows")
			output <- ResultRepository{Error: err}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var (
				article   model.Article
				desc, img sql.NullString
				modified  pq.NullTime
			)

			err = rows.Scan(&article.ID, &article.Title, &article.Summary, &desc, &img, &article.Created, &modified)
			if err != nil {
				utils.Log(log.ErrorLevel, err.Error(), ctxRepo, "err_rows")
				output <- ResultRepository{Error: err}
				return
			}

			if desc.Valid {
				article.Description = desc.String
			}

			if img.Valid {
				article.Image = img.String
			}

			if modified.Valid {
				article.Modified = modified.Time.Format(time.RFC3339)
			}

			articles = append(articles, article)
		}

		output <- ResultRepository{Result: articles}
	}()

	return output
}
