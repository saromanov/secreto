package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/crypto"
	"github.com/saromanov/secreto/internal/storage"
	log "github.com/sirupsen/logrus"
)

// GetSecret provides getting of the secret
func GetSecret(ctx context.Context, secretKey string, db storage.Storage) fiber.Handler {
	logger := log.WithContext(ctx)
	return func(c *fiber.Ctx) error {
		key := c.Query("key")
		if key == "" {
			logger.Error("key is not definded")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": "key is not defined"})
		}
		data, err := db.GetSecret(ctx, key)
		if err != nil {
			logger.WithError(err).Error("unable to get secret")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		result, err := crypto.Decrypt(data.Value, secretKey)
		if err != nil {
			logger.WithError(err).Error("unable to decrypt message")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "true", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"error": "false", "key": data.Key, "value": result})
	}
}
