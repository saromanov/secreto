package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/crypto"
	"github.com/saromanov/secreto/internal/models"
	"github.com/saromanov/secreto/internal/storage"
	log "github.com/sirupsen/logrus"
)

// CreateSecret provides creating of the secret
func CreateSecret(ctx context.Context, secretKey string, db storage.Storage) fiber.Handler {
	logger := log.WithContext(ctx)
	return func(c *fiber.Ctx) error {
		secret := new(models.SecretRESTPost)
		if err := c.BodyParser(&secret); err != nil {
			logger.WithError(err).Error("unable to parse request")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		if err := secret.Validate(); err != nil {
			logger.WithError(err).Error("unable to validate input data")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}

		result, err := crypto.Encrypt(secret.Value, secretKey)
		if err != nil {
			logger.WithError(err).Error("unable to encrypt data")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		if err := db.CreateSecret(ctx, models.Secret{
			Key:   secret.Key,
			Value: result,
		}); err != nil {
			logger.WithError(err).Error("unable to create secret")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"error": "false", "status": "ok"})
	}
}
