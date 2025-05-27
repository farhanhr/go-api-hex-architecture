package service

import (
	"context"
	"gonews/internal/adapter/repository"
	"gonews/internal/core/domain/entity"
	"gonews/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug

	err =  c.categoryRepository.CreateCategory(ctx, req)
	if err != nil {
		code = "[Service] GetCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteCategory implements CategoryService.
func (c *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// EditCategory implements CategoryService.
func (c *categoryService) EditCategory(ctx context.Context, req entity.CategoryEntity) error {
	categoryData, err := c.categoryRepository.GetCategoryByID(ctx, req.ID)
	if err != nil {
		code = "[SERVICE] EditCategoryByID - 1"
		log.Errorw(code, err)
		return err
	}
	slug := conv.GenerateSlug(req.Title)
	if categoryData.Title == req.Title {
		slug = categoryData.Slug
	}
	req.Slug = slug

	err = c.EditCategory(ctx, req)
	if err != nil {
		code = "[SERVICE] EditCategoryByID - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetCategories implements CategoryService.
func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	result, err := c.categoryRepository.GetCategories(ctx)
	if err != nil {
		code = "[Service] GetCategory - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, err
}

// GetCategoryByID implements CategoryService.
func (c *categoryService) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	result, err := c.categoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		code = "[Service] GetCategory - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepo,
	}
}
