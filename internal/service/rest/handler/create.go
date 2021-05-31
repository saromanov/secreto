package handler

import (
	"net/http"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/models"
	"github.com/saromanov/secreto/internal/storage"
)

// CreateSecret provides creating of the secret
func CreateSecret(ctx context.Context, secret string, db storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := new(models.SecretRESTPost)
		if err := c.BodyParser(&secret); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		if err := secret.Validate(); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}

		if err := db.CreateSecret(ctx, models.Secret{
			Key: secret.Key,
			Value: secret.Value,
		}); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"error": "false", "status": "ok"})
	}
}
