package handler

import (
	"net/http"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/storage"
)

// GetSecret provides getting of the secret
func GetSecret(ctx context.Context, db storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Query("key")
		if key == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": "key is not defined"})
		}
		data, err := db.GetSecret(ctx, key)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"error": "false", "key": data.Key, "value": data.Value})
	}
}
