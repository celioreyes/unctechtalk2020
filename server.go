package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type myServer struct {
	Port     int
	Verifier IVerifier
	Logger   *logrus.Logger

	httpServer *http.Server
	userMoods  map[string]Mood
	moods      map[int]Mood
}

func (ms *myServer) StartServer() error {
	ms.userMoods = make(map[string]Mood)
	ms.moods = map[int]Mood{
		1: {ID: 1, Name: "Happy"},
		2: {ID: 2, Name: "Excited"},
		3: {ID: 3, Name: "Sad"},
	}

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
		r.Use(CheckJWT(ms.Verifier))
		r.With(CheckPermissions("read:mood")).Get("/moods", ms.handleListMoods)
		r.With(CheckPermissions("read:user:mood")).Get("/users/{userId}/mood", ms.handleGetUserMood)
		r.With(CheckPermissions("write:user:mood")).Put("/users/{userId}/mood", ms.handlePutUserMood)
	})

	return r
}
