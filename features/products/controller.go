package products

import (
	"errors"

	"github.com/akhtarfath/domain"
	"github.com/gofiber/fiber/v3"

	"github.com/akhtarfath/features/categories"
)

const internalErrMessage = "Internal server error"

type Controller struct {
	service       *Service
	categoriesSvc *categories.Service
}

func NewController(service *Service, categoriesSvc *categories.Service) *Controller {
	return &Controller{service: service, categoriesSvc: categoriesSvc}
}

type productRequest struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
}

func (h *Controller) List(c fiber.Ctx) error {
	items, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Products retrieved successfully", items))
}

func (h *Controller) Get(c fiber.Ctx) error {
	product, err := h.service.Get(c.Params("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorMessage(fiber.StatusNotFound, "Product not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Product retrieved successfully", product))
}

func (h *Controller) Create(c fiber.Ctx) error {
	var req productRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Invalid request body"))
	}

	if req.CategoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Category ID is required"))
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Product name is required"))
	}
	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Product price cannot be negative"))
	}

	category, err := h.categoriesSvc.Get(req.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Category not found"))
	}

	product, err := h.service.Create(CreateInput{
		Name:     req.Name,
		Price:    req.Price,
		Category: category,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.NewCreatedResponse("Product created successfully", product))
}

func (h *Controller) Update(c fiber.Ctx) error {
	var req productRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Invalid request body"))
	}

	if req.CategoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Category ID is required"))
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Product name is required"))
	}
	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Product price cannot be negative"))
	}

	category, err := h.categoriesSvc.Get(req.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Category not found"))
	}

	product, err := h.service.Update(c.Params("id"), UpdateInput{
		Name:     req.Name,
		Price:    req.Price,
		Category: category,
	})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorMessage(fiber.StatusNotFound, "Product not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Product updated successfully", product))
}

func (h *Controller) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorMessage(fiber.StatusNotFound, "Product not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, internalErrMessage))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessMessage("Product deleted successfully"))
}
