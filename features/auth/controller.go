package auth

import (
	"errors"
	"time"

	"github.com/akhtarfath/domain"
	"github.com/gofiber/fiber/v3"
)

type Controller struct {
	tokens *TokenService
}

func NewController(tokens *TokenService) *Controller {
	return &Controller{tokens: tokens}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresAt string `json:"expires_at"`
}

func (h *Controller) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Invalid request body"))
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorMessage(fiber.StatusBadRequest, "Username and password are required"))
	}

	token, expiresAt, err := h.tokens.Authenticate(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorMessage(fiber.StatusUnauthorized, "Invalid username or password"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorMessage(fiber.StatusInternalServerError, "Internal server error"))
	}

	return c.Status(fiber.StatusOK).JSON(domain.NewSuccessResponse("Login successful", loginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}))
}
