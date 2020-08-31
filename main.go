package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
)

var (
	logger    *logrus.Logger
	tokenAuth *jwtauth.JWTAuth
)

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"userId": "a1b2c3d4"})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {
	logger = logrus.New()

	s := myServer{
		Port: 3000,
		// TODO: add real audience
		Auth0Audience: "my-cool-audience",

		Logger: logger,
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
