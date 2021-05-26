package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/service"
	"github.com/saromanov/secreto/internal/service/rest/handler"
	"github.com/saromanov/secreto/internal/storage"
	log "github.com/sirupsen/logrus"
)

type rest struct {
	srv *fiber.App
	st storage.Storage
}

func New(st storage.Storage) service.Service {
	return &rest{
		srv: newFiber(),
		st: st,
	}
}

func newFiber() *fiber.App {
	return fiber.New(fiber.Config{
		ReadTimeout:           10 * time.Minute,
		WriteTimeout:          5 * time.Minute,
		Prefork:               false,
		CaseSensitive:         false,
		StrictRouting:         true,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": "true", "Message": err.Error()})
		},
	})
}

func (s *rest) Run(ctx context.Context, ready func()) error {
	logger := log.WithContext(ctx)
	api := s.srv.Group("/api")
	api.Get("/secrets", handler.GetSecret(ctx, s.st))
	api.Post("/secrets", handler.CreateSecret(ctx, s.st))
	if err := s.srv.Listen(":8089"); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		logger.WithError(err).Error("Failed start listening")
		return err
	}
	logger.WithField("address", ":8089").Info("Start listening")
	defer func() {
		logger.Info("Stop listening")
	}()

	return nil
}

func (s *rest) Shutdown(ctx context.Context) error {
	return nil
}
