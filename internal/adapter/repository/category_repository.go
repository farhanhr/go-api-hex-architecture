package repository

import (
	"context"
	"errors"
	"gonews/internal/core/domain/entity"
	"gonews/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	panic("unimplemented")
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// EditCategory implements CategoryRepository.
func (c *categoryRepository) EditCategory(ctx context.Context, req entity.CategoryEntity) error {
	panic("unimplemented")
}

// GetCategories implements CategoryRepository.
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCategories []model.Category

	err = c.db.Order("created_at DESC").Preload("User").Find(&modelCategories).Error
	if err != nil {
		code = "[REPOSITORY] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCategories) == 0 {
		code = "[REPOSITORY] GetCategories - 2"
		err = errors.New("categories data is not found")
		log.Errorw(code, err)
		return nil, err
	}

	var resp []entity.CategoryEntity
	for _, val := range modelCategories {
		resp = append(resp, entity.CategoryEntity{
			ID: val.ID,
			Title: val.Title,
			Slug: val.Slug,
			User: entity.UserEntity{
				ID: val.User.ID,
				Name: val.User.Name,
				Email: val.User.Email,
			},
		})
	}

	return resp, nil
}

// GetCategoryByID implements CategoryRepository.
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	panic("unimplemented")
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
