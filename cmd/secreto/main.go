package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"

	"github.com/saromanov/secreto/internal/service"
	"github.com/saromanov/secreto/internal/service/rest"
	"github.com/saromanov/secreto/internal/storage"
	"github.com/saromanov/secreto/internal/storage/badger"
)

func main() {
	logger := logrus.New()
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

	st, err := badger.New(storage.Config{})
	if err != nil {
		log.WithError(err).Fatal("unable to init storage")
	}
	r := rest.New(st)

	s := service.Runner{}
	if err := s.SetupService(ctx, r, "rest", g); err != nil {
		log.WithError(err).Fatal("unable to setup service ")
	}
	logger.Info("Running of the service...")
	if err := g.Run(); err != nil {
		logger.WithError(err).Fatal("The service has been stopped with error")
	}
	logger.Info("Service is stopped")

}
