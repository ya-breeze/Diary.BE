// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Diary - OpenAPI 3.0
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Contact: ilya.korolev@outlook.com
 */

package goserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/ya-breeze/diary.be/pkg/config"
)

type CustomControllers struct {
	AuthAPIService AuthAPIService
	UserAPIService UserAPIService
}

func Serve(ctx context.Context, logger *slog.Logger, cfg *config.Config,
	controllers CustomControllers, extraRouters []Router, middlewares ...mux.MiddlewareFunc) (net.Addr, chan int, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to listen: %w", err)
	}
	logger.Info(fmt.Sprintf("Listening at port %d...", listener.Addr().(*net.TCPAddr).Port))

	AuthAPIService := NewAuthAPIService()
	if controllers.AuthAPIService != nil {
		AuthAPIService = controllers.AuthAPIService
	}
	AuthAPIController := NewAuthAPIController(AuthAPIService)

	UserAPIService := NewUserAPIService()
	if controllers.UserAPIService != nil {
		UserAPIService = controllers.UserAPIService
	}
	UserAPIController := NewUserAPIController(UserAPIService)

	routers := append(extraRouters, AuthAPIController, UserAPIController)
	router := NewRouter(routers...)

	router.Use(middlewares...)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	server := &http.Server{
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(router),
	}

	go func() {
		server.Serve(listener)
	}()

	finishChan := make(chan int, 1)
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server...")
		timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(timeout)
		finishChan <- 1
		logger.Info("Server stopped")
	}()

	return listener.Addr(), finishChan, nil
}
