package products

import (
	"github.com/gofiber/fiber/v3"
)

func Route(app *fiber.App, ctrl *Controller, authenticate fiber.Handler) {
	app.Get("/products", authenticate, ctrl.List)
	app.Get("/products/:id", authenticate, ctrl.Get)
	app.Post("/products", authenticate, ctrl.Create)
	app.Put("/products/:id", authenticate, ctrl.Update)
	app.Delete("/products/:id", authenticate, ctrl.Delete)
}
