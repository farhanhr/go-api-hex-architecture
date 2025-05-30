package repository

import (
	"context"
	"gonews/internal/core/domain/entity"
	"gonews/internal/core/domain/model"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ContentRepository interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	GetContentById(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

// CreateContent implements ContentRepository.
func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// DeleteContent implements ContentRepository.
func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetContentById implements ContentRepository.
func (c *contentRepository) GetContentById(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	panic("unimplemented")
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	var modelContents []model.Content

	err = c.db.Order("created_at DESC").Preload("User", "Category").Find(&modelContents).Error
	if err != nil {
		code = "[REPOSITORY] GetContents - 1"
		log.Errorw(code, err)
		return nil, err 
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags :=  strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Image:       val.Image,
			Tags:        tags,
			Status:      val.Description,
			CategoryID:  val.CategoryID,
			CreatedById: val.CreatedByID,
			CreatedAt:   val.CreatedAt,
			Category:    entity.CategoryEntity{
				ID:    val.CategoryID,
				Title: val.Title,
				Slug:  val.Description,
			},
			User:        entity.UserEntity{
				ID:       val.User.ID,
				Name:     val.User.Name,
			},
		}

		resps = append(resps, resp)
	}
	
	return resps, nil
}

// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
