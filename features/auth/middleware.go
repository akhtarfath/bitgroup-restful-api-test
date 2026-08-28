package auth

import (
	"strings"

	"github.com/akhtarfath/domain"
	"github.com/gofiber/fiber/v3"
)

const bearerPrefix = "Bearer "

func Authenticate(tokens *TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")

		if !strings.HasPrefix(header, bearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorMessage(fiber.StatusUnauthorized, "Unauthorized"))
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
		if token == "" || !tokens.IsValid(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorMessage(fiber.StatusUnauthorized, "Unauthorized"))
		}

		return c.Next()
	}
}
