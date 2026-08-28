package categories

import (
	"errors"

	"github.com/akhtarfath/domain"
	"github.com/gofiber/fiber/v3"
)

const internalErrMessage = "Internal server error"

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

type categoryRequest struct {
	Name string `json:"name"`
}

func (h *Controller) List(c fiber.Ctx) error {
	items, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Categories retrieved successfully", items))
}

func (h *Controller) Get(c fiber.Ctx) error {
	category, err := h.service.Get(c.Params("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorMessage(fiber.StatusNotFound, "Category not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Category retrieved successfully", category))
}

func (h *Controller) Create(c fiber.Ctx) error {
	var req categoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Invalid request body"))
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Category name is required"))
	}

	category, err := h.service.Create(CreateInput{Name: req.Name})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.NewCreatedResponse("Category created successfully", category))
}

func (h *Controller) Update(c fiber.Ctx) error {
	var req categoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Invalid request body"))
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Category name is required"))
	}

	category, err := h.service.Update(c.Params("id"), UpdateInput{Name: req.Name})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorMessage(fiber.StatusNotFound, "Category not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Category updated successfully", category))
}

func (h *Controller) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorMessage(fiber.StatusNotFound, "Category not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessMessage("Category deleted successfully"))
}
