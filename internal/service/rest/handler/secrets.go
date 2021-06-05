package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/storage"
)

// GenerateSecret provides generation of teh secrets
func GenerateSecret(ctx context.Context, secretKey string, db storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"error": "false", "result": "ok"})
	}
}
