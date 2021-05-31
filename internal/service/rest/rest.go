package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/saromanov/secreto/internal/service"
	"github.com/saromanov/secreto/internal/service/rest/handler"
	"github.com/saromanov/secreto/internal/storage"
	log "github.com/sirupsen/logrus"
)

type rest struct {
	cfg Config
	srv *fiber.App
	st  storage.Storage
}

// New provides initialization of the rest service
func New(cfg Config, st storage.Storage) service.Service {
	return &rest{
		cfg: cfg,
		srv: newFiber(),
		st:  st,
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
	api.Get("/secrets", handler.GetSecret(ctx, s.cfg.Secret, s.st))
	api.Post("/secrets", handler.CreateSecret(ctx, s.cfg.Secret, s.st))
	logger.WithField("address", s.cfg.Port).Info("Start listening")
	if err := s.srv.Listen(fmt.Sprintf(":%d", s.cfg.Port)); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		logger.WithError(err).Error("Failed start listening")
		return err
	}
	defer func() {
		logger.Info("Stop listening")
	}()

	return nil
}

func (s *rest) Shutdown(ctx context.Context) error {
	return nil
}
