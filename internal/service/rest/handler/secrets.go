package handler

import (
	"context"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GenerateSecretKey provides generation of the secrets
func GenerateSecretKey(ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result := strings.Replace(uuid.New().String(), "-", "", -1)
		return c.JSON(fiber.Map{"error": "false", "result": result})
	}
}
