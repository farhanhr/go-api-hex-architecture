package handler

import (
	"gonews/internal/adapter/handler/request"
	"gonews/internal/adapter/handler/response"
	"gonews/internal/core/domain/entity"
	"gonews/internal/core/service"
	"gonews/lib/conv"
	validatorLib "gonews/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultSuccessResponse response.DefaultSuccessResponse

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// CreateCategory implements CategoryHandler.
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateCategory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateCategory -3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.CreateCategory(c.Context(), reqEntity) 
	if err != nil {
		code = "[HANDLER] CreateCategory - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Data = nil
	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Category created successfuly"
	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

// DeleteCategory implements CategoryHandler.
func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	panic("unimplemented")
}

// EditCategory implements CategoryHandler.
func (ch *categoryHandler) EditCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] EditCategoriesByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] EditCategoriesByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] EditCategoriesByID -3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParams := c.Params("categoryID")
	id, err := conv.StringToInt64(idParams)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		ID: id,
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.EditCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Data = nil
	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Category updated successfuly"
	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

// GetCategories implements CategoryHandler.
func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userId := claims.UserID
	if userId == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID: result.ID,
			Title: result.Title,
			Slug: result.Slug,
			CreatedByName: result.User.Name,
		}
		categoryResponses = append(categoryResponses, categoryResponse)
	}
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Message = "categories fetched successfuly"
	defaultSuccessResponse.Data = categoryResponses

	return c.JSON(defaultSuccessResponse)
}

// GetCategoryByID implements CategoryHandler.
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userId := claims.UserID
	if userId == 0 {
		code = "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParams := c.Params("categoryID")
	id, err := conv.StringToInt64(idParams)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID: id,
		Title: result.Title,
		Slug: result.Slug,
		CreatedByName: result.User.Name,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Message = "categories detail fetched successfuly"
	defaultSuccessResponse.Data = categoryResponse

	return c.JSON(defaultSuccessResponse)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
