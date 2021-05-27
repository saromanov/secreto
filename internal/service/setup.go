package service

import (
	"context"
	"sync"

	"github.com/oklog/run"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	prev chan struct{}
}

// SetupService provides setup of the service
func (o *Runner) SetupService(ctx context.Context, srv Service, role string, g *run.Group) error {
	logger := log.New()
	chPrev := o.prev
	chNext := make(chan struct{})
	var once sync.Once
	closer := func() {
		once.Do(func() {
			close(chNext)
		})
	}
	g.Add(func() error {
		if chPrev != nil {
			<-chPrev
		}
		logger.Info("Running the service...")
		defer logger.Info("˜the service is stopped")
		if err := ctx.Err(); err != nil {
			return nil
		}
		return srv.Run(ctx, closer)
	}, func(error) {
		closer()
		logger.Info("Shutdowning the service...")
		defer logger.Info("˜the service is shutdown")
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("Cannot shutdown the service properly")
		}
	})
	o.prev = chNext
	return nil
}
