package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"

	"github.com/saromanov/secreto/internal/service"
	"github.com/saromanov/secreto/internal/service/rest"
	"github.com/saromanov/secreto/internal/storage"
	"github.com/saromanov/secreto/internal/storage/badger"
)

type config struct {
	Storage storage.Config
	Rest    rest.Config
}

func main() {
	var cfg config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		logrus.WithError(err).Fatal("Cannot decode config envs")
	}
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	ctx, cancel := context.WithCancel(context.Background())
	g := &run.Group{}
	{
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}

	st, err := badger.New(cfg.Storage)
	if err != nil {
		logger.WithError(err).Fatal("unable to init storage")
	}
	r := rest.New(cfg.Rest, st)

	s := service.Runner{}
	if err := s.SetupService(ctx, r, "rest", g); err != nil {
		logger.WithError(err).Fatal("unable to setup service ")
	}
	logger.Info("Running of the service...")
	if err := g.Run(); err != nil {
		logger.WithError(err).Fatal("The service has been stopped with error")
	}
	logger.Info("Service is stopped")

}
