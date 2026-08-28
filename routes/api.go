package routes

import (
	"github.com/akhtarfath/features/auth"
	"github.com/akhtarfath/features/categories"
	"github.com/akhtarfath/features/products"
	"github.com/gofiber/fiber/v3"
)

func New(app *fiber.App) {
	tokens := auth.NewTokenService()
	authController := auth.NewController(tokens)
	auth.AuthRoute(app, authController)

	authenticate := auth.Authenticate(tokens)

	categoryStore := categories.NewStore()
	categoryService := categories.NewService(categoryStore)
	categoryController := categories.NewController(categoryService)
	categories.Route(app, categoryController, authenticate)

	productStore := products.NewStore()
	productService := products.NewService(productStore)
	productController := products.NewController(productService, categoryService)
	products.Route(app, productController, authenticate)
}
