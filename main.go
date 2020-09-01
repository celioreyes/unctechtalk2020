package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func main() {
	logger = logrus.New()

	port := os.Getenv("PORT")

	verifier := &Verifier{
		Resource:  os.Getenv("AuthResource"),
		TenantURL: os.Getenv("AuthTenantURL"),
	}

	s := myServer{
		Port:     port,
		Verifier: verifier,
		Logger:   logger,
	}

	go func() {
		if err := s.StartServer(); err != nil {
			logger.WithError(err).Error("could not start server")
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Info("gracefully shutting down server")
	s.httpServer.Shutdown(context.Background())
}
