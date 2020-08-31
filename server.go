package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type myServer struct {
	Port          int
	Auth0Audience string
	Logger        *logrus.Logger

	httpServer   *http.Server
	userStatuses map[string]string
}

func (ms *myServer) StartServer() error {
	ms.userStatuses = make(map[string]string)

	ms.httpServer = &http.Server{
		Addr:         ":3000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      ms.registerRoutes(),
	}

	ms.Logger.Printf("server starting on port: %d", ms.Port)

	err := ms.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (ms *myServer) GracefullyShutdown(ctx context.Context) error {
	return ms.httpServer.Shutdown(ctx)
}

func (ms *myServer) registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yessir"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(Authenticator)
		r.Get("/users/status", ms.handleListUserStatus)
		r.Get("/users/{userId}/status", ms.handleGetUserStatus)
		r.Put("/users/{userId}/status", ms.handlePutUserStatus)
	})

	return r
}
