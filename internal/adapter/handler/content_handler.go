package handler

import (
	"gonews/internal/adapter/handler/response"
	"gonews/internal/core/domain/entity"
	"gonews/internal/core/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentById(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// CreateContent implements ContentHandler.
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteContent implements ContentHandler.
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContentById implements ContentHandler.
func (ch *contentHandler) GetContentById(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContents implements ContentHandler.
func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	result, err := ch.contentService.GetContents(c.Context())
	if err != nil {
		code = "[HANDLER] GetContents - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Message = "contents fetched successfuly"

	respContents := []response.ContentResponse{}
	for _, content := range result {
		respContent := response.ContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedById:  content.CreatedById,
			CreatedAt:    content.CreatedAt.Format(time.RFC3339),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}

		respContents = append(respContents, respContent)
	}

	defaultSuccessResponse.Data = respContents

	return c.JSON(defaultSuccessResponse)
}

// UpdateContent implements ContentHandler.
func (ch *contentHandler) UpdateContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentHandler.
func (ch *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
