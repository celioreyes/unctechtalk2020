package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type myServer struct {
	Port     string
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
		Addr:         fmt.Sprintf(":%s", ms.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      ms.registerRoutes(),
	}

	ms.Logger.Printf("server starting on port: %s", ms.Port)

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

	// Basic CORS -- for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// HTTP Logger
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
