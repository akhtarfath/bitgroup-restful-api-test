package auth

import (
	"github.com/akhtarfath/config"
	"github.com/akhtarfath/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func AuthRoute(app *fiber.App, ctrl *Controller) {
	app.Post("/login", limiter.New(limiter.Config{
		Max:        config.RateLimitMax(),
		Expiration: config.RateLimitWindow(),
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(domain.NewErrorMessage(fiber.StatusTooManyRequests, "Too many login attempts, please try again later"))
		},
	}), ctrl.Login)
}
