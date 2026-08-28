package categories

import (
	"github.com/gofiber/fiber/v3"
)

func Route(app *fiber.App, ctrl *Controller, authenticate fiber.Handler) {
	app.Get("/categories", authenticate, ctrl.List)
	app.Get("/categories/:id", authenticate, ctrl.Get)
	app.Post("/categories", authenticate, ctrl.Create)
	app.Put("/categories/:id", authenticate, ctrl.Update)
	app.Delete("/categories/:id", authenticate, ctrl.Delete)
}
